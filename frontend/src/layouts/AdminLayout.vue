<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="220px" class="admin-sidebar">
        <div class="sidebar-header">
          <div class="logo">NetFileSys Admin</div>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          router
        >
          <el-menu-item index="/admin/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>Dashboard</span>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <span>User Management</span>
          </el-menu-item>
          <el-menu-item index="/admin/orgs">
            <el-icon><OfficeBuilding /></el-icon>
            <span>Organizations</span>
          </el-menu-item>
          <el-menu-item index="/admin/roles">
            <el-icon><Key /></el-icon>
            <span>Roles & Permissions</span>
          </el-menu-item>
          <el-menu-item index="/admin/files">
            <el-icon><Document /></el-icon>
            <span>File Management</span>
          </el-menu-item>
          <el-menu-item index="/admin/shares">
            <el-icon><Share /></el-icon>
            <span>Share Management</span>
          </el-menu-item>
          <el-menu-item index="/admin/logs">
            <el-icon><Notebook /></el-icon>
            <span>Audit Logs</span>
          </el-menu-item>
          <el-menu-item index="/admin/settings">
            <el-icon><Setting /></el-icon>
            <span>System Settings</span>
          </el-menu-item>
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <span>Back to User Portal</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header class="admin-header">
          <div class="header-left">
            <!-- Breadcrumb could go here -->
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="el-dropdown-link">
                Admin <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">Logout</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main class="admin-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import {
  Odometer,
  User,
  Document,
  Share,
  Notebook,
  HomeFilled,
  ArrowDown,
  OfficeBuilding,
  Key,
  Setting
} from '@element-plus/icons-vue';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const activeMenu = computed(() => route.path);

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout();
    router.push('/login');
  }
};
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  display: flex;
}

.admin-sidebar {
  background-color: #304156;
  color: white;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b3649;
}

.logo {
  font-weight: bold;
  font-size: 18px;
}

.el-menu-vertical {
  border-right: none;
  flex: 1;
}

.admin-header {
  background-color: white;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}

.admin-main {
  background-color: #f0f2f5;
  padding: 0;
  overflow-y: auto;
}

/* Global admin page styles */
:deep(.page-container) {
  padding: 20px;
}

:deep(.page-card) {
  min-height: calc(100vh - 140px);
}

:deep(.card-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.card-header h2) {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

:deep(.header-actions) {
  display: flex;
  align-items: center;
  gap: 10px;
}

:deep(.pagination-wrapper) {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

:deep(.el-table) {
  --el-table-border-color: #ebeef5;
}

:deep(.text-muted) {
  color: #909399;
  font-size: 12px;
}
</style>
