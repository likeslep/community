<template>
  <div>
    <h2>标签管理</h2>
    <el-table :data="tags" border>
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="name" label="标签名" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../../api'

const tags = ref([])

onMounted(async () => {
  const res = await api.adminListTags()
  tags.value = res.data.tags || []
})
</script>
