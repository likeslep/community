<template>
  <div>
    <Navbar />
    <div class="container">
      <h1>{{ article.title }}</h1>
      <div class="meta">作者 ID: {{ article.author_id }} · 状态: {{ article.status }}</div>
      <div class="actions">
        <el-button size="small" @click="like">👍 点赞</el-button>
        <el-button size="small" @click="collect">⭐ 收藏</el-button>
      </div>
      <pre class="content">{{ article.content }}</pre>

      <el-divider>评论</el-divider>
      <div v-if="comments.length === 0" class="empty">暂无评论</div>
      <div v-for="c in comments" :key="c.id" class="comment">
        <span class="comment-content">{{ c.content }}</span>
      </div>
      <div class="comment-form">
        <el-input v-model="commentText" placeholder="写下你的评论..." />
        <el-button type="primary" size="small" @click="addComment" style="margin-top: 8px">
          发表评论
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from '../components/Navbar.vue'
import { api } from '../api'

const route = useRoute()
const id = route.params.id
const article = ref({})
const comments = ref([])
const commentText = ref('')

onMounted(async () => {
  try {
    const [a, c] = await Promise.all([
      api.getArticle(id),
      api.listComments({ target_type: 'article', target_id: id })
    ])
    article.value = a.data
    comments.value = c.data.comments || []
    api.view({ target_type: 'article', target_id: Number(id) }).catch(() => {})
  } catch (e) {
    ElMessage.error('加载失败')
  }
})

async function like() {
  await api.like({ target_type: 'article', target_id: Number(id) })
  ElMessage.success('已点赞')
}
async function collect() {
  await api.collect({ target_type: 'article', target_id: Number(id) })
  ElMessage.success('已收藏')
}
async function addComment() {
  if (!commentText.value.trim()) return
  await api.createComment({ target_type: 'article', target_id: Number(id), content: commentText.value })
  commentText.value = ''
  const c = await api.listComments({ target_type: 'article', target_id: id })
  comments.value = c.data.comments || []
}
</script>

<style scoped>
.container {
  max-width: 760px;
  margin: 24px auto;
  padding: 0 16px;
}
.meta {
  color: #909399;
  font-size: 14px;
  margin-bottom: 12px;
}
.actions {
  margin-bottom: 16px;
}
.content {
  white-space: pre-wrap;
  word-break: break-word;
  background: #fafafa;
  padding: 16px;
  border-radius: 4px;
  font-family: inherit;
}
.comment {
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}
.empty {
  color: #c0c4cc;
}
.comment-form {
  margin-top: 16px;
}
</style>
