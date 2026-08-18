<template>
  <div>
    <Navbar />
    <div class="container">
      <h2>写文章</h2>
      <el-form>
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="文章标题" />
        </el-form-item>
        <el-form-item label="内容（Markdown）">
          <el-input v-model="form.content" type="textarea" :rows="12" placeholder="支持 Markdown 语法" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="tagsText" placeholder="用逗号分隔，如 go,kafka" />
        </el-form-item>
        <el-button type="primary" @click="submit" :loading="loading">发布</el-button>
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
    const res = await api.createArticle({ title: form.title, content: form.content, tags })
    ElMessage.success('文章已创建')
    router.push(`/articles/${res.data.id}`)
  } catch (e) {
    ElMessage.error(e.response?.data?.message || '创建失败')
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
