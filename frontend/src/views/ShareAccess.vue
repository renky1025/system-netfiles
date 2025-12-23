<template>
  <div class="share-page">
    <header class="share-header">
      <div class="share-logo">NetFileSys</div>
      <div class="share-title">Share</div>
    </header>

    <main class="share-main">
      <div class="share-card">
        <div class="header">
          <h2>NetFileSys Share</h2>
        </div>

      <div v-if="loading" class="loading">
        <el-skeleton :rows="3" animated />
      </div>

        <div v-else-if="error" class="error">
        <el-result
          icon="error"
          title="Access Failed"
          :sub-title="error"
        />
      </div>

        <div v-else-if="requirePassword && !authenticated" class="password-form">
        <el-result
          icon="warning"
          title="Password Protected"
          sub-title="Please enter password to access this file"
        >
          <template #extra>
            <el-form @submit.prevent="handlePasswordSubmit">
              <el-form-item>
                <el-input
                  v-model="password"
                  type="password"
                  placeholder="Enter Password"
                  show-password
                />
              </el-form-item>
              <el-button type="primary" @click="handlePasswordSubmit" :loading="validating">
                Access
              </el-button>
            </el-form>
          </template>
        </el-result>
      </div>

        <div v-else class="file-info">
          <div class="file-icon">
            <el-icon :size="64" color="#409eff"><Document /></el-icon>
          </div>
          <h3 class="file-name">{{ displayFileName }}</h3>
          <p class="file-meta">Size: {{ formatSize(displayFileSize) }}</p>
          <p class="file-meta" v-if="displayExpiredAt">Expires: {{ formatDate(displayExpiredAt) }}</p>
          
          <div class="actions">
            <el-button type="primary" size="large" @click="handleDownload">
              <el-icon><Download /></el-icon> Download
            </el-button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { getShareInfo, validateShare, downloadShareFile } from '../api/share';
import { Document, Download } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

const route = useRoute();
const code = route.params.code as string;

const loading = ref(true);
const error = ref('');
const requirePassword = ref(false);
const authenticated = ref(false);
const password = ref('');
const validating = ref(false);
const shareInfo = ref<any>({});

const displayFileName = computed(() => {
  const s = shareInfo.value || {};
  return s.file_name || s.FileName || s.Name || 'Unknown File';
});

const displayFileSize = computed(() => {
  const s = shareInfo.value || {};
  return s.file_size || s.FileSize || s.Size || 0;
});

const displayExpiredAt = computed(() => {
  const s = shareInfo.value || {};
  return s.expired_at || s.ExpiredAt || '';
});

const fetchShareInfo = async () => {
  try {
    const res = await getShareInfo(code);
    if (res.data?.share) {
      const data = res.data.share;
      shareInfo.value = data;
      
      // Check if password is required (Type 2)
      if (data.Type === 2 || data.type === 2) {
        requirePassword.value = true;
      } else {
        authenticated.value = true;
      }
    } else {
      error.value = res.msg || res.error || 'Share not found or expired';
    }
  } catch (err: any) {
    console.error(err);
    error.value = err.response?.data?.msg || err.response?.data?.error || 'Share not found or expired';
  } finally {
    loading.value = false;
  }
};

const handlePasswordSubmit = async () => {
  if (!password.value) return;
  
  validating.value = true;
  try {
    const res = await validateShare({
      share_id: shareInfo.value.ID || shareInfo.value.id,
      password: password.value,
    });
    if (res.code === 200 || res.code === 0) {
      authenticated.value = true;
    } else {
      ElMessage.error(res.msg || res.error || 'Invalid password');
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Invalid password');
  } finally {
    validating.value = false;
  }
};

const handleDownload = async () => {
  // TODO: Implement public download endpoint
  // Current backend implementation is missing a public download endpoint.
  // We will try to use a direct link if available or call a new endpoint.
  // For now, we'll simulate the call or use the authenticated one if user is logged in (which they aren't).
  
  try {
    const blob = await downloadShareFile(code, requirePassword.value ? password.value : undefined);
    
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', displayFileName.value || 'download');
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    
    ElMessage.success('Download started');
  } catch (err: any) {
    console.error(err);
    let message = 'Download failed';

    const resp = err?.response;
    if (resp?.data instanceof Blob) {
      try {
        const text = await resp.data.text();
        const json = JSON.parse(text);
        message = json.msg || json.error || message;
      } catch (parseErr) {
        // ignore parse error, keep default message
      }
    } else if (resp?.data) {
      message = resp.data.msg || resp.data.error || message;
    }

    ElMessage.error(message);
  }
};

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString();
};

onMounted(() => {
  fetchShareInfo();
});
</script>

<style scoped>
.share-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

.share-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #409eff;
  color: #fff;
  padding: 0 20px;
  height: 60px;
}

.share-logo {
  font-size: 20px;
  font-weight: bold;
}

.share-title {
  font-size: 16px;
  opacity: 0.9;
}

.share-main {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 40px 20px;
}

.share-card {
  width: 100%;
  max-width: 720px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
  text-align: center;
}

.header h2 {
  margin-top: 0;
  color: #303133;
}

.file-icon {
  margin: 20px 0;
}

.file-name {
  margin: 10px 0;
  font-size: 18px;
  color: #303133;
}

.file-meta {
  color: #909399;
  font-size: 14px;
  margin: 5px 0;
}

.actions {
  margin-top: 30px;
}
</style>
