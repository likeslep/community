<template>
  <div>
    <h2>数据统计</h2>
    <el-row :gutter="16">
      <el-col :span="6" v-for="s in stats" :key="s.label">
        <el-card shadow="hover">
          <div class="num">{{ s.value }}</div>
          <div class="label">{{ s.label }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../../api'

const stats = ref([])

onMounted(async () => {
  try {
    const res = await api.adminStatistics()
    const d = res.data
    stats.value = [
      { label: '用户数', value: d.user_count },
      { label: '举报数', value: d.report_count },
      { label: '标签数', value: d.tag_count },
      { label: '敏感词数', value: d.sensitive_word_count },
      { label: '审计日志', value: d.audit_log_count }
    ]
  } catch (e) {
    console.error(e)
  }
})
</script>

<style scoped>
.num {
  font-size: 32px;
  font-weight: 700;
  color: #409eff;
}
.label {
  color: #909399;
  margin-top: 8px;
}
</style>
