<template>
  <div class="recycle-bin">
    <el-page-header @back="$router.back()" title="Back">
      <template #content>
        <h2>Recycle Bin</h2>
      </template>
    </el-page-header>

    <div class="toolbar">
      <el-button @click="clearAll" type="danger" :disabled="files.length === 0">
        <el-icon><Delete /></el-icon>
        Clear All
      </el-button>
      <el-button @click="fetchFiles">
        <el-icon><Refresh /></el-icon>
        Refresh
      </el-button>
    </div>

    <el-table :data="files" style="width: 100%">
      <el-table-column prop="Name" label="Name" width="300" />
      <el-table-column prop="Size" label="Size" width="120">
        <template #default="scope">
          {{ formatSize(scope.row.Size) }}
        </template>
      </el-table-column>
      <el-table-column prop="DeletedAt" label="Deleted At" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.DeletedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="200">
        <template #default="scope">
          <el-button size="small" type="success" @click="restoreFile(scope.row)">
            Restore
          </el-button>
          <el-button size="small" type="danger" @click="permanentDelete(scope.row)">
            Delete Forever
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchFiles"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getRecycleList, restoreFile as restoreFileApi, permanentDeleteFile, clearRecycleBin } from '../api/recycle';
import { Delete, Refresh } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const files = ref<any[]>([]);
const currentPage = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchFiles = async () => {
  try {
    const res = await getRecycleList(currentPage.value, pageSize.value);
    if (res.data) {
      files.value = res.data.data || [];
      total.value = res.data.total || 0;
    } else {
      files.value = [];
      total.value = 0;
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch recycle bin');
  }
};

const restoreFile = async (file: any) => {
  try {
    await ElMessageBox.confirm(
      `Restore "${file.Name}"?`,
      'Confirm Restore',
      {
        confirmButtonText: 'Restore',
        cancelButtonText: 'Cancel',
        type: 'info',
      }
    );

    const res = await restoreFileApi(file.ID);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('File restored successfully');
      fetchFiles();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to restore file');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to restore file');
    }
  }
};

const permanentDelete = async (file: any) => {
  try {
    await ElMessageBox.confirm(
      `Permanently delete "${file.Name}"? This action cannot be undone!`,
      'Confirm Permanent Delete',
      {
        confirmButtonText: 'Delete Forever',
        cancelButtonText: 'Cancel',
        type: 'error',
      }
    );

    const res = await permanentDeleteFile(file.ID);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('File permanently deleted');
      fetchFiles();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to delete file');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to delete file');
    }
  }
};

const clearAll = async () => {
  try {
    await ElMessageBox.confirm(
      'Clear all files in recycle bin? This action cannot be undone!',
      'Confirm Clear All',
      {
        confirmButtonText: 'Clear All',
        cancelButtonText: 'Cancel',
        type: 'error',
      }
    );

    const res = await clearRecycleBin();
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Recycle bin cleared');
      fetchFiles();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to clear recycle bin');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to clear recycle bin');
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
.recycle-bin {
  padding: 20px;
}

.toolbar {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: center;
}
</style>
