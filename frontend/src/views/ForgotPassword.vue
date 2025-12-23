<template>
  <div class="auth-container">
    <el-card class="auth-card">
      <template #header>
        <div class="card-header">
          <span>Forgot Password</span>
        </div>
      </template>
      
      <div v-if="submitted" class="success-message">
        <el-result
          icon="success"
          title="Request Sent"
          sub-title="If an account exists with this email, you will receive password reset instructions."
        >
          <template #extra>
            <el-button type="primary" @click="$router.push('/login')">Back to Login</el-button>
          </template>
        </el-result>
      </div>

      <el-form v-else ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="Email" prop="email">
          <el-input v-model="form.email" placeholder="Enter your email" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit" style="width: 100%">
            Send Reset Link
          </el-button>
        </el-form-item>
        
        <div class="links">
          <router-link to="/login">Back to Login</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { requestResetPassword } from '../api/auth';
import { ElMessage } from 'element-plus';

const formRef = ref();
const loading = ref(false);
const submitted = ref(false);

const form = reactive({
  email: '',
});

const rules = {
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' },
    { type: 'email', message: 'Please enter valid email', trigger: 'blur' },
  ],
};

const handleSubmit = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true;
      try {
        const res = await requestResetPassword({ email: form.email });
        if (res.code === 200 || res.code === 0) {
          submitted.value = true;
        } else {
          ElMessage.error(res.msg || res.error || 'Request failed');
        }
      } catch (err: any) {
        console.error(err);
        // For security, don't reveal if email exists or not, or just show generic error
        // But here we might want to show error if it's a system error
        ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Request failed');
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
  min-height: 100vh;
  padding: 40px 16px;
  background-color: #f0f2f5;
}

.auth-card {
  width: 100%;
  max-width: 520px;
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

.links a {
  color: #409eff;
  text-decoration: none;
}
</style>
