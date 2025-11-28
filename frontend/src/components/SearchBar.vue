<template>
  <div class="search-bar">
    <el-input
      v-model="searchQuery"
      placeholder="搜索文件和文件夹..."
      class="search-input"
      clearable
      @input="onSearchInput"
      @keyup.enter="performSearch"
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
      <template #append>
        <el-button :icon="Search" @click="performSearch">搜索</el-button>
      </template>
    </el-input>

    <!-- Search Suggestions -->
    <div v-if="showSuggestions && suggestions.length > 0" class="suggestions-dropdown">
      <div
        v-for="(suggestion, index) in suggestions"
        :key="index"
        class="suggestion-item"
        @click="selectSuggestion(suggestion)"
      >
        <el-icon><Document /></el-icon>
        <span>{{ suggestion }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Search, Document } from '@element-plus/icons-vue';
import { getSearchSuggestions } from '../api/search';

const searchQuery = ref('');
const suggestions = ref<string[]>([]);
const showSuggestions = ref(false);
let suggestionTimer: NodeJS.Timeout | null = null;

const emit = defineEmits(['search']);

const onSearchInput = () => {
  if (suggestionTimer) {
    clearTimeout(suggestionTimer);
  }

  if (searchQuery.value.trim().length < 2) {
    suggestions.value = [];
    showSuggestions.value = false;
    return;
  }

  suggestionTimer = setTimeout(async () => {
    try {
      const res = await getSearchSuggestions(searchQuery.value);
      if (res.data) {
        suggestions.value = res.data;
        showSuggestions.value = true;
      }
    } catch (error) {
      console.error('Failed to get suggestions:', error);
    }
  }, 300);
};

const performSearch = () => {
  if (searchQuery.value.trim()) {
    showSuggestions.value = false;
    emit('search', searchQuery.value);
  }
};

const selectSuggestion = (suggestion: string) => {
  searchQuery.value = suggestion;
  showSuggestions.value = false;
  performSearch();
};

// Close suggestions when clicking outside
watch(showSuggestions, (newVal) => {
  if (newVal) {
    document.addEventListener('click', closeSuggestions);
  } else {
    document.removeEventListener('click', closeSuggestions);
  }
});

const closeSuggestions = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.search-bar')) {
    showSuggestions.value = false;
  }
};
</script>

<style scoped>
.search-bar {
  position: relative;
  width: 100%;
  max-width: 600px;
}

.search-input {
  width: 100%;
}

.suggestions-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-top: 4px;
  max-height: 300px;
  overflow-y: auto;
  z-index: 1000;
}

.suggestion-item {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.suggestion-item:hover {
  background-color: #f5f7fa;
}

.suggestion-item .el-icon {
  margin-right: 10px;
  color: #909399;
}

.suggestion-item span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
