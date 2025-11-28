<template>
  <div class="register-container">
    <el-card class="register-card">
      <h2>Register</h2>
      <el-form :model="form" label-width="80px">
        <el-form-item label="Username">
          <el-input v-model="form.username" placeholder="Username" />
        </el-form-item>
        <el-form-item label="Email">
          <el-input v-model="form.email" placeholder="Email" />
        </el-form-item>
        <el-form-item label="Password">
          <el-input v-model="form.password" type="password" placeholder="Password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister">Register</el-button>
          <el-button @click="$router.push('/login')">Login</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { register } from '../api/auth';
import { ElMessage } from 'element-plus';

const router = useRouter();
const form = ref({
  username: '',
  password: '',
  email: '',
});

const handleRegister = async () => {
  try {
    const res = await register(form.value);
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('Registration successful');
      router.push('/login');
    } else {
      ElMessage.error(res.msg || res.error || 'Registration failed');
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Registration failed');
  }
};
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}
.register-card {
  width: 400px;
}
</style>
