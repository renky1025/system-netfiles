<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>File Management</h2>
          <el-button @click="fetchFiles">Refresh</el-button>
        </div>
      </template>

      <el-table :data="files" style="width: 100%" v-loading="loading" border>
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Name" label="Name" min-width="200" show-overflow-tooltip>
        <template #default="scope">
          <div class="file-name">
            <el-icon><Document /></el-icon>
            <span>{{ scope.row.Name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="Size" label="Size" width="120">
        <template #default="scope">
          {{ formatSize(scope.row.Size) }}
        </template>
      </el-table-column>
      <el-table-column prop="MD5" label="MD5" width="200" show-overflow-tooltip />
      <el-table-column label="Uploaded At" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="150">
        <template #default="scope">
          <el-button size="small" type="danger" @click="handleForceDelete(scope.row)">
            Force Delete
          </el-button>
        </template>
      </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchFiles"
          @size-change="fetchFiles"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getAllFiles, forceDeleteFile } from '../../api/admin';
import { Document } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const files = ref<any[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchFiles = async () => {
  loading.value = true;
  try {
    const res = await getAllFiles(currentPage.value, pageSize.value);
    if (res.data) {
      files.value = res.data.list || [];
      // Note: The API response structure for file list might need adjustment if total is not returned directly
      // Assuming standard pagination response structure here
      total.value = (res.data as any).total || 0; 
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Failed to fetch files');
  } finally {
    loading.value = false;
  }
};

const handleForceDelete = async (file: any) => {
  try {
    await ElMessageBox.confirm(
      `Permanently delete file "${file.Name}"? This cannot be undone and will affect all users!`,
      'Confirm Force Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'error',
      }
    );

    await forceDeleteFile(file.ID);
    ElMessage.success('File force deleted successfully');
    fetchFiles();
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to delete file');
    }
  }
};

const formatSize = (bytes: number) => {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString();
};

onMounted(() => {
  fetchFiles();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
    min-width: 1024px;
}

.page-card {
  min-height: calc(100vh - 140px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
