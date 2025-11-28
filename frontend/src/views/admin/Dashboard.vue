<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>System Dashboard</h2>
          <el-button @click="fetchData">Refresh</el-button>
        </div>
      </template>
    
      <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Total Users</span>
              <el-icon><User /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ systemStats.total_users }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Total Files</span>
              <el-icon><Document /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ systemStats.total_files }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Total Storage</span>
              <el-icon><Odometer /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ formatSize(systemStats.total_storage) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>Active Users (Today)</span>
              <el-icon><Timer /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ systemStats.active_users }}</div>
        </el-card>
      </el-col>
    </el-row>

      <el-row :gutter="20" class="chart-row">
        <el-col :span="12">
          <el-card shadow="hover">
            <template #header>
              <div class="stat-card-header">
                <span>Storage Usage</span>
              </div>
            </template>
            <div class="storage-info">
              <el-progress type="dashboard" :percentage="storageStats.usage_percent" :color="colors" />
              <div class="storage-details">
                <p>Used: {{ formatSize(storageStats.used_space) }}</p>
                <p>Free: {{ formatSize(storageStats.free_space) }}</p>
                <p>Total: {{ formatSize(storageStats.total_space) }}</p>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getSystemStats, getStorageStats } from '../../api/admin';
import { User, Document, Odometer, Timer } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

const systemStats = ref({
      active_shares: 0,
    active_users: 0,
    recycle_bin_files: 0,
    today_downloads: 0,
    today_uploads: 0,
    total_files: 0,
    total_storage: 0,
    total_users: 0
});

const storageStats = ref({
  total_space: 0,
  used_space: 0,
  free_space: 0,
  usage_percent: 0,
});

const colors = [
  { color: '#67C23A', percentage: 60 },
  { color: '#E6A23C', percentage: 80 },
  { color: '#F56C6C', percentage: 100 },
];

const fetchData = async () => {
  try {
    const sysRes = await getSystemStats();
    if (sysRes.data) {
      systemStats.value = sysRes.data;
    }

    const storageRes = await getStorageStats();
    if (storageRes.data) {
      storageStats.value = storageRes.data;
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Failed to fetch dashboard data');
  }
};

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

onMounted(() => {
  fetchData();
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

.stat-cards {
  margin-bottom: 20px;
}

.stat-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-value {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
  padding: 15px 0;
}

.chart-row {
  margin-top: 20px;
}

.storage-info {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 20px 0;
}

.storage-details p {
  margin: 8px 0;
  color: #606266;
}
</style>
