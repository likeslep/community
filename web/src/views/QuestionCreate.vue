<template>
  <div>
    <Navbar />
    <div class="container">
      <h2>提问</h2>
      <el-form>
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="问题标题" />
        </el-form-item>
        <el-form-item label="描述（Markdown）">
          <el-input v-model="form.content" type="textarea" :rows="12" placeholder="详细描述你的问题" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="tagsText" placeholder="用逗号分隔，如 go,grpc" />
        </el-form-item>
        <el-button type="primary" @click="submit" :loading="loading">发布问题</el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from '../components/Navbar.vue'
import { api } from '../api'

const router = useRouter()
const loading = ref(false)
const tagsText = ref('')
const form = reactive({ title: '', content: '' })

async function submit() {
  loading.value = true
  try {
    const tags = tagsText.value.split(',').map((t) => t.trim()).filter(Boolean)
    const res = await api.createQuestion({ title: form.title, content: form.content, tags })
    ElMessage.success('问题已发布')
    router.push(`/questions/${res.data.id}`)
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '发布失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.container {
  max-width: 760px;
  margin: 24px auto;
  padding: 0 16px;
}
</style>
