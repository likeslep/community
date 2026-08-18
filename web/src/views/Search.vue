<template>
  <div>
    <Navbar />
    <div class="container">
      <div class="search-bar">
        <el-input v-model="keyword" placeholder="搜索文章/问题/回答..." @keyup.enter="doSearch" />
        <el-button type="primary" @click="doSearch">搜索</el-button>
      </div>
      <div class="result-meta">共 {{ total }} 条结果</div>
      <el-card v-for="r in results" :key="r.id" class="result" shadow="hover">
        <div class="result-type">{{ r.type }}</div>
        <div class="result-title">{{ r.title }}</div>
        <div v-if="r.snippet" class="result-snippet" v-html="r.snippet"></div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Navbar from '../components/Navbar.vue'
import { api } from '../api'

const keyword = ref('')
const results = ref([])
const total = ref(0)

async function doSearch() {
  if (!keyword.value.trim()) return
  const res = await api.search({ keyword: keyword.value, page: 1, page_size: 20 })
  results.value = res.data.results || []
  total.value = res.data.total || 0
}
</script>

<style scoped>
.container {
  max-width: 760px;
  margin: 24px auto;
  padding: 0 16px;
}
.search-bar {
  display: flex;
  gap: 8px;
}
.result-meta {
  color: #909399;
  font-size: 14px;
  margin: 16px 0;
}
.result {
  margin-bottom: 12px;
}
.result-type {
  color: #409eff;
  font-size: 12px;
}
.result-title {
  font-size: 16px;
  font-weight: 600;
  margin: 4px 0;
}
.result-snippet {
  color: #606266;
  font-size: 14px;
}
</style>
