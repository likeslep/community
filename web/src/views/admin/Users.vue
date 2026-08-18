<template>
  <div>
    <h2>用户管理</h2>
    <el-table :data="users" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="role" label="角色" width="120" />
      <el-table-column prop="status" label="状态" width="120" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button v-if="row.status !== 'banned'" size="small" type="danger" @click="ban(row.id)">
            封禁
          </el-button>
          <span v-else style="color: #909399">已封禁</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../../api'

const users = ref([])

onMounted(load)

async function load() {
  const res = await api.adminListUsers({ limit: 100 })
  users.value = res.data.users || []
}

async function ban(id) {
  await api.adminBanUser(id)
  ElMessage.success('已封禁')
  await load()
}
</script>
