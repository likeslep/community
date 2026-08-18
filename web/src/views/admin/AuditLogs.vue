<template>
  <div>
    <h2>审计日志</h2>
    <el-table :data="logs" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="admin_id" label="管理员" width="100" />
      <el-table-column prop="action" label="操作" width="180" />
      <el-table-column prop="target_type" label="目标类型" width="120" />
      <el-table-column prop="target_id" label="目标 ID" width="100" />
      <el-table-column prop="detail" label="详情" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../../api'

const logs = ref([])

onMounted(async () => {
  const res = await api.adminListAuditLogs({ limit: 100 })
  logs.value = res.data.logs || []
})
</script>
