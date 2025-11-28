<template>
  <div class="my-shares">
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

    <el-table :data="shares" style="width: 100%">
      <el-table-column prop="Code" label="Share Code" width="150" />
      <el-table-column label="Type" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.Type === 2 ? 'warning' : 'success'">
            {{ scope.row.Type === 2 ? 'Password' : 'Public' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ClickCount" label="Views" width="80" />
      <el-table-column label="Expires At" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.ExpiredAt) }}
        </template>
      </el-table-column>
      <el-table-column label="Created" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="Link" min-width="200">
        <template #default="scope">
          <el-input
            :value="getShareLink(scope.row.Code)"
            readonly
            size="small"
          >
            <template #append>
              <el-button @click="copyLink(scope.row.Code)">Copy</el-button>
            </template>
          </el-input>
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="120">
        <template #default="scope">
          <el-button size="small" type="danger" @click="deleteShare(scope.row)">
            Delete
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchShares"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getShareList, deleteShare as deleteShareApi } from '../api/share';
import { Refresh } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const shares = ref<any[]>([]);
const currentPage = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchShares = async () => {
  try {
    const res = await getShareList(currentPage.value, pageSize.value);
    if (res.data) {
      shares.value = res.data.data || [];
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

onMounted(() => {
  fetchShares();
});
</script>

<style scoped>
.my-shares {
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
