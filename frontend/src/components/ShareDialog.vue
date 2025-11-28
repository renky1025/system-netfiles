<template>
  <el-dialog
    v-model="dialogVisible"
    title="Share File"
    width="500px"
    @close="handleClose"
  >
    <div v-if="shareCode" class="share-result">
      <el-alert
        title="Share Created Successfully"
        type="success"
        :closable="false"
        show-icon
      />
      <div class="share-link-box">
        <p>Share Link:</p>
        <el-input v-model="shareLink" readonly>
          <template #append>
            <el-button @click="copyLink">Copy</el-button>
          </template>
        </el-input>
      </div>
      <div class="share-code-box">
        <p>Share Code: <strong>{{ shareCode }}</strong></p>
        <p v-if="form.password">Password: <strong>{{ form.password }}</strong></p>
      </div>
    </div>

    <el-form v-else :model="form" label-width="120px">
      <el-form-item label="Password">
        <el-switch v-model="enablePassword" />
        <el-input 
          v-if="enablePassword" 
          v-model="form.password" 
          placeholder="Enter password" 
          type="password" 
          show-password
          style="margin-top: 10px"
        />
      </el-form-item>

      <el-form-item label="Expiration">
        <el-switch v-model="enableExpiration" />
        <el-date-picker
          v-if="enableExpiration"
          v-model="form.expires_at"
          type="datetime"
          placeholder="Select expiration time"
          style="margin-top: 10px; width: 100%"
        />
      </el-form-item>

      <el-form-item label="Download Limit">
        <el-switch v-model="enableLimit" />
        <el-input-number 
          v-if="enableLimit" 
          v-model="form.max_downloads" 
          :min="1" 
          style="margin-top: 10px"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">Close</el-button>
        <el-button v-if="!shareCode" type="primary" @click="handleCreate">Create Share</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, defineProps, defineEmits } from 'vue';
import { createShare } from '../api/share';
import { ElMessage } from 'element-plus';

const props = defineProps<{
  visible: boolean;
  fileId: number | null;
  fileName: string;
}>();

const emit = defineEmits(['update:visible']);

const dialogVisible = ref(false);
const shareCode = ref('');
const enablePassword = ref(false);
const enableExpiration = ref(false);
const enableLimit = ref(false);

const form = ref({
  password: '',
  expires_at: '',
  max_downloads: 10,
});

const shareLink = computed(() => {
  if (!shareCode.value) return '';
  const baseUrl = window.location.origin;
  return `${baseUrl}/share/${shareCode.value}`;
});

watch(() => props.visible, (val) => {
  dialogVisible.value = val;
  if (val) {
    resetForm();
  }
});

watch(dialogVisible, (val) => {
  emit('update:visible', val);
});

const resetForm = () => {
  shareCode.value = '';
  enablePassword.value = false;
  enableExpiration.value = false;
  enableLimit.value = false;
  form.value = {
    password: '',
    expires_at: '',
    max_downloads: 10,
  };
};

const handleClose = () => {
  dialogVisible.value = false;
};

const handleCreate = async () => {
  if (!props.fileId) return;

  try {
    const payload: any = {
      file_id: props.fileId,
    };

    if (enablePassword.value && form.value.password) {
      payload.password = form.value.password;
    }
    if (enableExpiration.value && form.value.expires_at) {
      payload.expires_at = form.value.expires_at;
    }
    if (enableLimit.value) {
      payload.max_downloads = form.value.max_downloads;
    }

    const res = await createShare(payload);
    if (res.data?.share_code) {
      shareCode.value = res.data.share_code;
      ElMessage.success('Share link created');
    } else {
      ElMessage.error(res.msg || res.error || 'Failed to create share');
    }
  } catch (err: any) {
    console.error(err);
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Failed to create share');
  }
};

const copyLink = () => {
  navigator.clipboard.writeText(shareLink.value);
  ElMessage.success('Link copied to clipboard');
};
</script>

<style scoped>
.share-result {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.share-link-box, .share-code-box {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}
.share-link-box p, .share-code-box p {
  margin: 0 0 10px 0;
  font-weight: bold;
  color: #606266;
}
</style>
