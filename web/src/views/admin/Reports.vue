<template>
  <div>
    <h2>举报处理</h2>
    <el-table :data="reports" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="reporter_id" label="举报人" width="100" />
      <el-table-column prop="target_type" label="目标类型" width="120" />
      <el-table-column prop="target_id" label="目标 ID" width="100" />
      <el-table-column prop="reason" label="原因" />
      <el-table-column prop="status" label="状态" width="100" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <template v-if="row.status === 'pending'">
            <el-button size="small" type="success" @click="handle(row.id, 'approved')">通过</el-button>
            <el-button size="small" type="danger" @click="handle(row.id, 'rejected')">驳回</el-button>
          </template>
          <span v-else style="color: #909399">已处理</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'

const reports = ref([])

onMounted(load)

async function load() {
  const res = await api.adminListReports({ limit: 100 })
  reports.value = res.data.reports || []
}

async function handle(id, action) {
  await api.adminHandleReport(id, { action })
  ElMessage.success('已处理')
  await load()
}
</script>
