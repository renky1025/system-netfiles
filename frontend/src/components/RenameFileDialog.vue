<template>
  <el-dialog
    v-model="dialogVisible"
    title="Rename File"
    width="400px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="100px" @submit.prevent="handleSubmit">
      <el-form-item label="New Name" required>
        <el-input v-model="form.name" placeholder="Enter new file name" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">Cancel</el-button>
        <el-button type="primary" @click="handleSubmit">Rename</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue';
import { renameFile } from '../api/file';
import { ElMessage } from 'element-plus';

const props = defineProps<{
  visible: boolean;
  fileId: number | null;
  fileName: string;
}>();

const emit = defineEmits(['update:visible', 'success']);

const dialogVisible = ref(false);
const form = ref({
  name: '',
});

watch(() => props.visible, (val) => {
  dialogVisible.value = val;
  if (val) {
    form.value.name = props.fileName;
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
    ElMessage.warning('Please enter file name');
    return;
  }

  if (!props.fileId) return;

  try {
    const res = await renameFile(props.fileId, {
      new_name: form.value.name,
    });
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('File renamed successfully');
      emit('success');
      handleClose();
    } else {
      ElMessage.error(res.msg || res.error || 'Operation failed');
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
