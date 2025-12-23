<template>
  <div class="home-container">
    <el-header class="header">
      <div class="logo">NetFileSys</div>
      <div class="nav-menu">
        <el-button text @click="$router.push('/recycle')">
          <el-icon><Delete /></el-icon>
          Recycle Bin
        </el-button>
        <el-button text @click="$router.push('/shares')">
          <el-icon><Share /></el-icon>
          My Shares
        </el-button>
      </div>
      <div class="user-info">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link" style="color: white; cursor: pointer; display: flex; align-items: center;">
            User <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="admin">Admin Portal</el-dropdown-item>
              <el-dropdown-item command="change-password">Change Password</el-dropdown-item>
              <el-dropdown-item command="logout">Logout</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-container>
      <el-aside width="250px" class="sidebar">
        <div class="sidebar-header">
          <h3>Folders</h3>
        </div>
        <FolderTree @folder-selected="handleFolderSelected" />
        <div class="sidebar-quota">
          <StorageQuotaCard />
        </div>
      </el-aside>

      <el-main>
        <div class="page-container">
          <el-page-header @back="$router.back()" title="Back">
            <template #content>
              <h2>My Shares</h2>
            </template>
          </el-page-header>

          <div class="toolbar">
            <el-button @click="fetchShares">
              <el-icon><Refresh /></el-icon>
              Refresh
            </el-button>
          </div>

          <div class="table-card">
            <el-table
              :data="shares"
              style="width: 100%"
              border
              stripe
              class="data-table"
            >
              <el-table-column prop="Code" label="Share Code" width="180" />
              <el-table-column label="Type" width="120" align="center">
                <template #default="scope">
                  <el-tag :type="scope.row.Type === 2 ? 'warning' : 'success'">
                    {{ scope.row.Type === 2 ? 'Password' : 'Public' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="ClickCount" label="Views" width="100" align="center" />
              <el-table-column label="Expires At" width="200">
                <template #default="scope">
                  {{ formatDate(scope.row.ExpiredAt) }}
                </template>
              </el-table-column>
              <el-table-column label="Created" width="200">
                <template #default="scope">
                  {{ formatDate(scope.row.CreatedAt) }}
                </template>
              </el-table-column>
              <el-table-column label="Link" min-width="260">
                <template #default="scope">
                  <el-input
                    :value="getShareLink(scope.row.Code)"
                    readonly
                    size="small"
                  >
                    <template #append>
                      <el-button size="small" type="primary" @click="copyLink(scope.row.Code)">
                        Copy
                      </el-button>
                    </template>
                  </el-input>
                </template>
              </el-table-column>
              <el-table-column label="Actions" width="140" align="center">
                <template #default="scope">
                  <el-button size="small" type="danger" @click="deleteShare(scope.row)">
                    Delete
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <div class="pagination-wrapper">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="total"
                layout="total, prev, pager, next"
                @current-change="fetchShares"
              />
            </div>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getShareList, deleteShare as deleteShareApi } from '../api/share';
import { Refresh, Delete, Share, ArrowDown } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useUserStore } from '../store/user';
import FolderTree from '../components/FolderTree.vue';
import StorageQuotaCard from '../components/StorageQuotaCard.vue';

const router = useRouter();
const userStore = useUserStore();

const shares = ref<any[]>([]);
const currentPage = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchShares = async () => {
  try {
    const res = await getShareList(currentPage.value, pageSize.value);
    if (res.data) {
      shares.value = res.data.list || [];
      total.value = res.data.total || 0;
    } else {
      shares.value = [];
      total.value = 0;
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch shares');
  }
};

const getShareLink = (code: string) => {
  return `${window.location.origin}/share/${code}`;
};

const copyLink = async (code: string) => {
  try {
    await navigator.clipboard.writeText(getShareLink(code));
    ElMessage.success('Link copied to clipboard');
  } catch (err) {
    ElMessage.error('Failed to copy link');
  }
};

const deleteShare = async (share: any) => {
  try {
    await ElMessageBox.confirm(
      `Delete share "${share.Code}"?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    const res = await deleteShareApi(share.ID);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Share deleted successfully');
      fetchShares();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to delete share');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to delete share');
    }
  }
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'Never';
  return new Date(dateStr).toLocaleString();
};

const handleFolderSelected = (folderId: number | null) => {
  if (folderId) {
    router.push(`/folder/${folderId}`);
  } else {
    router.push('/');
  }
};

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout();
    router.push('/login');
  } else if (command === 'change-password') {
    router.push('/change-password');
  } else if (command === 'admin') {
    router.push('/admin');
  }
};

onMounted(() => {
  fetchShares();
});
</script>

<style scoped>
.home-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #409eff;
  color: white;
  padding: 0 20px;
  height: 60px;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  min-width: 150px;
}

.nav-menu {
  display: flex;
  gap: 10px;
  flex: 1;
  margin-left: 40px;
}

.nav-menu .el-button {
  color: white;
}

.user-info {
  display: flex;
  align-items: center;
}

.page-container {
  padding: 20px;
}

.sidebar {
  background-color: #f5f7fa;
  padding: 20px 10px;
  border-right: 1px solid #dcdfe6;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 16px;
}

.sidebar-quota {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #dcdfe6;
}

.toolbar {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.table-card {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.data-table :deep(.el-table__header-wrapper) {
  background-color: #f5f7fa;
}

.data-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}
</style>
