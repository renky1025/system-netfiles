<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>Organization Management</h2>
          <div class="header-actions">
            <el-button type="primary" @click="openCreateDialog(null)">Add Root Organization</el-button>
            <el-button @click="fetchOrgTree">Refresh</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="10">
          <el-card shadow="never" class="inner-card">
            <template #header>
              <div class="inner-card-header">
                <span>Organization Structure</span>
                <el-tag type="info" size="small">{{ flatOrgList.length }} total</el-tag>
              </div>
            </template>
            <el-tree
              ref="treeRef"
              :data="orgTree"
              :props="defaultProps"
              node-key="ID"
              highlight-current
              default-expand-all
              @node-click="handleNodeClick"
            >
              <template #default="{ node, data }">
                <span class="custom-tree-node">
                  <span>{{ node.label }}</span>
                  <span class="tree-actions">
                    <el-button link type="primary" size="small" @click.stop="openCreateDialog(data)">
                      Add Sub
                    </el-button>
                    <el-button link type="danger" size="small" @click.stop="handleDelete(data)">
                      Delete
                    </el-button>
                  </span>
                </span>
              </template>
            </el-tree>
          </el-card>
        </el-col>
        
        <el-col :span="14">
          <el-card v-if="selectedOrg" shadow="never" class="inner-card">
            <template #header>
              <div class="inner-card-header">
                <span><strong>{{ selectedOrg.Name }}</strong> - Details</span>
                <el-button type="primary" size="small" @click="handleUpdate">Save Changes</el-button>
              </div>
            </template>
            <el-form :model="editForm" label-width="120px">
              <el-form-item label="Name">
                <el-input v-model="editForm.name" />
              </el-form-item>
              <el-form-item label="Type">
                <el-select v-model="editForm.type" placeholder="Select type">
                  <el-option label="Company" value="company" />
                  <el-option label="Department" value="department" />
                  <el-option label="Team" value="team" />
                </el-select>
              </el-form-item>
              <el-form-item label="Manager ID">
                <el-input v-model.number="editForm.manager_id" type="number" placeholder="User ID" />
              </el-form-item>
              <el-divider content-position="left">Quota Settings</el-divider>
              <el-form-item label="Storage Quota">
                <el-input-number v-model="editForm.quotaGB" :min="0" :max="10000" />
                <span style="margin-left: 8px;">GB (0 = inherit)</span>
              </el-form-item>
              <el-form-item label="Download Limit">
                <el-input-number v-model="editForm.rateLimitMB" :min="0" :max="1000" />
                <span style="margin-left: 8px;">MB/s (0 = inherit)</span>
              </el-form-item>
              <el-form-item label="Created At">
                <span>{{ formatDate(selectedOrg.CreatedAt) }}</span>
              </el-form-item>
            </el-form>
          </el-card>
          <el-card v-else shadow="never" class="inner-card empty-card">
            <el-empty description="Select an organization to view details" />
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <!-- Create Dialog -->
    <el-dialog
      v-model="createDialogVisible"
      :title="createParent ? `Add Sub-Organization to ${createParent.Name}` : 'Add Root Organization'"
      width="500px"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="Name" required>
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item label="Type" required>
          <el-select v-model="createForm.type">
            <el-option label="Company" value="company" />
            <el-option label="Department" value="department" />
            <el-option label="Team" value="team" />
          </el-select>
        </el-form-item>
        <el-form-item label="Manager ID">
          <el-input v-model.number="createForm.manager_id" type="number" placeholder="Optional User ID" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="handleCreate">Create</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrgTree, createOrg, updateOrg, deleteOrg } from '../../api/org';
import type { OrgInfo } from '../../api/org';
import { setOrgQuota } from '../../api/quota';

const orgTree = ref<OrgInfo[]>([]);
const flatOrgList = ref<OrgInfo[]>([]);
const selectedOrg = ref<OrgInfo | null>(null);
const treeRef = ref();

const createDialogVisible = ref(false);
const createParent = ref<OrgInfo | null>(null);

