<template>
  <header class="navbar">
    <router-link to="/" class="logo">Community</router-link>
    <nav class="links">
      <router-link to="/">首页</router-link>
      <router-link to="/search">搜索</router-link>
      <router-link to="/articles/new">写文章</router-link>
      <router-link to="/questions/new">提问</router-link>
      <router-link v-if="auth.isAdmin" to="/admin">后台管理</router-link>
      <router-link v-if="auth.isLoggedIn" to="/profile">个人</router-link>
      <a v-if="!auth.isLoggedIn" href="/login">登录</a>
      <a v-else href="#" @click.prevent="logout">退出</a>
    </nav>
  </header>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}
.logo {
  font-size: 20px;
  font-weight: 700;
  color: #409eff;
  text-decoration: none;
}
.links {
  display: flex;
  gap: 20px;
}
.links a {
  color: #606266;
  text-decoration: none;
}
.links a:hover {
  color: #409eff;
}
</style>
