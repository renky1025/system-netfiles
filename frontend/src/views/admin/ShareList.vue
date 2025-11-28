<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>Share Management</h2>
          <el-button @click="fetchShares">Refresh</el-button>
        </div>
      </template>

      <el-table :data="shares" style="width: 100%" v-loading="loading" border>
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Code" label="Code" width="120" />
      <el-table-column prop="file_name" label="File" min-width="200" show-overflow-tooltip />
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
      <el-table-column label="Created At" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="150">
        <template #default="scope">
          <el-button size="small" type="danger" @click="handleDisable(scope.row)">
            Disable
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
          @current-change="fetchShares"
          @size-change="fetchShares"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getAllShares, disableShare } from '../../api/admin';
import { ElMessage, ElMessageBox } from 'element-plus';

const shares = ref<any[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(20);
const total = ref(0);

const fetchShares = async () => {
  loading.value = true;
  try {
    const res = await getAllShares(currentPage.value, pageSize.value);
    if (res.data) {
      shares.value = res.data.list || [];
      total.value = res.data.total || 0;
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Failed to fetch shares');
  } finally {
    loading.value = false;
  }
};

const handleDisable = async (share: any) => {
  try {
    await ElMessageBox.confirm(
      `Disable share "${share.Code}"? It will no longer be accessible.`,
      'Confirm Disable',
      {
        confirmButtonText: 'Disable',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    await disableShare(share.ID);
    ElMessage.success('Share disabled successfully');
    fetchShares();
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to disable share');
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

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
