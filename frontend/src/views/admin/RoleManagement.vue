<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>Role & Permission Management</h2>
          <div class="header-actions">
            <el-button type="primary" @click="openCreateDialog">
              <el-icon><Plus /></el-icon>
              Create Role
            </el-button>
            <el-button @click="fetchRoles">Refresh</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <!-- Roles List -->
        <el-col :span="10">
          <el-card shadow="never" class="inner-card">
            <template #header>
              <div class="inner-card-header">
                <span>Roles</span>
                <el-tag type="info" size="small">{{ roles.length }} total</el-tag>
              </div>
            </template>
            <el-table 
              :data="roles" 
              style="width: 100%" 
              highlight-current-row 
              @current-change="handleRoleSelect"
              v-loading="loading"
              border
            >
              <el-table-column prop="name" label="Name" min-width="120" />
              <el-table-column prop="description" label="Description" min-width="150" show-overflow-tooltip />
              <el-table-column label="Actions" width="120" fixed="right">
                <template #default="scope">
                  <el-button link type="primary" size="small" @click.stop="openEditDialog(scope.row)">
                    Edit
                  </el-button>
                  <el-button link type="danger" size="small" @click.stop="handleDelete(scope.row)">
                    Delete
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="!loading && roles.length === 0" description="No roles found. Create one to get started." />
          </el-card>
        </el-col>

        <!-- Role Details & Permissions -->
        <el-col :span="14">
          <el-card v-if="selectedRole" shadow="never" class="inner-card">
            <template #header>
              <div class="inner-card-header">
                <span><strong>{{ selectedRole.Name }}</strong> - Permissions</span>
                <el-button type="primary" size="small" @click="savePermissions" :loading="saving">
                  Save Permissions
                </el-button>
              </div>
            </template>
            <div class="permissions-grid">
              <el-checkbox-group v-model="selectedPermissionIds">
                <el-row :gutter="10">
                  <el-col :span="8" v-for="perm in allPermissions" :key="perm.ID">
                    <el-checkbox :label="perm.ID" class="permission-item">
                      <div class="perm-info">
                        <span class="perm-name">{{ perm.Name }}</span>
                        <span class="perm-desc">{{ perm.Description }}</span>
                      </div>
                    </el-checkbox>
                  </el-col>
                </el-row>
              </el-checkbox-group>
              <el-empty v-if="allPermissions.length === 0" description="No permissions available" />
            </div>
          </el-card>
          <el-card v-else shadow="never" class="inner-card empty-card">
            <el-empty description="Select a role from the left to manage its permissions" />
          </el-card>
        </el-col>
      </el-row>
    </el-card>

    <!-- Create/Edit Role Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? 'Edit Role' : 'Create Role'"
      width="500px"
    >
      <el-form :model="roleForm" label-width="100px">
        <el-form-item label="Name" required>
          <el-input v-model="roleForm.name" placeholder="Role name" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="roleForm.description" type="textarea" :rows="3" placeholder="Role description" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? 'Update' : 'Create' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import {
  getRoles,
  getPermissions,
  createRole,
  updateRole,
  deleteRole,
  type RoleInfo,
  type PermissionInfo,
} from '../../api/role';

const roles = ref<RoleInfo[]>([]);
const allPermissions = ref<PermissionInfo[]>([]);
const selectedRole = ref<RoleInfo | null>(null);
const selectedPermissionIds = ref<number[]>([]);
const loading = ref(false);

const dialogVisible = ref(false);
const isEdit = ref(false);
const editingRoleId = ref<number | null>(null);
const submitting = ref(false);
const saving = ref(false);

const roleForm = reactive({
  name: '',
  description: '',
});

