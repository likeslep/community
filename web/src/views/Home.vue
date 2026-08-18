<template>
  <div>
    <Navbar />
    <div class="container">
      <el-tabs v-model="tab">
        <el-tab-pane label="最新文章" name="articles">
          <el-empty v-if="articles.length === 0" description="暂无文章" />
          <el-card v-for="a in articles" :key="a.id" class="item" shadow="hover">
            <router-link :to="`/articles/${a.id}`" class="title">{{ a.title }}</router-link>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="最新问题" name="questions">
          <el-empty v-if="questions.length === 0" description="暂无问题" />
          <el-card v-for="q in questions" :key="q.id" class="item" shadow="hover">
            <router-link :to="`/questions/${q.id}`" class="title">{{ q.title }}</router-link>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '../components/Navbar.vue'
import { api } from '../api'

const tab = ref('articles')
const articles = ref([])
const questions = ref([])

onMounted(async () => {
  try {
    const [a, q] = await Promise.all([
      api.listArticles({ limit: 20, status: 'published' }),
      api.listQuestions({ limit: 20 })
    ])
    articles.value = a.data.articles || []
    questions.value = q.data.questions || []
  } catch (e) {
    console.error(e)
  }
})
</script>

<style scoped>
.container {
  max-width: 760px;
  margin: 24px auto;
  padding: 0 16px;
}
.item {
  margin-bottom: 12px;
}
.title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  text-decoration: none;
}
.title:hover {
  color: #409eff;
}
</style>
