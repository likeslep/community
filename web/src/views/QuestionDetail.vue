<template>
  <div>
    <Navbar />
    <div class="container">
      <h1>{{ question.title }}</h1>
      <div class="meta">提问者 ID: {{ question.author_id }}</div>
      <pre class="content">{{ question.content }}</pre>

      <el-divider>回答</el-divider>
      <div v-if="answers.length === 0" class="empty">暂无回答</div>
      <el-card v-for="a in answers" :key="a.id" class="answer" :class="{ accepted: a.accepted }" shadow="hover">
        <div v-if="a.accepted" class="accepted-tag">✓ 已采纳</div>
        <pre class="answer-content">{{ a.content }}</pre>
        <el-button v-if="!a.accepted" size="small" @click="accept(a.id)">采纳</el-button>
      </el-card>

      <div class="answer-form">
        <el-input v-model="answerText" type="textarea" :rows="4" placeholder="写下你的回答..." />
        <el-button type="primary" size="small" @click="addAnswer" style="margin-top: 8px">提交回答</el-button>
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
const question = ref({})
const answers = ref([])
const answerText = ref('')

onMounted(async () => {
  await load()
})

async function load() {
  try {
    const [q, a] = await Promise.all([
      api.getQuestion(id),
      api.listAnswers(id)
    ])
    question.value = q.data
    answers.value = a.data.answers || []
  } catch (e) {
    ElMessage.error('加载失败')
  }
}

async function addAnswer() {
  if (!answerText.value.trim()) return
  await api.createAnswer(id, { content: answerText.value })
  answerText.value = ''
  await load()
}

async function accept(answerId) {
  await api.acceptAnswer(id, { answer_id: answerId })
  ElMessage.success('已采纳')
  await load()
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
.content {
  white-space: pre-wrap;
  word-break: break-word;
  background: #fafafa;
  padding: 16px;
  border-radius: 4px;
  font-family: inherit;
}
.answer {
  margin-bottom: 12px;
}
.answer.accepted {
  border-color: #67c23a;
}
.accepted-tag {
  color: #67c23a;
  font-weight: 600;
  margin-bottom: 8px;
}
.answer-content {
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  margin: 0;
}
.empty {
  color: #c0c4cc;
}
.answer-form {
  margin-top: 20px;
}
</style>
