<template>
  <el-card class="storage-quota-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <el-icon><Coin /></el-icon>
        <span>Storage Usage</span>
      </div>
    </template>
    
    <div class="quota-content" v-loading="loading">
      <el-progress 
        type="dashboard" 
        :percentage="quotaInfo.usage_percent" 
        :color="getProgressColor"
        :stroke-width="12"
      >
        <template #default>
          <span class="percentage-value">{{ quotaInfo.usage_percent.toFixed(1) }}%</span>
        </template>
      </el-progress>
      
      <div class="quota-details">
        <div class="detail-item">
          <span class="label">Used:</span>
          <span class="value">{{ formatSize(quotaInfo.used_storage) }}</span>
        </div>
        <div class="detail-item">
          <span class="label">Free:</span>
          <span class="value">{{ formatSize(quotaInfo.free_storage) }}</span>
        </div>
        <div class="detail-item">
          <span class="label">Total:</span>
          <span class="value">{{ formatSize(quotaInfo.storage_quota) }}</span>
        </div>
        <div class="detail-item source">
          <span class="label">Source:</span>
          <el-tag size="small" :type="getSourceTagType(quotaInfo.quota_source)">
            {{ quotaInfo.quota_source }}
          </el-tag>
        </div>
      </div>

      <div v-if="quotaInfo.download_rate_limit > 0" class="rate-limit-info">
        <el-icon><Download /></el-icon>
        <span>Download limit: {{ formatSize(quotaInfo.download_rate_limit) }}/s</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { Coin, Download } from '@element-plus/icons-vue';
import { getMyQuota, type QuotaInfo } from '../api/quota';

const loading = ref(false);
const quotaInfo = ref<QuotaInfo>({
  user_id: 0,
  storage_quota: 0,
  used_storage: 0,
  free_storage: 0,
  usage_percent: 0,
  quota_source: 'system',
  rate_limit_source: 'system',
  download_rate_limit: 0
});

const getProgressColor = computed(() => {
  const percent = quotaInfo.value.usage_percent;
  if (percent >= 90) return '#F56C6C';
  if (percent >= 70) return '#E6A23C';
  return '#409EFF';
});

const getSourceTagType = (source: string) => {
  switch (source) {
    case 'user': return 'success';
    case 'role': return 'primary';
    case 'organization': return 'warning';
    case 'admin': return 'danger';
    default: return 'info';
  }
};

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const fetchQuota = async () => {
  loading.value = true;
  try {
    const res = await getMyQuota();
    if (res.data) {
      quotaInfo.value = res.data;
    }
  } catch (err) {
    console.error('Failed to fetch quota:', err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchQuota();
});

// Expose refresh method
defineExpose({ refresh: fetchQuota });
</script>

<style scoped>
.storage-quota-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.quota-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.percentage-value {
  font-size: 20px;
  font-weight: bold;
}

.quota-details {
  width: 100%;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
}

.detail-item.source {
  grid-column: 1 / -1;
  border-top: 1px solid #dcdfe6;
  padding-top: 8px;
  margin-top: 4px;
}

.label {
  color: #909399;
  font-size: 13px;
}

.value {
  font-weight: 500;
  color: #303133;
}

.rate-limit-info {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 12px;
  padding: 8px 12px;
  background: #fdf6ec;
  border-radius: 4px;
}
</style>
