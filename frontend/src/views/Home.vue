<template>
  <div class="home-container">
    <el-header class="header">
      <div class="logo">NetFileSys</div>
      <div class="search-bar">
        <el-autocomplete
          v-model="searchQuery"
          :fetch-suggestions="fetchSearchSuggestions"
          placeholder="Search files and folders..."
          prefix-icon="Search"
          clearable
          class="search-input"
          @select="handleSearchSelect"
          @keyup.enter="performSearch"
          @clear="clearSearch"
        >
          <template #suffix>
            <el-button v-if="searchQuery" link @click="performSearch">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-autocomplete>
      </div>
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
      <!-- Left Sidebar: Folder Tree -->
      <el-aside width="250px" class="sidebar">
        <div class="sidebar-header">
          <h3>Folders</h3>
          <el-button size="small" type="primary" @click="showCreateFolderDialog">
            <el-icon><FolderAdd /></el-icon>
          </el-button>
        </div>
        <FolderTree ref="folderTreeRef" @folder-selected="handleFolderSelected" />
        
        <!-- Storage Quota Display -->
        <div class="sidebar-quota">
          <StorageQuotaCard ref="quotaCardRef" />
        </div>
      </el-aside>

      <!-- Main Content -->
      <el-main>
        <!-- Breadcrumb Navigation -->
        <Breadcrumb :folder-id="currentFolderId" @navigate="handleBreadcrumbNavigate" />

        <!-- Toolbar -->
        <div class="toolbar">
          <div class="toolbar-group">
            <el-button @click="showCreateFolderDialog" type="primary">
              <el-icon><FolderAdd /></el-icon>
              New Folder
            </el-button>
            <FileUploader :current-folder-id="currentFolderId" @upload-success="refreshAll" />
          </div>
          
          <div class="toolbar-group" v-if="selectedItems.length > 0">
             <el-button-group>
                <el-button type="primary" :icon="Rank" @click="openBatchMove">Move</el-button>
                <el-button type="primary" :icon="CopyDocument" @click="openBatchCopy">Copy</el-button>
                <el-button type="danger" :icon="Delete" @click="handleBatchDelete">Delete</el-button>
             </el-button-group>
          </div>

          <div class="toolbar-group" v-else>
            <el-button @click="createShare" :disabled="!selectedFile">
              <el-icon><Share /></el-icon>
              Share
            </el-button>
            <el-button @click="downloadFileHandler" :disabled="!selectedFile" type="success">
              <el-icon><Download /></el-icon>
              Download
            </el-button>
          </div>
        </div>

        <!-- Search Mode Indicator -->
        <div v-if="isSearchMode" class="search-mode-bar">
          <span>Search results for: "{{ searchQuery }}"</span>
          <el-button size="small" @click="clearSearch">Clear Search</el-button>
        </div>

        <!-- Files and Folders Table -->
        <el-table 
          :data="displayItems" 
          style="width: 100%" 
          @selection-change="handleSelectionChange"
          @row-dblclick="handleRowDoubleClick"
          v-loading="loading"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="Name" min-width="300">
            <template #default="scope">
              <div class="file-name-cell" @click="handleItemClick(scope.row)">
                <el-icon v-if="scope.row.type === 'folder'" class="folder-icon">
                  <Folder />
                </el-icon>
                <el-icon v-else class="file-icon">
                  <Document />
                </el-icon>
                <span class="file-name">{{ scope.row.name }}</span>
                <el-tag v-if="scope.row.isSearchResult" size="small" type="info" style="margin-left: 8px;">
                  {{ scope.row.path }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="size" label="Size" width="120">
            <template #default="scope">
              {{ scope.row.type === 'folder' ? '-' : formatSize(scope.row.size) }}
            </template>
          </el-table-column>
          <el-table-column label="Modified" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.updated_at || scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="280" fixed="right">
            <template #default="scope">
              <div class="action-buttons">
                <!-- File actions -->
                <template v-if="scope.row.type === 'file'">
                  <el-button size="small" type="primary" :icon="View" @click="handlePreview(scope.row)" title="Preview" />
                  <el-button size="small" type="success" :icon="Download" @click="handleDownload(scope.row)" title="Download" />
                  <el-button size="small" :icon="Share" @click="handleShare(scope.row)" title="Share" />
                  <el-button size="small" :icon="Edit" @click="handleRename(scope.row)" title="Rename" />
                  <el-button size="small" :icon="Clock" @click="handleVersionHistory(scope.row)" title="Version History" />
                  <el-button size="small" type="danger" :icon="Delete" @click="handleDeleteFile(scope.row)" title="Delete" />
                </template>
                <!-- Folder actions -->
                <template v-else>
                  <el-button size="small" type="primary" @click="handleOpenFolder(scope.row)" title="Open">Open</el-button>
                  <el-button size="small" :icon="Edit" @click="handleRenameFolder(scope.row)" title="Rename" />
                  <el-button size="small" type="danger" :icon="Delete" @click="handleDeleteFolder(scope.row)" title="Delete" />
                </template>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
    </el-container>

    <!-- Folder Dialog -->
    <FolderDialog
      v-model:visible="folderDialogVisible"
      :folder-id="editingFolderId"
      :parent-id="currentFolderId"
      :folder-name="editingFolderName"
      :is-edit="isEditMode"
      @success="handleFolderOperationSuccess"
    />

    <!-- Share Dialog -->
    <!-- Share Dialog -->
    <ShareDialog
      v-model:visible="shareDialogVisible"
      :file-id="selectedFile?.ID || selectedFile?.id"
      :file-name="selectedFile?.Name || selectedFile?.name"
    />

    <!-- Move/Copy Dialog -->
    <MoveCopyDialog
      v-model:visible="moveCopyDialogVisible"
      :mode="moveCopyMode"
      :items="moveCopyItems"
      @success="refreshAll"
    />

    <!-- Rename File Dialog -->
    <RenameFileDialog
      v-model:visible="renameFileDialogVisible"
      :file-id="renamingFileId"
      :file-name="renamingFileName"
      @success="refreshAll"
    />

    <!-- File Preview -->
    <FilePreview
      v-model:visible="previewVisible"
      :file-id="previewFileId"
      :file-name="previewFileName"
      :file-size="previewFileSize"
    />

    <!-- Version History -->
    <VersionHistoryDialog
      v-model:visible="versionHistoryVisible"
      :file-id="versionFileId"
      @success="refreshAll"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  Document, 
  Folder, 
  FolderAdd, 
  Share, 
  Download, 
  Delete,
  Edit, 
  Rank, 
  CopyDocument,
  ArrowDown,
  Search,
  View,
  Clock
} from '@element-plus/icons-vue';

