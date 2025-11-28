<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="400px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="100px">
      <el-form-item label="Folder Name" required>
        <el-input v-model="form.name" placeholder="Enter folder name" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">Cancel</el-button>
        <el-button type="primary" @click="handleSubmit">
          {{ isEdit ? 'Update' : 'Create' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { createFolder, updateFolder } from '../api/folder';
import { ElMessage } from 'element-plus';

interface FolderForm {
  name: string;
}

const props = defineProps<{
  visible: boolean;
  folderId?: number | null;
  parentId?: number | null;
  folderName?: string;
  isEdit?: boolean;
}>();

const emit = defineEmits(['update:visible', 'success']);

const dialogVisible = ref(false);
const form = ref<FolderForm>({
  name: '',
});

const dialogTitle = ref('Create Folder');

watch(() => props.visible, (val) => {
  dialogVisible.value = val;
  if (val) {
    if (props.isEdit && props.folderName) {
      form.value.name = props.folderName;
      dialogTitle.value = 'Rename Folder';
    } else {
      form.value.name = '';
      dialogTitle.value = 'Create Folder';
    }
  }
});

watch(dialogVisible, (val) => {
  emit('update:visible', val);
});

const handleClose = () => {
  dialogVisible.value = false;
  form.value.name = '';
};

const handleSubmit = async () => {
  if (!form.value.name.trim()) {
    ElMessage.warning('Please enter folder name');
    return;
  }

  try {
    if (props.isEdit && props.folderId) {
      // Update folder
      const res = await updateFolder(props.folderId, {
        name: form.value.name,
      });
      if (res.code === 200 || res.code === 0) {
        ElMessage.success('Folder renamed successfully');
        emit('success');
        handleClose();
      } else {
        ElMessage.error(res.msg || res.error || 'Operation failed');
      }
    } else {
      // Create folder
      const res = await createFolder({
        name: form.value.name,
        parent_id: props.parentId || null,
      });
      if (res.code === 200 || res.code === 0) {
        ElMessage.success('Folder created successfully');
        emit('success');
        handleClose();
      } else {
        ElMessage.error(res.msg || res.error || 'Operation failed');
      }
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Operation failed');
  }
};
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
