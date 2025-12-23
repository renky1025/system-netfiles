<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>Audit Logs</h2>
        </div>
      </template>
    
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <el-tab-pane label="File Operations" name="file-ops">
        <div class="tab-content">
          <div class="toolbar">
            <el-button @click="fetchFileLogs">Refresh</el-button>
          </div>
          
          <el-table :data="fileLogs" style="width: 100%" v-loading="fileLoading">
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="User" width="150">
              <template #default="scope">
                {{ scope.row.User?.username || scope.row.User?.Username || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="OpType" label="Operation" width="140">
              <template #default="scope">
                <el-tag>{{ scope.row.OpType }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="File ID" width="120">
              <template #default="scope">
                <span v-if="scope.row.FileID">{{ scope.row.FileID }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="Details" label="Details" min-width="200" show-overflow-tooltip />
            <el-table-column prop="ClientIP" label="IP" width="140" />
            <el-table-column label="Time" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.CreatedAt) }}
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="filePage"
            v-model:page-size="filePageSize"
            :total="fileTotal"
            layout="total, prev, pager, next"
            @current-change="fetchFileLogs"
            class="pagination"
          />
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="Login Logs" name="login-logs">
        <div class="tab-content">
          <div class="toolbar">
            <el-button @click="fetchLoginLogs">Refresh</el-button>
          </div>

          <el-table :data="loginLogs" style="width: 100%" v-loading="loginLoading">
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="User" width="150">
              <template #default="scope">
                {{ scope.row.User?.username || scope.row.User?.Username || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="Status" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.Success ? 'success' : 'danger'">
                  {{ scope.row.Success ? 'Success' : 'Failed' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ClientIP" label="IP" width="140" />
            <el-table-column prop="UserAgent" label="User Agent" min-width="200" show-overflow-tooltip />
            <el-table-column label="Time" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.CreatedAt) }}
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="loginPage"
            v-model:page-size="loginPageSize"
            :total="loginTotal"
            layout="total, prev, pager, next"
            @current-change="fetchLoginLogs"
            class="pagination"
          />
        </div>
      </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getFileOpLogs, getLoginLogs } from '../../api/admin';
import { ElMessage } from 'element-plus';

const activeTab = ref('file-ops');

// File Logs State
const fileLogs = ref<any[]>([]);
const fileLoading = ref(false);
const filePage = ref(1);
const filePageSize = ref(20);
const fileTotal = ref(0);

// Login Logs State
const loginLogs = ref<any[]>([]);
const loginLoading = ref(false);
const loginPage = ref(1);
const loginPageSize = ref(20);
const loginTotal = ref(0);

const fetchFileLogs = async () => {
  fileLoading.value = true;
  try {
    const res = await getFileOpLogs(filePage.value, filePageSize.value);
    if (res.data) {
      fileLogs.value = res.data.list || [];
      fileTotal.value = res.data.total || 0;
    }
  } catch (err) {
    ElMessage.error('Failed to fetch file logs');
  } finally {
    fileLoading.value = false;
  }
};

const fetchLoginLogs = async () => {
  loginLoading.value = true;
  try {
    const res = await getLoginLogs(loginPage.value, loginPageSize.value);
    if (res.data) {
      loginLogs.value = res.data.list || [];
      loginTotal.value = res.data.total || 0;
    }
  } catch (err) {
    ElMessage.error('Failed to fetch login logs');
  } finally {
    loginLoading.value = false;
  }
};

const handleTabClick = (tab: any) => {
  if (tab.props.name === 'file-ops' && fileLogs.value.length === 0) {
    fetchFileLogs();
  } else if (tab.props.name === 'login-logs' && loginLogs.value.length === 0) {
    fetchLoginLogs();
  }
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString();
};

onMounted(() => {
  fetchFileLogs();
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

.toolbar {
  margin-bottom: 15px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
