<template>
  <div class="folder-tree">
    <el-tree
      :data="treeData"
      :props="treeProps"
      node-key="id"
      :default-expanded-keys="expandedKeys"
      :highlight-current="true"
      @node-click="handleNodeClick"
    >
      <template #default="{ node, data }">
        <span class="custom-tree-node">
          <el-icon><Folder /></el-icon>
          <span class="node-label">{{ node.label }}</span>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getFolderTree } from '../api/folder';
import { Folder } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

interface FolderNode {
  id: number;
  name: string;
  parent_id: number | null;
  children: FolderNode[];
}

const emit = defineEmits(['folder-selected']);

const treeData = ref<FolderNode[]>([]);
const expandedKeys = ref<number[]>([]);
const treeProps = {
  children: 'children',
  label: 'name',
};

const fetchFolderTree = async () => {
  try {
    const res = await getFolderTree();
    if (res.data?.tree) {
      // Convert FolderInfo[] to FolderNode[] format
      const convertToFolderNode = (folder: any): FolderNode => ({
        id: folder.id || folder.ID || 0,
        name: folder.name || folder.Name || '',
        parent_id: folder.parent_id || folder.ParentID || null,
        children: folder.children ? folder.children.map(convertToFolderNode) : [],
      });
      treeData.value = res.data.tree.map(convertToFolderNode);
    } else {
      treeData.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch folder tree');
  }
};

const handleNodeClick = (data: FolderNode) => {
  emit('folder-selected', data.id);
};

onMounted(() => {
  fetchFolderTree();
});

defineExpose({
  refresh: fetchFolderTree,
});
</script>

<style scoped>
.folder-tree {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  max-height: 400px;
  overflow-y: auto;
}

.custom-tree-node {
  display: flex;
  align-items: center;
  gap: 5px;
}

.node-label {
  font-size: 14px;
}
</style>
