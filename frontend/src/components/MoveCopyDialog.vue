<template>
  <el-dialog
    :title="mode === 'move' ? 'Move to' : 'Copy to'"
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    width="500px"
  >
    <div class="dialog-content">
      <p>Select destination folder:</p>
      <div class="folder-tree-container">
        <el-tree
          ref="treeRef"
          :data="folderTree"
          :props="defaultProps"
          node-key="id"
          highlight-current
          :expand-on-click-node="false"
          @current-change="handleNodeClick"
          default-expand-all
        >
          <template #default="{ node, data }">
            <span class="custom-tree-node">
              <el-icon><Folder /></el-icon>
              <span style="margin-left: 8px">{{ node.label }}</span>
            </span>
          </template>
        </el-tree>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="$emit('update:visible', false)">Cancel</el-button>
        <el-button type="primary" @click="handleConfirm" :disabled="!selectedFolderId">
          {{ mode === 'move' ? 'Move' : 'Copy' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits, watch } from 'vue';
import { Folder } from '@element-plus/icons-vue';
import { getFolderTree } from '../api/folder';
import { moveFile, copyFile, batchMoveFiles, batchCopyFiles } from '../api/file';
import { moveFolder } from '../api/folder';
import { ElMessage } from 'element-plus';

const props = defineProps<{
  visible: boolean;
  mode: 'move' | 'copy';
  items: any[]; // Array of files or folders to move/copy
}>();

const emit = defineEmits(['update:visible', 'success']);

const folderTree = ref<any[]>([]);
const selectedFolderId = ref<number | null>(null);
const treeRef = ref();

const defaultProps = {
  children: 'children',
  label: 'name',
};

const fetchFolderTree = async () => {
  try {
    const res = await getFolderTree();
    if (res.data?.tree) {
      folderTree.value = res.data.tree;
    } else {
      folderTree.value = [];
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Failed to load folder tree');
  }
};

const handleNodeClick = (data: any) => {
  selectedFolderId.value = data.id;
};

const handleConfirm = async () => {
  if (!selectedFolderId.value) return;

  try {
    const folderId = selectedFolderId.value;
    const itemIds = props.items.map(item => item.id);
    const isBatch = props.items.length > 1;

    if (props.mode === 'move') {
      if (isBatch) {
        const res = await batchMoveFiles({
          file_ids: itemIds, // Note: This only handles files for now based on API. Need to check if folder move is supported in batch.
          folder_id: folderId,
        });
        if (res.code === 200 || res.code === 0) {
          ElMessage.success('Moved successfully');
          emit('success');
          emit('update:visible', false);
        } else {
          ElMessage.error(res.msg || res.error || 'Move failed');
        }
      } else {
        const item = props.items[0];
        try {
          if (item.type === 'file') {
            const res = await moveFile(item.id, { folder_id: folderId });
            if (res.code === 200 || res.code === 0) {
              ElMessage.success('Moved successfully');
              emit('success');
              emit('update:visible', false);
            } else {
              ElMessage.error(res.msg || res.error || 'Move failed');
            }
          } else {
            const res = await moveFolder(item.id, { parent_id: folderId });
            if (res.code === 200 || res.code === 0) {
              ElMessage.success('Moved successfully');
              emit('success');
              emit('update:visible', false);
            } else {
              ElMessage.error(res.msg || res.error || 'Move failed');
            }
          }
        } catch (err: any) {
          console.error(err);
          ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Move failed');
        }
      }
    } else {
      // Copy
      if (isBatch) {
        // Filter out folders as folder copy is not supported yet
        const fileIds = props.items.filter(i => i.type === 'file').map(i => i.id);
        if (fileIds.length === 0) {
           ElMessage.warning('No files to copy (folder copy not supported)');
           return;
        }

        const res = await batchCopyFiles({
          file_ids: fileIds,
          folder_id: folderId,
        });

        if (res.code === 200 || res.code === 0) {
          ElMessage.success('Copied successfully');
          emit('success');
          emit('update:visible', false);
        } else {
          ElMessage.error(res.msg || res.error || 'Copy failed');
        }
      } else {
        // Single item copy
        const item = props.items[0];
        if (item.type === 'file') {
            const res = await copyFile(item.id, { 
              folder_id: folderId,
              new_name: item.name 
            });
            if (res.code === 200 || res.code === 0) {
              ElMessage.success('Copied successfully');
              emit('success');
              emit('update:visible', false);
            } else {
              ElMessage.error(res.msg || res.error || 'Copy failed');
            }
        } else {
             ElMessage.warning('Folder copy is not supported yet');
        }
      }
    }
  } catch (err) {
    console.error(err);
    ElMessage.error('Operation failed');
  }
};

watch(() => props.visible, (newVal) => {
  if (newVal) {
    fetchFolderTree();
    selectedFolderId.value = null;
  }
});
</script>

<style scoped>
.folder-tree-container {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  margin-top: 10px;
}
.custom-tree-node {
  display: flex;
  align-items: center;
}
</style>