import { getFileList, downloadFile, deleteFile as deleteFileApi, batchDeleteFiles } from '../api/file';
import { search, getSearchSuggestions, type SearchResult } from '../api/search';
import { getFolderList, deleteFolder } from '../api/folder';
import { useUserStore } from '../store/user';
import FileUploader from '../components/FileUploader.vue';
import FolderTree from '../components/FolderTree.vue';
import Breadcrumb from '../components/Breadcrumb.vue';
import FolderDialog from '../components/FolderDialog.vue';
import ShareDialog from '../components/ShareDialog.vue';
import MoveCopyDialog from '../components/MoveCopyDialog.vue';
import RenameFileDialog from '../components/RenameFileDialog.vue';
import FilePreview from '../components/FilePreview.vue';
import VersionHistoryDialog from '../components/VersionHistoryDialog.vue';
import StorageQuotaCard from '../components/StorageQuotaCard.vue';

const userStore = useUserStore();
const router = useRouter();
const route = useRoute();

// Data
const loading = ref(false);
const files = ref<any[]>([]);
const folders = ref<any[]>([]);
const currentFolderId = ref<number | null>(null);
const selectedFile = ref<any>(null);
const selectedItems = ref<any[]>([]);
const folderTreeRef = ref();
const quotaCardRef = ref();
const searchQuery = ref('');
const isSearchMode = ref(false);
const searchResults = ref<SearchResult[]>([]);
const searchLoading = ref(false);

// Dialogs State
const shareDialogVisible = ref(false);

const folderDialogVisible = ref(false);
const editingFolderId = ref<number | null>(null);
const editingFolderName = ref('');
const isEditMode = ref(false);

const moveCopyDialogVisible = ref(false);
const moveCopyMode = ref<'move' | 'copy'>('move');
const moveCopyItems = ref<any[]>([]);

const renameFileDialogVisible = ref(false);
const renamingFileId = ref<number | null>(null);
const renamingFileName = ref('');

const previewVisible = ref(false);
const previewFileId = ref<number | null>(null);
const previewFileName = ref('');
const previewFileSize = ref(0);

const versionHistoryVisible = ref(false);
const versionFileId = ref<number | null>(null);