const fetchRoles = async () => {
  loading.value = true;
  try {
    const res = await getRoles();
    console.log('Roles response:', res);
    if (res.data?.roles) {
      roles.value = res.data.roles;
      console.log('Loaded roles:', roles.value);
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Failed to fetch roles');
  } finally {
    loading.value = false;
  }
};

const fetchPermissions = async () => {
  try {
    const res = await getPermissions();
    console.log('Permissions response:', res);
    if (res.data?.permissions) {
      allPermissions.value = res.data.permissions;
      console.log('Loaded permissions:', allPermissions.value);
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Failed to fetch permissions');
  }
};

const handleRoleSelect = (role: RoleInfo | null) => {
  selectedRole.value = role;
  if (role) {
    // Handle both Permissions and permissions (case sensitivity)
    const perms = (role as any).Permissions || (role as any).permissions || [];
    selectedPermissionIds.value = perms.map((p: any) => p.ID || p.id);
    console.log('Selected role permissions:', selectedPermissionIds.value);
  } else {
    selectedPermissionIds.value = [];
  }
};

const openCreateDialog = () => {
  isEdit.value = false;
  editingRoleId.value = null;
  roleForm.name = '';
  roleForm.description = '';
  dialogVisible.value = true;
};

const openEditDialog = (role: RoleInfo) => {
  isEdit.value = true;
  editingRoleId.value = role.ID;
  roleForm.name = role.Name;
  roleForm.description = role.Description || '';
  dialogVisible.value = true;
};

const handleSubmit = async () => {
  if (!roleForm.name.trim()) {
    ElMessage.warning('Role name is required');
    return;
  }

  submitting.value = true;
  try {
    if (isEdit.value && editingRoleId.value) {
      const res = await updateRole(editingRoleId.value, {
        name: roleForm.name,
        description: roleForm.description,
      });
      if (res.code === 200 || res.code === 0) {
        ElMessage.success('Role updated');
        dialogVisible.value = false;
        fetchRoles();
      } else {
        ElMessage.error(res.msg || 'Failed to update role');
      }
    } else {
      const res = await createRole({
        name: roleForm.name,
        description: roleForm.description,
      });
      if (res.code === 200 || res.code === 0) {
        ElMessage.success('Role created');
        dialogVisible.value = false;
        fetchRoles();
      } else {
        ElMessage.error(res.msg || 'Failed to create role');
      }
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.msg || 'Operation failed');
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (role: RoleInfo) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete role "${role.Name}"?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    const res = await deleteRole(role.ID);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Role deleted');
      if (selectedRole.value?.ID === role.ID) {
        selectedRole.value = null;
      }
      fetchRoles();
    } else {
      ElMessage.error(res.msg || 'Failed to delete role');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.response?.data?.msg || 'Failed to delete role');
    }
  }
};

const savePermissions = async () => {
  if (!selectedRole.value) return;

  saving.value = true;
  try {
    console.log('Saving permissions:', selectedPermissionIds.value);
    const res = await updateRole(selectedRole.value.ID, {
      permission_ids: selectedPermissionIds.value,
    });
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Permissions updated');
      // Refresh roles and re-select current role
      const currentRoleId = selectedRole.value.ID;
      await fetchRoles();
      // Re-select the role to refresh permissions display
      const updatedRole = roles.value.find(r => r.ID === currentRoleId);
      if (updatedRole) {
        handleRoleSelect(updatedRole);
      }
    } else {
      ElMessage.error(res.msg || 'Failed to update permissions');
    }
  } catch (err: any) {
    console.error('Save permissions error:', err);
    ElMessage.error(err.response?.data?.msg || 'Failed to save permissions');
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  fetchRoles();
  fetchPermissions();
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

.permissions-grid {
  padding: 10px 0;
}

.permission-item {
  margin-bottom: 15px;
  display: flex;
  align-items: flex-start;
}

.perm-info {
  display: flex;
  flex-direction: column;
}

.perm-name {
  font-weight: 500;
}

.perm-desc {
  font-size: 12px;
  color: #909399;
}

:deep(.el-table) {
  --el-table-border-color: #ebeef5;
}

:deep(.el-card__body) {
  padding: 15px;
}
</style>
