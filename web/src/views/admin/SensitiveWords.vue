<template>
  <div>
    <h2>敏感词管理</h2>
    <div class="add-form">
      <el-input v-model="form.word" placeholder="敏感词" style="width: 200px" />
      <el-select v-model="form.level" placeholder="级别" style="width: 140px; margin: 0 8px">
        <el-option label="拦截" value="block" />
        <el-option label="人工审核" value="review" />
      </el-select>
      <el-button type="primary" @click="add">添加</el-button>
    </div>
    <el-table :data="words" border style="margin-top: 16px">
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="word" label="敏感词" />
      <el-table-column prop="level" label="级别" width="120" />
    </el-table>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'

const words = ref([])
const form = reactive({ word: '', level: 'review' })

onMounted(load)

async function load() {
  const res = await api.adminListSensitiveWords()
  words.value = res.data.words || []
}

async function add() {
  if (!form.word.trim()) return
  await api.adminCreateSensitiveWord({ word: form.word, level: form.level })
  ElMessage.success('已添加')
  form.word = ''
  await load()
}
</script>

<style scoped>
.add-form {
  display: flex;
  align-items: center;
}
</style>