// Computed
const displayItems = computed(() => {
  // If in search mode, show search results
  if (isSearchMode.value && searchResults.value.length > 0) {
    return searchResults.value.map(r => ({
      ...r,
      type: r.type,
      name: r.name,
      size: r.size || 0,
      id: r.id,
      ID: r.id,
      isSearchResult: true,
    }));
  }

  const folderItems = folders.value.map(f => ({
    ...f,
    type: 'folder',
    name: f.Name || f.name,
    size: 0,
    updated_at: f.UpdatedAt || f.updated_at,
  }));
  
  const fileItems = files.value.map(f => ({
    ...f,
    type: 'file',
    name: f.Name || f.name,
    size: f.Size || f.size,
    updated_at: f.UpdatedAt || f.updated_at,
  }));
  
  return [...folderItems, ...fileItems];
});

// Search functions
const fetchSearchSuggestions = async (queryString: string, cb: (suggestions: any[]) => void) => {
  if (!queryString || queryString.length < 2) {
    cb([]);
    return;
  }
  try {
    const res = await getSearchSuggestions(queryString, 5);
    const suggestions = (res.data || []).map((s: string) => ({ value: s }));
    cb(suggestions);
  } catch (err) {
    cb([]);
  }
};

const performSearch = async () => {
  if (!searchQuery.value || searchQuery.value.length < 2) {
    ElMessage.warning('Please enter at least 2 characters to search');
    return;
  }
  
  searchLoading.value = true;
  isSearchMode.value = true;
  try {
    const res = await search({ query: searchQuery.value, type: 'all', page: 1, page_size: 50 });
    if (res.data?.results) {
      searchResults.value = res.data.results;
      if (searchResults.value.length === 0) {
        ElMessage.info('No results found');
      }
    } else {
      searchResults.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error('Search failed');
    searchResults.value = [];
  } finally {
    searchLoading.value = false;
  }
};

const handleSearchSelect = (item: any) => {
  searchQuery.value = item.value;
  performSearch();
};

const clearSearch = () => {
  isSearchMode.value = false;
  searchResults.value = [];
  searchQuery.value = '';
};

// API Calls
const fetchFiles = async () => {
  try {
    const res = await getFileList(currentFolderId.value);
    if (res.data?.files) {
      files.value = res.data.files;
    } else {
      files.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch files');
  }
};

const fetchFolders = async () => {
  try {
    const res = await getFolderList(currentFolderId.value);
    if (res.data?.folders) {
      folders.value = res.data.folders;
    } else {
      folders.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch folders');
  }
};

const refreshAll = async () => {
  loading.value = true;
  try {
    await Promise.all([fetchFiles(), fetchFolders()]);
    if (folderTreeRef.value) {
      folderTreeRef.value.refresh();
    }
  } finally {
    loading.value = false;
  }
  selectedItems.value = [];
  selectedFile.value = null;
};

// Handlers
const handleFolderSelected = (folderId: number) => {
  currentFolderId.value = folderId;
  // Update URL to reflect folder navigation
  if (folderId) {
    router.push(`/folder/${folderId}`);
  } else {
    router.push('/');
  }
  refreshAll();
};

const handleBreadcrumbNavigate = (folderId: number | null) => {
  currentFolderId.value = folderId;
  // Update URL to reflect folder navigation
  if (folderId) {
    router.push(`/folder/${folderId}`);
  } else {
    router.push('/');
  }
  refreshAll();
};

const handleRowDoubleClick = (row: any) => {
  if (row.type === 'folder') {
    currentFolderId.value = row.ID || row.id;
    refreshAll();
  } else {
    // Preview file
    previewFileId.value = row.ID || row.id;
    previewFileName.value = row.Name || row.name;
    previewFileSize.value = row.Size || row.size;
    previewVisible.value = true;
  }
};

const handleSelectionChange = (val: any[]) => {
  selectedItems.value = val;
  if (val.length === 1 && val[0].type === 'file') {
    selectedFile.value = val[0];
  } else {
    selectedFile.value = null;
  }
};

// Folder Operations
const showCreateFolderDialog = () => {
  isEditMode.value = false;
  editingFolderId.value = null;
  editingFolderName.value = '';
  folderDialogVisible.value = true;
};

const handleRenameFolder = (folder: any) => {
  isEditMode.value = true;
  editingFolderId.value = folder.ID || folder.id;
  editingFolderName.value = folder.Name || folder.name;
  folderDialogVisible.value = true;
};

const handleDeleteFolder = async (folder: any) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete folder "${folder.Name || folder.name}"?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );
    
    const res = await deleteFolder(folder.ID || folder.id);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Folder deleted successfully');
      refreshAll();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to delete folder');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to delete folder');
    }
  }
};

const handleFolderOperationSuccess = () => {
  refreshAll();
};

// File Operations
const createShare = async () => {
  if (!selectedFile.value) return;
  shareDialogVisible.value = true;
};

// ...



