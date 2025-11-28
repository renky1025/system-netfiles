<template>
  <el-dialog
    v-model="dialogVisible"
    title="Version History"
    width="600px"
    destroy-on-close
  >
    <div v-loading="loading">
      <el-table :data="versions" style="width: 100%">
        <el-table-column prop="version" label="Ver" width="80">
          <template #default="scope">
            <el-tag>v{{ scope.row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="Size" width="100">
          <template #default="scope">
            {{ formatSize(scope.row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="Created At" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="Actions" min-width="150">
          <template #default="scope">
            <el-button 
              size="small" 
              type="warning" 
              @click="handleRollback(scope.row)"
            >
              Rollback
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(scope.row)"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-empty v-if="versions.length === 0 && !loading" description="No history versions found" />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { getFileVersions, rollbackVersion, deleteVersion } from '../api/version';
import { ElMessage, ElMessageBox } from 'element-plus';

const props = defineProps<{
  visible: boolean;
  fileId: number | null;
}>();

const emit = defineEmits(['update:visible', 'success']);

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const versions = ref<any[]>([]);
const loading = ref(false);

const fetchVersions = async () => {
  if (!props.fileId) return;
  
  loading.value = true;
  try {
    const res = await getFileVersions(props.fileId);
    if (res.data) {
      versions.value = res.data.versions || [];
    }
  } catch (err) {
    ElMessage.error('Failed to fetch version history');
  } finally {
    loading.value = false;
  }
};

const handleRollback = async (version: any) => {
  if (!props.fileId) return;

  try {
    await ElMessageBox.confirm(
      `Rollback to version v${version.version}? Current content will be overwritten.`,
      'Confirm Rollback',
      {
        confirmButtonText: 'Rollback',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    await rollbackVersion(props.fileId, version.id);
    ElMessage.success('File rolled back successfully');
    emit('success');
    dialogVisible.value = false;
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to rollback version');
    }
  }
};

const handleDelete = async (version: any) => {
  if (!props.fileId) return;

  try {
    await ElMessageBox.confirm(
      `Delete version v${version.version}?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'error',
      }
    );

    await deleteVersion(props.fileId, version.id);
    ElMessage.success('Version deleted successfully');
    fetchVersions();
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to delete version');
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

watch(() => props.visible, (val) => {
  if (val) {
    fetchVersions();
  }
});
</script>
