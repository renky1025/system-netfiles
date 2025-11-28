<template>
  <div class="file-preview-enhanced">
    <el-dialog
      v-model="dialogVisible"
      :title="fileName"
      :fullscreen="true"
      destroy-on-close
      class="preview-dialog"
    >
      <div class="preview-toolbar">
        <el-button-group>
          <el-button v-if="fileType === 'pdf'" @click="zoomIn" :icon="ZoomIn">放大</el-button>
          <el-button v-if="fileType === 'pdf'" @click="zoomOut" :icon="ZoomOut">缩小</el-button>
          <el-button @click="download" :icon="Download">下载</el-button>
        </el-button-group>
      </div>

      <div class="preview-content" v-loading="loading">
        <!-- Image -->
        <div v-if="fileType === 'image'" class="image-preview">
          <img :src="fileUrl" alt="Preview" />
        </div>

        <!-- Video -->
        <div v-if="fileType === 'video'" class="video-preview">
          <video controls autoplay>
            <source :src="fileUrl" />
            Your browser does not support the video tag.
          </video>
        </div>

        <!-- PDF -->
        <div v-if="fileType === 'pdf'" class="pdf-preview">
          <canvas ref="pdfCanvas" :style="{ transform: `scale(${zoomLevel})` }"></canvas>
          <div class="pdf-controls">
            <el-button @click="prevPage" :disabled="currentPage <= 1">上一页</el-button>
            <span>{{ currentPage }} / {{ totalPages }}</span>
            <el-button @click="nextPage" :disabled="currentPage >= totalPages">下一页</el-button>
          </div>
        </div>

        <!-- Office Documents (using Office Online Viewer) -->
        <div v-if="fileType === 'office'" class="office-preview">
          <iframe 
            :src="`https://view.officeapps.live.com/op/embed.aspx?src=${encodeURIComponent(fileUrl)}`"
            width="100%" 
            height="100%"
          ></iframe>
        </div>

        <!-- Text/Code -->
        <div v-if="fileType === 'text'" class="text-preview">
          <pre><code>{{ textContent }}</code></pre>
        </div>

        <!-- Unsupported -->
        <div v-if="fileType === 'unknown'" class="unknown-preview">
          <el-empty description="此文件类型暂不支持预览">
            <el-button type="primary" @click="download">下载文件</el-button>
          </el-empty>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue';
import { downloadFile } from '../api/file';
import { ElMessage } from 'element-plus';
import { Download, ZoomIn, ZoomOut } from '@element-plus/icons-vue';
import * as pdfjsLib from 'pdfjs-dist';

// Set PDF.js worker
pdfjsLib.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjsLib.version}/pdf.worker.min.js`;

const props = defineProps<{
  visible: boolean;
  fileId: number | null;
  fileName: string;
  fileSize: number;
}>();

const emit = defineEmits(['update:visible']);

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const fileUrl = ref('');
const textContent = ref('');
const loading = ref(false);
const pdfCanvas = ref<HTMLCanvasElement | null>(null);
const pdfDoc = ref<any>(null);
const currentPage = ref(1);
const totalPages = ref(0);
const zoomLevel = ref(1);

const fileType = computed(() => {
  const ext = props.fileName.split('.').pop()?.toLowerCase();
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp'].includes(ext || '')) return 'image';
  if (['mp4', 'webm', 'ogg', 'mov'].includes(ext || '')) return 'video';
  if (['pdf'].includes(ext || '')) return 'pdf';
  if (['doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'].includes(ext || '')) return 'office';
  if (['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'go', 'py', 'java', 'c', 'cpp', 'h'].includes(ext || '')) return 'text';
  return 'unknown';
});

const loadFile = async () => {
  if (!props.fileId) return;
  
  loading.value = true;
  try {
    const blob = await downloadFile(props.fileId);
    fileUrl.value = URL.createObjectURL(blob);

    if (fileType.value === 'text') {
      textContent.value = await blob.text();
    } else if (fileType.value === 'pdf') {
      await loadPDF(blob);
    }
  } catch (err) {
    ElMessage.error('文件预览加载失败');
  } finally {
    loading.value = false;
  }
};

const loadPDF = async (blob: Blob) => {
  const arrayBuffer = await blob.arrayBuffer();
  pdfDoc.value = await pdfjsLib.getDocument({ data: arrayBuffer }).promise;
  totalPages.value = pdfDoc.value.numPages;
  currentPage.value = 1;
  await renderPage(1);
};

const renderPage = async (pageNum: number) => {
  if (!pdfDoc.value || !pdfCanvas.value) return;
  
  const page = await pdfDoc.value.getPage(pageNum);
  const viewport = page.getViewport({ scale: 1.5 });
  
  const canvas = pdfCanvas.value;
  const context = canvas.getContext('2d');
  canvas.height = viewport.height;
  canvas.width = viewport.width;
  
  await page.render({
    canvasContext: context,
    viewport: viewport
  }).promise;
};

const prevPage = async () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    await renderPage(currentPage.value);
  }
};

const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    await renderPage(currentPage.value);
  }
};

const zoomIn = () => {
  zoomLevel.value = Math.min(zoomLevel.value + 0.1, 3);
};

const zoomOut = () => {
  zoomLevel.value = Math.max(zoomLevel.value - 0.1, 0.5);
};

const download = () => {
  const link = document.createElement('a');
  link.href = fileUrl.value;
  link.download = props.fileName;
  link.click();
};

watch(() => props.visible, (val) => {
  if (val) {
    nextTick(() => loadFile());
  } else {
    // Cleanup
    if (fileUrl.value) {
      URL.revokeObjectURL(fileUrl.value);
      fileUrl.value = '';
    }
    textContent.value = '';
    pdfDoc.value = null;
    currentPage.value = 1;
    totalPages.value = 0;
    zoomLevel.value = 1;
  }
});
</script>

<style scoped>
.preview-toolbar {
  padding: 10px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-content {
  height: calc(100vh - 150px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  overflow: auto;
}

.image-preview {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.image-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.video-preview {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.video-preview video {
  max-width: 100%;
  max-height: 100%;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.pdf-preview {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  overflow: auto;
}

.pdf-preview canvas {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  background: white;
  margin-bottom: 20px;
  transition: transform 0.3s ease;
}

.pdf-controls {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.office-preview {
  width: 100%;
  height: 100%;
}

.text-preview {
  width: 100%;
  height: 100%;
  padding: 20px;
  background-color: white;
  overflow: auto;
}

.text-preview pre {
  margin: 0;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.unknown-preview {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
