import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', component: () => import('../views/Login.vue') },
  { path: '/register', component: () => import('../views/Register.vue') },
  { path: '/', component: () => import('../views/Home.vue') },
  { path: '/articles/new', component: () => import('../views/ArticleCreate.vue') },
  { path: '/articles/:id', component: () => import('../views/ArticleDetail.vue') },
  { path: '/questions/new', component: () => import('../views/QuestionCreate.vue') },
  { path: '/questions/:id', component: () => import('../views/QuestionDetail.vue') },
  { path: '/search', component: () => import('../views/Search.vue') },
  { path: '/profile', component: () => import('../views/Profile.vue') },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    children: [
      { path: '', component: () => import('../views/admin/Dashboard.vue') },
      { path: 'users', component: () => import('../views/admin/Users.vue') },
      { path: 'reports', component: () => import('../views/admin/Reports.vue') },
      { path: 'tags', component: () => import('../views/admin/Tags.vue') },
      { path: 'sensitive-words', component: () => import('../views/admin/SensitiveWords.vue') },
      { path: 'audit-logs', component: () => import('../views/admin/AuditLogs.vue') }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
