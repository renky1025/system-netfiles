<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item :to="{ path: '/' }" @click="handleClick(null)">
      <el-icon><HomeFilled /></el-icon>
      <span>Root</span>
    </el-breadcrumb-item>
    <el-breadcrumb-item
      v-for="item in breadcrumbs"
      :key="item.id"
      @click="handleClick(item.id)"
    >
      {{ item.name }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { getBreadcrumb } from '../api/folder';
import { HomeFilled } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

interface BreadcrumbItem {
  id: number;
  name: string;
  path: string;
}

const props = defineProps<{
  folderId: number | null;
}>();

const emit = defineEmits(['navigate']);

const breadcrumbs = ref<BreadcrumbItem[]>([]);

const fetchBreadcrumb = async (folderId: number | null) => {
  if (!folderId) {
    breadcrumbs.value = [];
    return;
  }

  try {
    const res = await getBreadcrumb(folderId);
    if (res.data?.breadcrumbs) {
      const items = res.data.breadcrumbs;
      // Remove the first "Root" item as we have it in template
      breadcrumbs.value = items.filter((item: BreadcrumbItem) => item.id !== 0);
    } else {
      breadcrumbs.value = [];
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to fetch breadcrumb');
  }
};

const handleClick = (folderId: number | null) => {
  emit('navigate', folderId);
};

watch(() => props.folderId, (newId) => {
  fetchBreadcrumb(newId);
}, { immediate: true });
</script>

<style scoped>
.el-breadcrumb {
  margin-bottom: 20px;
  font-size: 14px;
}

.el-breadcrumb-item {
  cursor: pointer;
}

.el-icon {
  vertical-align: middle;
  margin-right: 4px;
}
</style>
