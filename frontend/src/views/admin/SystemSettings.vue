<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>System Settings</h2>
        </div>
      </template>

      <el-row :gutter="20">
      <!-- Storage Statistics -->
      <el-col :span="12">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>Storage Statistics</span>
              <el-button link type="primary" @click="fetchStorageStats">Refresh</el-button>
            </div>
          </template>
          <div class="storage-stats" v-loading="loadingStorage">
            <el-progress
              type="dashboard"
              :percentage="storageStats.usage_percent"
              :color="getStorageColor(storageStats.usage_percent)"
            >
              <template #default>
                <span class="percentage-value">{{ storageStats.usage_percent }}%</span>
                <span class="percentage-label">Used</span>
              </template>
            </el-progress>
            <div class="storage-details">
              <div class="stat-item">
                <span class="label">Total Space:</span>
                <span class="value">{{ formatSize(storageStats.total_space) }}</span>
              </div>
              <div class="stat-item">
                <span class="label">Used Space:</span>
                <span class="value">{{ formatSize(storageStats.used_space) }}</span>
              </div>
              <div class="stat-item">
                <span class="label">Free Space:</span>
                <span class="value">{{ formatSize(storageStats.free_space) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- System Statistics -->
      <el-col :span="12">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>System Statistics</span>
              <el-button link type="primary" @click="fetchSystemStats">Refresh</el-button>
            </div>
          </template>
          <div class="system-stats" v-loading="loadingSystem">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-statistic title="Total Users" :value="systemStats.user_count">
                  <template #prefix>
                    <el-icon><User /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="12">
                <el-statistic title="Total Files" :value="systemStats.file_count">
                  <template #prefix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
            </el-row>
            <el-row :gutter="20" style="margin-top: 20px;">
              <el-col :span="12">
                <el-statistic title="Active Users Today" :value="systemStats.active_users_today">
                  <template #prefix>
                    <el-icon><UserFilled /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="12">
                <el-statistic title="Total Storage" :value="formatSize(systemStats.total_storage)">
                  <template #prefix>
                    <el-icon><FolderOpened /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Configuration Settings -->
    <el-card class="config-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>Configuration</span>
        </div>
      </template>
      <el-form label-width="200px" class="config-form">
        <el-form-item label="Storage Type">
          <el-tag>{{ configInfo.storage_type || 'Local' }}</el-tag>
        </el-form-item>
        <el-form-item label="Max Upload Size">
          <span>{{ formatSize(configInfo.max_upload_size || 104857600) }}</span>
        </el-form-item>
        <el-form-item label="Allowed File Types">
          <el-tag v-for="type in (configInfo.allowed_types || ['*'])" :key="type" style="margin-right: 5px;">
            {{ type }}
          </el-tag>
        </el-form-item>
        <el-form-item label="Share Expiry Default">
          <span>{{ configInfo.share_expiry_days || 7 }} days</span>
        </el-form-item>
        <el-form-item label="Recycle Bin Retention">
          <span>{{ configInfo.recycle_retention_days || 30 }} days</span>
        </el-form-item>
      </el-form>
      <el-alert
        title="Configuration changes require server restart"
        type="info"
        :closable="false"
        show-icon
        style="margin-top: 20px;"
      />
    </el-card>

    <!-- Quick Actions -->
    <el-card class="actions-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>Quick Actions</span>
        </div>
      </template>
      <div class="actions-grid">
        <el-button type="primary" @click="clearRecycleBin" :loading="clearingRecycle">
          <el-icon><Delete /></el-icon>
          Clear All Recycle Bins
        </el-button>
        <el-button type="warning" @click="refreshCache">
          <el-icon><Refresh /></el-icon>
          Refresh Cache
        </el-button>
      </div>
    </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { User, Document, UserFilled, FolderOpened, Delete, Refresh } from '@element-plus/icons-vue';
import { getSystemStats, getStorageStats } from '../../api/admin';

const loadingStorage = ref(false);
const loadingSystem = ref(false);
const clearingRecycle = ref(false);

const storageStats = reactive({
  total_space: 0,
  used_space: 0,
  free_space: 0,
  usage_percent: 0,
});

const systemStats = reactive({
  user_count: 0,
  file_count: 0,
  total_storage: 0,
  active_users_today: 0,
});

const configInfo = reactive({
  storage_type: 'Local',
  max_upload_size: 104857600,
  allowed_types: ['*'],
  share_expiry_days: 7,
  recycle_retention_days: 30,
});

const fetchStorageStats = async () => {
  loadingStorage.value = true;
  try {
    const res = await getStorageStats();
    if (res.data) {
      Object.assign(storageStats, res.data);
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Failed to fetch storage stats');
  } finally {
    loadingStorage.value = false;
  }
};

const fetchSystemStats = async () => {
  loadingSystem.value = true;
  try {
    const res = await getSystemStats();
    if (res.data) {
      Object.assign(systemStats, res.data);
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Failed to fetch system stats');
  } finally {
    loadingSystem.value = false;
  }
};

const getStorageColor = (percent: number) => {
  if (percent < 60) return '#67c23a';
  if (percent < 80) return '#e6a23c';
  return '#f56c6c';
};

const formatSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const clearRecycleBin = async () => {
  try {
    await ElMessageBox.confirm(
      'This will permanently delete all files in all users\' recycle bins. Continue?',
      'Warning',
      {
        confirmButtonText: 'Clear All',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );
    
    clearingRecycle.value = true;
    // Note: This would need a backend endpoint to clear all recycle bins
    ElMessage.info('This feature requires backend implementation');
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('Operation failed');
    }
  } finally {
    clearingRecycle.value = false;
  }
};

const refreshCache = () => {
  ElMessage.success('Cache refresh requested');
};

onMounted(() => {
  fetchStorageStats();
  fetchSystemStats();
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

.inner-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-card {
  min-height: 300px;
}

.storage-stats {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.percentage-value {
  font-size: 24px;
  font-weight: bold;
}

.percentage-label {
  font-size: 12px;
  color: #909399;
}

.storage-details {
  width: 100%;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #ebeef5;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-item .label {
  color: #606266;
}

.stat-item .value {
  font-weight: 500;
}

.system-stats {
  padding: 20px 0;
}

.config-form {
  max-width: 600px;
}

.actions-grid {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}
</style>
