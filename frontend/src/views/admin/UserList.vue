<template>
  <div class="page-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <h2>User Management</h2>
          <div class="header-actions">
            <el-input
              v-model="searchQuery"
              placeholder="Search users..."
              clearable
              style="width: 200px; margin-right: 10px;"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
            <el-button type="primary" @click="openCreateDialog">
              <el-icon><Plus /></el-icon>
              Add User
            </el-button>
            <el-button @click="fetchUsers">Refresh</el-button>
          </div>
        </div>
      </template>

      <el-table :data="users" style="width: 100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="Username" min-width="120" />
        <el-table-column prop="email" label="Email" min-width="180" />
        <el-table-column label="Roles" min-width="150">
          <template #default="scope">
            <el-tag 
              v-for="role in (scope.row.roles || [])" 
              :key="role.id" 
              size="small" 
              style="margin-right: 4px;"
            >
              {{ role.name }}
            </el-tag>
            <span v-if="!scope.row.roles?.length" class="text-muted">No roles</span>
          </template>
        </el-table-column>
        <el-table-column label="Organizations" min-width="150">
          <template #default="scope">
            <el-tag 
              v-for="org in (scope.row.organizations || [])" 
              :key="org.id" 
              size="small" 
              type="warning"
              style="margin-right: 4px;"
            >
              {{ org.name }}
              <el-icon v-if="org.is_primary" style="margin-left: 2px;"><Star /></el-icon>
            </el-tag>
            <span v-if="!scope.row.organizations?.length" class="text-muted">No org</span>
          </template>
        </el-table-column>
        <el-table-column label="Status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 1 ? 'Active' : 'Frozen' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Created At" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="320" fixed="right">
          <template #default="scope">
            <el-button-group>
              <el-button size="small" type="primary" @click="openEditDialog(scope.row)" title="Edit">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button size="small" @click="openAssignRoleDialog(scope.row)" title="Assign Role">
                <el-icon><User /></el-icon>
              </el-button>
              <el-button size="small" @click="openAssignOrgDialog(scope.row)" title="Assign Org">
                <el-icon><OfficeBuilding /></el-icon>
              </el-button>
              <el-button 
                size="small" 
                :type="scope.row.status === 1 ? 'warning' : 'success'"
                @click="handleStatusChange(scope.row)"
                :title="scope.row.status === 1 ? 'Freeze' : 'Unfreeze'"
              >
                <el-icon><Lock v-if="scope.row.status === 1" /><Unlock v-else /></el-icon>
              </el-button>
              <el-button size="small" @click="handleResetPassword(scope.row)" title="Reset Password">
                <el-icon><Key /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row)" title="Delete">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
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
          @current-change="fetchUsers"
          @size-change="fetchUsers"
        />
      </div>
    </el-card>

    <!-- Create/Edit User Dialog -->
    <el-dialog v-model="userDialogVisible" :title="isEditUser ? 'Edit User' : 'Create User'" width="500px">
      <el-form :model="userForm" label-width="100px" :rules="userRules" ref="userFormRef">
        <el-form-item label="Username" prop="username">
          <el-input v-model="userForm.username" placeholder="Enter username" :disabled="isEditUser" />
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model="userForm.email" placeholder="Enter email" />
        </el-form-item>
        <el-form-item v-if="!isEditUser" label="Password" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="Enter password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleUserSubmit" :loading="submitting">
          {{ isEditUser ? 'Update' : 'Create' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Assign Role Dialog -->
    <el-dialog v-model="assignRoleDialogVisible" title="Assign Role" width="500px">
      <div v-if="selectedUser">
        <p>User: <strong>{{ selectedUser.username }}</strong></p>
        <p>Current Roles: 
          <el-tag v-for="role in (selectedUser.roles || [])" :key="role.id" size="small" style="margin-right: 4px;">
            {{ role.name }}
          </el-tag>
          <span v-if="!selectedUser.roles?.length">None</span>
        </p>
        <el-divider />
        <el-form label-width="100px">
          <el-form-item label="Add Role">
            <el-select v-model="selectedRoleId" placeholder="Select a role" style="width: 100%;">
              <el-option 
                v-for="role in availableRoles" 
                :key="role.id" 
                :label="role.name" 
                :value="role.id"
                :disabled="userHasRole(role.id)"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="assignRoleDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleAssignRole" :disabled="!selectedRoleId">
          Assign
        </el-button>
      </template>
    </el-dialog>

    <!-- Assign Organization Dialog -->
    <el-dialog v-model="assignOrgDialogVisible" title="Assign Organization" width="500px">
      <div v-if="selectedUser">
        <p>User: <strong>{{ selectedUser.username }}</strong></p>
        <p>Current Organizations: 
          <el-tag v-for="org in (selectedUser.organizations || [])" :key="org.id" size="small" type="warning" style="margin-right: 4px;">
            {{ org.name }}
          </el-tag>
          <span v-if="!selectedUser.organizations?.length">None</span>
        </p>
        <el-divider />
        <el-form label-width="120px">
          <el-form-item label="Organization">
            <el-tree-select
              v-model="selectedOrgId"
              :data="filteredOrgTree"
              :props="{ label: 'name', value: 'id', children: 'children' }"
              placeholder="Select organization"
              style="width: 100%;"
              check-strictly
            />
          </el-form-item>
          <el-form-item label="Is Primary">
            <el-switch v-model="isPrimaryOrg" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="assignOrgDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleAssignOrg" :disabled="!selectedOrgId">
          Assign
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { getUserList, freezeUser, unfreezeUser, adminResetPassword, deleteUser, createUser, updateUser } from '../../api/admin';
import { getRoles, assignRoleToUser, type RoleInfo } from '../../api/role';
import { getOrgTree, type OrgInfo } from '../../api/org';
import api from '../../api/axios';
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus';
import { Plus, Edit, Delete, User, Lock, Unlock, Key, OfficeBuilding, Star } from '@element-plus/icons-vue';

const users = ref<any[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const searchQuery = ref('');

// Create/Edit user
const userDialogVisible = ref(false);
const isEditUser = ref(false);
const editingUserId = ref<number | null>(null);
const submitting = ref(false);
const userFormRef = ref<FormInstance>();
const userForm = reactive({
  username: '',
  email: '',
  password: '',
});
const userRules = reactive<FormRules>({
  username: [
    { required: true, message: 'Username is required', trigger: 'blur' },
    { min: 3, max: 50, message: 'Length should be 3 to 50', trigger: 'blur' },
  ],
  email: [
    { required: true, message: 'Email is required', trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email', trigger: 'blur' },
  ],
  password: [
    { required: true, message: 'Password is required', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' },
  ],
});

// Role assignment
const assignRoleDialogVisible = ref(false);
const selectedUser = ref<any>(null);
const selectedRoleId = ref<number | null>(null);
const availableRoles = ref<RoleInfo[]>([]);

// Org assignment
const assignOrgDialogVisible = ref(false);
const selectedOrgId = ref<number | null>(null);
const orgTree = ref<OrgInfo[]>([]);
const isPrimaryOrg = ref(false);

// Filter out organizations that user already belongs to
const filteredOrgTree = computed(() => {
  if (!selectedUser.value) return orgTree.value;
  
  const userOrgIds = (selectedUser.value.organizations || []).map((o: any) => o.id);
  
  const filterTree = (nodes: OrgInfo[]): OrgInfo[] => {
    return nodes
      .filter((node: any) => !userOrgIds.includes(node.id))
      .map((node: any) => ({
        ...node,
        children: node.children ? filterTree(node.children) : []
      }));
  };
  
  return filterTree(orgTree.value);
});

const fetchUsers = async () => {
  loading.value = true;
  try {
    const res = await getUserList(currentPage.value, pageSize.value, searchQuery.value);
    if (res.data) {
      users.value = res.data.list || [];
      total.value = res.data.total || 0;
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Failed to fetch users');
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  currentPage.value = 1;
  fetchUsers();
};

const openCreateDialog = () => {
  isEditUser.value = false;
  editingUserId.value = null;
  userForm.username = '';
  userForm.email = '';
  userForm.password = '';
  userDialogVisible.value = true;
};

const openEditDialog = (user: any) => {
  isEditUser.value = true;
  editingUserId.value = user.id;
  userForm.username = user.username;
  userForm.email = user.email;
  userForm.password = '';
  userDialogVisible.value = true;
};

const handleUserSubmit = async () => {
  if (!userFormRef.value) return;
  
  await userFormRef.value.validate(async (valid) => {
    if (!valid) return;
    
    submitting.value = true;
    try {
      if (isEditUser.value && editingUserId.value) {
        const res = await updateUser(editingUserId.value, {
          username: userForm.username,
          email: userForm.email,
        });
        if (res.code === 200 || res.code === 0) {
          ElMessage.success('User updated successfully');
          userDialogVisible.value = false;
          fetchUsers();
        } else {
          ElMessage.error(res.msg || 'Failed to update user');
        }
      } else {
        const res = await createUser({
          username: userForm.username,
          email: userForm.email,
          password: userForm.password,
        });
        if (res.code === 200 || res.code === 0) {
          ElMessage.success('User created successfully');
          userDialogVisible.value = false;
          fetchUsers();
        } else {
          ElMessage.error(res.msg || 'Failed to create user');
        }
      }
    } catch (err: any) {
      ElMessage.error(err.response?.data?.msg || 'Operation failed');
    } finally {
      submitting.value = false;
    }
  });
};

const fetchRoles = async () => {
  try {
    const res = await getRoles();
    if (res.data?.roles) {
      availableRoles.value = res.data.roles;
    }
  } catch (err) {
    console.error(err);
  }
};

const buildOrgTree = (list: OrgInfo[]): OrgInfo[] => {
  const map = new Map<number, OrgInfo>();
  const roots: OrgInfo[] = [];
  list.forEach(item => {
    item.children = [];
    map.set(item.ID, item);
  });
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
    if (res.data?.list) {
      orgTree.value = buildOrgTree(res.data.list);
    }
  } catch (err) {
    console.error(err);
  }
};

const userHasRole = (roleId: number): boolean => {
  if (!selectedUser.value?.roles) return false;
  return selectedUser.value.roles.some((r: any) => r.id === roleId);
};

const openAssignRoleDialog = (user: any) => {
  selectedUser.value = user;
  selectedRoleId.value = null;
  assignRoleDialogVisible.value = true;
  fetchRoles();
};

const openAssignOrgDialog = (user: any) => {
  selectedUser.value = user;
  selectedOrgId.value = null;
  isPrimaryOrg.value = false;
  assignOrgDialogVisible.value = true;
  fetchOrgTree();
};

const handleAssignRole = async () => {
  if (!selectedUser.value || !selectedRoleId.value) return;
  try {
    const res = await assignRoleToUser({
      user_id: selectedUser.value.id,
      role_id: selectedRoleId.value,
    });
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Role assigned successfully');
      assignRoleDialogVisible.value = false;
      fetchUsers();
    } else {
      ElMessage.error(res.msg || 'Failed to assign role');
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.msg || 'Failed to assign role');
  }
};

const handleAssignOrg = async () => {
  if (!selectedUser.value || !selectedOrgId.value) return;
  try {
    const res = await api.post('/api/org/user/add', {
      user_id: selectedUser.value.id,
      organization_id: selectedOrgId.value,
      is_primary: isPrimaryOrg.value,
    });
    if (res.data.code === 200 || res.data.code === 0) {
      ElMessage.success('Organization assigned successfully');
      assignOrgDialogVisible.value = false;
      fetchUsers();
    } else {
      ElMessage.error(res.data.msg || 'Failed to assign organization');
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.msg || 'Failed to assign organization');
  }
};

const handleStatusChange = async (user: any) => {
  const action = user.status === 1 ? 'freeze' : 'unfreeze';
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to ${action} user "${user.username}"?`,
      'Confirm',
      {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    if (action === 'freeze') {
      await freezeUser(user.id);
    } else {
      await unfreezeUser(user.id);
    }
    
    ElMessage.success(`User ${action}d successfully`);
    fetchUsers();
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(`Failed to ${action} user`);
    }
  }
};

const handleResetPassword = async (user: any) => {
  try {
    await ElMessageBox.confirm(
      `Reset password for user "${user.username}"?`,
      'Confirm Reset',
      {
        confirmButtonText: 'Reset',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );

    const res = await adminResetPassword(user.id);
    if (res.code === 200 || res.code === 0) {
      ElMessageBox.alert(
        `Password reset successfully. New password: <strong>${res.data.new_password}</strong>`,
        'Success',
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: 'OK',
        }
      );
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to reset password');
    }
  }
};

const handleDelete = async (user: any) => {
  try {
    await ElMessageBox.confirm(
      `Permanently delete user "${user.username}"? This cannot be undone!`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'error',
      }
    );

    await deleteUser(user.id);
    ElMessage.success('User deleted successfully');
    fetchUsers();
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('Failed to delete user');
    }
  }
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString();
};

onMounted(() => {
  fetchUsers();
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
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.text-muted {
  color: #909399;
  font-size: 12px;
}

:deep(.el-table) {
  --el-table-border-color: #ebeef5;
}
</style>
