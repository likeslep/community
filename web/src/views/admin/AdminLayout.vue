<template>
  <el-container class="admin-layout">
    <el-aside width="200px">
      <div class="logo">后台管理</div>
      <el-menu router :default-active="$route.path" class="menu">
        <el-menu-item index="/admin">数据统计</el-menu-item>
        <el-menu-item index="/admin/users">用户管理</el-menu-item>
        <el-menu-item index="/admin/reports">举报处理</el-menu-item>
        <el-menu-item index="/admin/tags">标签管理</el-menu-item>
        <el-menu-item index="/admin/sensitive-words">敏感词管理</el-menu-item>
        <el-menu-item index="/admin/audit-logs">审计日志</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <router-link to="/">返回前台</router-link>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.admin-layout {
  height: 100vh;
}
.logo {
  height: 56px;
  line-height: 56px;
  text-align: center;
  font-weight: 700;
  color: #409eff;
}
.menu {
  border-right: none;
}
.header {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  border-bottom: 1px solid #e4e7ed;
}
</style>
