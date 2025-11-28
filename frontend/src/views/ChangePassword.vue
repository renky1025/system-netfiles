<template>
  <div class="auth-container">
    <el-card class="auth-card">
      <template #header>
        <div class="card-header">
          <span>Change Password</span>
        </div>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="Old Password" prop="oldPassword">
          <el-input v-model="form.oldPassword" type="password" show-password placeholder="Enter old password" />
        </el-form-item>

        <el-form-item label="New Password" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" show-password placeholder="Enter new password" />
        </el-form-item>

        <el-form-item label="Confirm New Password" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password placeholder="Confirm new password" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit" style="width: 100%">
            Change Password
          </el-button>
        </el-form-item>
        
        <div class="links">
          <el-button text @click="$router.push('/')">Cancel</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { changePassword } from '../api/auth';
import { ElMessage } from 'element-plus';

const router = useRouter();
const formRef = ref();
const loading = ref(false);

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});

const validatePass2 = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('Please input the password again'));
  } else if (value !== form.newPassword) {
    callback(new Error("Two inputs don't match!"));
  } else {
    callback();
  }
};

const rules = {
  oldPassword: [
    { required: true, message: 'Please enter old password', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: 'Please enter new password', trigger: 'blur' },
    { min: 6, message: 'Length should be at least 6 characters', trigger: 'blur' },
  ],
  confirmPassword: [
    { validator: validatePass2, trigger: 'blur' },
  ],
};

const handleSubmit = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true;
      try {
        const res = await changePassword({
          old_password: form.oldPassword,
          new_password: form.newPassword,
        });
        if (res.code === 200 || res.code === 0) {
          ElMessage.success('Password changed successfully');
          router.push('/');
        } else {
          ElMessage.error(res.msg || res.error || 'Change password failed');
        }
      } catch (err: any) {
        console.error(err);
        ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Change password failed');
      } finally {
        loading.value = false;
      }
    }
  });
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.auth-card {
  width: 400px;
}

.card-header {
  text-align: center;
  font-size: 18px;
  font-weight: bold;
}

.links {
  text-align: center;
  margin-top: 10px;
}
</style>