const downloadFileHandler = async () => {
  if (!selectedFile.value) return;
  try {
    const fileId = selectedFile.value.ID || selectedFile.value.id;
    const blob = await downloadFile(fileId);
    
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', selectedFile.value.Name || selectedFile.value.name);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    
    ElMessage.success('Download started');
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to download file');
  }
};

// Single file delete (from toolbar) - handled by handleDeleteFile

const handleDeleteFile = async (file: any) => {
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete "${file.Name || file.name}"?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );
    
    const fileId = file.ID || file.id;
    const res = await deleteFileApi(fileId);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('File deleted successfully');
      refreshAll();
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to delete file');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to delete file');
    }
  }
};

// Batch & Dialog Handlers
const openMoveDialog = (item: any) => {
  moveCopyMode.value = 'move';
  moveCopyItems.value = [item];
  moveCopyDialogVisible.value = true;
};

const openCopyDialog = (item: any) => {
  moveCopyMode.value = 'copy';
  moveCopyItems.value = [item];
  moveCopyDialogVisible.value = true;
};

const openBatchMove = () => {
  moveCopyMode.value = 'move';
  moveCopyItems.value = selectedItems.value;
  moveCopyDialogVisible.value = true;
};

const openBatchCopy = () => {
  moveCopyMode.value = 'copy';
  moveCopyItems.value = selectedItems.value;
  moveCopyDialogVisible.value = true;
};

const handleBatchDelete = async () => {
  const fileItems = selectedItems.value.filter(item => item.type === 'file');
  if (fileItems.length === 0) {
    ElMessage.warning('Please select files to delete');
    return;
  }
  
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete ${fileItems.length} file(s)?`,
      'Confirm Delete',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );
    
    const fileIds = fileItems.map(item => item.ID || item.id);
    const res = await batchDeleteFiles({ file_ids: fileIds });
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Files deleted successfully');
      refreshAll();
    } else {
      ElMessage.error(res.msg || 'Failed to delete files');
    }
  } catch (err: any) {
    if (err !== 'cancel') {
      console.error(err);
      ElMessage.error('Failed to delete files');
    }
  }
};

const handleCommand = (command: string) => {
  if (command === 'logout') {
    handleLogout();
  } else if (command === 'change-password') {
    router.push('/change-password');
  } else if (command === 'admin') {
    router.push('/admin');
  }
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString();
};

// Item click handlers
const handleItemClick = (item: any) => {
  if (item.type === 'folder') {
    handleOpenFolder(item);
  }
};

const handleOpenFolder = (folder: any) => {
  const folderId = folder.ID || folder.id;
  currentFolderId.value = folderId;
  router.push(`/folder/${folderId}`);
};

const handlePreview = (file: any) => {
  previewFileId.value = file.ID || file.id;
  previewFileName.value = file.Name || file.name;
  previewFileSize.value = file.Size || file.size;
  previewVisible.value = true;
};

const handleDownload = async (file: any) => {
  try {
    const fileId = file.ID || file.id;
    const res = await downloadFile(fileId);
    const blob = res instanceof Blob ? res : new Blob([res]);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', file.Name || file.name);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (err) {
    ElMessage.error('Download failed');
  }
};

const handleShare = (file: any) => {
  selectedFile.value = file;
  shareDialogVisible.value = true;
};

const handleRename = (file: any) => {
  renamingFileId.value = file.ID || file.id;
  renamingFileName.value = file.Name || file.name;
  renameFileDialogVisible.value = true;
};

const handleVersionHistory = (file: any) => {
  versionFileId.value = file.ID || file.id;
  versionHistoryVisible.value = true;
};

// Watch route changes to handle folder navigation
watch(() => route.params.folderId, (newFolderId) => {
  if (newFolderId) {
    currentFolderId.value = parseInt(newFolderId as string);
  } else {
    currentFolderId.value = null;
  }
  refreshAll();
}, { immediate: false });

onMounted(() => {
  // Initialize from route params
  if (route.params.folderId) {
    currentFolderId.value = parseInt(route.params.folderId as string);
  }
  refreshAll();
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

.search-bar {
  flex: 1;
  max-width: 400px;
  margin: 0 20px;
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

.toolbar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
}

.toolbar-group {
  display: flex;
  gap: 10px;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.file-name-cell:hover .file-name {
  color: #409eff;
  text-decoration: underline;
}

.folder-icon {
  color: #409eff;
  font-size: 18px;
}

.file-icon {
  color: #67c23a;
  font-size: 18px;
}

.action-buttons {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  padding: 4px 8px;
}

.search-mode-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #ecf5ff;
  border-radius: 4px;
  margin-bottom: 15px;
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.sidebar-quota {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #dcdfe6;
}
</style>
