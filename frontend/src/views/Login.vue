<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>Login</h2>
      <el-form :model="form" label-width="80px">
        <el-form-item label="Username">
          <el-input v-model="form.username" placeholder="Username" />
        </el-form-item>
        <el-form-item label="Password">
          <el-input v-model="form.password" type="password" placeholder="Password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin">Login</el-button>
          <el-button @click="$router.push('/register')">Register</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { login } from '../api/auth';
import { ElMessage } from 'element-plus';

const router = useRouter();
const userStore = useUserStore();
const form = ref({
  username: '',
  password: '',
});

const handleLogin = async () => {
  try {
    const response = await login(form.value);
    userStore.setToken(response.data.token);
    userStore.setUser(response.data.user);
    router.push('/');
  } catch (err) {
    ElMessage.error(err.response?.data?.msg || err.response?.data?.error || 'Login failed');
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}
.login-card {
  width: 400px;
}
</style>
