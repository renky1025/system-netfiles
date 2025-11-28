<template>
  <div class="auth-container">
    <el-card class="auth-card">
      <template #header>
        <div class="card-header">
          <span>Reset Password</span>
        </div>
      </template>

      <div v-if="submitted" class="success-message">
        <el-result
          icon="success"
          title="Password Reset Successful"
          sub-title="You can now login with your new password."
        >
          <template #extra>
            <el-button type="primary" @click="$router.push('/login')">Go to Login</el-button>
          </template>
        </el-result>
      </div>

      <el-form v-else ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="New Password" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="Enter new password" />
        </el-form-item>

        <el-form-item label="Confirm Password" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password placeholder="Confirm new password" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit" style="width: 100%">
            Reset Password
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { resetPassword } from '../api/auth';
import { ElMessage } from 'element-plus';

const route = useRoute();
const router = useRouter();
const formRef = ref();
const loading = ref(false);
const submitted = ref(false);
const token = ref('');

const form = reactive({
  password: '',
  confirmPassword: '',
});

const validatePass2 = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('Please input the password again'));
  } else if (value !== form.password) {
    callback(new Error("Two inputs don't match!"));
  } else {
    callback();
  }
};

const rules = {
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' },
    { min: 6, message: 'Length should be at least 6 characters', trigger: 'blur' },
  ],
  confirmPassword: [
    { validator: validatePass2, trigger: 'blur' },
  ],
};

onMounted(() => {
  token.value = route.query.token as string;
  if (!token.value) {
    ElMessage.error('Invalid or missing reset token');
    // router.push('/login'); // Optional: redirect
  }
});

const handleSubmit = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      if (!token.value) {
        ElMessage.error('Missing reset token');
        return;
      }

      loading.value = true;
      try {
        const res = await resetPassword({
          token: token.value,
          new_password: form.password,
        });
        if (res.code === 200 || res.code === 0) {
          submitted.value = true;
        } else {
          ElMessage.error(res.msg || res.error || 'Reset failed');
        }
      } catch (err: any) {
        console.error(err);
        ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Reset failed');
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
</style>
