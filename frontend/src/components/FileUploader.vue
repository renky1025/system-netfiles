<template>
  <el-upload
    class="upload-demo"
    action="#"
    :http-request="customUpload"
    :show-file-list="false"
    multiple
  >
    <el-button type="primary">Upload File</el-button>
  </el-upload>
</template>

<script setup lang="ts">
import { defineEmits, defineProps } from 'vue';
import { checkFile, instantUpload, uploadChunk, mergeChunks } from '../api/file';
import { ElMessage } from 'element-plus';
import { v4 as uuidv4 } from 'uuid';
import SparkMD5 from 'spark-md5';

const props = defineProps<{
  currentFolderId: number | null;
}>();

const emit = defineEmits(['upload-success']);

const CHUNK_SIZE = 2 * 1024 * 1024; // 2MB

const calculateMD5 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const blobSlice = File.prototype.slice || (File.prototype as any).mozSlice || (File.prototype as any).webkitSlice;
    const chunks = Math.ceil(file.size / CHUNK_SIZE);
    let currentChunk = 0;
    const spark = new SparkMD5.ArrayBuffer();
    const fileReader = new FileReader();

    fileReader.onload = function (e) {
      spark.append(e.target?.result as ArrayBuffer);
      currentChunk++;

      if (currentChunk < chunks) {
        loadNext();
      } else {
        resolve(spark.end());
      }
    };

    fileReader.onerror = function () {
      reject('MD5 calculation failed');
    };

    function loadNext() {
      const start = currentChunk * CHUNK_SIZE;
      const end = ((start + CHUNK_SIZE) >= file.size) ? file.size : start + CHUNK_SIZE;
      fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
    }

    loadNext();
  });
};

const customUpload = async (options: any) => {
  const file = options.file;
  
  try {
    // 1. Calculate MD5
    const md5 = await calculateMD5(file);

    // 2. Check if file exists
    const checkRes = await checkFile({ md5 });
    
    if (checkRes.data?.exists) {
      // 3. Instant Upload
      const uploadRes = await instantUpload({
        md5,
        file_name: file.name,
        file_size: file.size,
        folder_id: props.currentFolderId,
      });
      if (uploadRes.code === 200 || uploadRes.code === 0) {
        ElMessage.success('Instant upload successful');
        emit('upload-success');
        return;
      } else {
        ElMessage.error(uploadRes.msg || uploadRes.error || 'Instant upload failed');
        return;
      }
    }

    // 4. Normal Upload
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    const uploadId = uuidv4();

    for (let i = 0; i < totalChunks; i++) {
      const start = i * CHUNK_SIZE;
      const end = Math.min(file.size, start + CHUNK_SIZE);
      const chunk = file.slice(start, end);

      await uploadChunk({
        upload_id: uploadId,
        index: i,
        file: chunk,
      });
    }

    // Merge
    const mergeRes = await mergeChunks({
      upload_id: uploadId,
      file_name: file.name,
      total_chunks: totalChunks,
      folder_id: props.currentFolderId,
    });

    if (mergeRes.code === 200 || mergeRes.code === 0) {
      ElMessage.success('Upload successful');
      emit('upload-success');
    } else {
      ElMessage.error(mergeRes.msg || mergeRes.error || 'Upload failed');
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Upload failed');
  }
};
</script>