const createForm = reactive({
  name: '',
  type: 'department',
  manager_id: undefined as number | undefined,
});

const editForm = reactive({
  name: '',
  type: '',
  manager_id: undefined as number | undefined,
  quotaGB: 0,
  rateLimitMB: 0,
});

const defaultProps = {
  children: 'children',
  label: 'Name',
};

const buildTree = (list: OrgInfo[]): OrgInfo[] => {
  const map = new Map<number, OrgInfo>();
  const roots: OrgInfo[] = [];

  // First pass: create map and initialize children
  list.forEach(item => {
    item.children = [];
    map.set(item.ID, item);
  });

  // Second pass: link children to parents
  list.forEach(item => {
    if (item.ParentID && map.has(item.ParentID)) {
      map.get(item.ParentID)!.children!.push(item);
    } else {
      roots.push(item);
    }
  });

  return roots;
};

const fetchOrgTree = async () => {
  try {
    const res = await getOrgTree();
    if (res.data && res.data.list) {
      flatOrgList.value = res.data.list;
      orgTree.value = buildTree(res.data.list);
    } else {
      flatOrgList.value = [];
      orgTree.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Failed to fetch organization tree');
  }
};

const handleNodeClick = (data: OrgInfo) => {
  selectedOrg.value = data;
  editForm.name = data.Name;
  editForm.type = data.Type;
  editForm.manager_id = data.ManagerID || undefined;
  editForm.quotaGB = (data as any).StorageQuota ? Math.round((data as any).StorageQuota / (1024 * 1024 * 1024)) : 0;
  editForm.rateLimitMB = (data as any).DownloadRateLimit ? Math.round((data as any).DownloadRateLimit / (1024 * 1024)) : 0;
};

const openCreateDialog = (parent: OrgInfo | null) => {
  createParent.value = parent;
  createForm.name = '';
  createForm.type = 'department';
  createForm.manager_id = undefined;
  createDialogVisible.value = true;
};

const handleCreate = async () => {
  if (!createForm.name) {
    ElMessage.warning('Name is required');
    return;
  }

  try {
    const res = await createOrg({
      name: createForm.name,
      type: createForm.type,
      parent_id: createParent.value ? createParent.value.ID : null,
      manager_id: createForm.manager_id || null,
    });

    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Organization created');
      createDialogVisible.value = false;
      fetchOrgTree();
    } else {
      ElMessage.error(res.msg || 'Failed to create');
    }
  } catch (err: any) {
    ElMessage.error('Failed to create organization');
  }
};

const handleUpdate = async () => {
  if (!selectedOrg.value) return;

  try {
    const res = await updateOrg(selectedOrg.value.ID, {
      name: editForm.name,
      type: editForm.type,
      manager_id: editForm.manager_id || null,
    });

    if (res.code === 200 || res.code === 0) {
      // Update quota
      const quotaBytes = editForm.quotaGB * 1024 * 1024 * 1024;
      const rateLimitBytes = editForm.rateLimitMB * 1024 * 1024;
      await setOrgQuota(selectedOrg.value.ID, quotaBytes, rateLimitBytes);
      
      ElMessage.success('Organization updated');
      fetchOrgTree();
    } else {
      ElMessage.error(res.msg || 'Failed to update');
    }
  } catch (err: any) {
    ElMessage.error('Failed to update organization');
  }
};

const handleDelete = async (data: OrgInfo) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete "${data.Name}"?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    const res = await deleteOrg(data.ID);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Organization deleted');
      if (selectedOrg.value && selectedOrg.value.ID === data.ID) {
        selectedOrg.value = null;
      }
      fetchOrgTree();
    } else {
      ElMessage.error(res.msg || 'Failed to delete');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.response?.data?.msg || 'Failed to delete organization');
    }
  }
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString();
};

onMounted(() => {
  fetchOrgTree();
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.inner-card {
  min-height: 400px;
}

.inner-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-card {
  display: flex;
  align-items: center;
  justify-content: center;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 8px;
}

.tree-actions {
  display: none;
}

.custom-tree-node:hover .tree-actions {
  display: inline-block;
}
</style>
