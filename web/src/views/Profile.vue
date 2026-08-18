<template>
  <div>
    <Navbar />
    <div class="container">
      <h2>个人资料</h2>
      <el-descriptions v-if="auth.user" border :column="1">
        <el-descriptions-item label="用户名">{{ auth.user.username }}</el-descriptions-item>
        <el-descriptions-item label="角色">{{ auth.user.role }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ profile.email || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="简介">{{ profile.bio || '未设置' }}</el-descriptions-item>
      </el-descriptions>

      <h3 style="margin-top: 24px">更新资料</h3>
      <el-form>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="邮箱" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.bio" type="textarea" placeholder="个人简介" />
        </el-form-item>
        <el-button type="primary" @click="submit">保存</el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Navbar from '../components/Navbar.vue'
import { api } from '../api'
import { useAuthStore } from '../store/auth'

const auth = useAuthStore()
const profile = ref({})
const form = reactive({ email: '', bio: '' })

onMounted(async () => {
  if (!auth.user) return
  try {
    const res = await api.getProfile(auth.user.id)
    profile.value = res.data
    form.email = res.data.email || ''
    form.bio = res.data.bio || ''
  } catch (e) {
    console.error(e)
  }
})

async function submit() {
  try {
    await api.updateProfile({ email: form.email, bio: form.bio })
    ElMessage.success('已保存')
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '保存失败')
  }
}
</script>

<style scoped>
.container {
  max-width: 560px;
  margin: 24px auto;
  padding: 0 16px;
}
</style>
