import http from './http'

export const api = {
  // 认证
  register: (data) => http.post('/auth/register', data),
  login: (data) => http.post('/auth/login', data),

  // 用户
  getProfile: (id) => http.get(`/users/${id}`),
  updateProfile: (data) => http.put('/users/me', data),

  // 文章
  listArticles: (params) => http.get('/articles', { params }),
  getArticle: (id) => http.get(`/articles/${id}`),
  createArticle: (data) => http.post('/articles', data),

  // 问答
  listQuestions: (params) => http.get('/questions', { params }),
  getQuestion: (id) => http.get(`/questions/${id}`),
  createQuestion: (data) => http.post('/questions', data),
  createAnswer: (id, data) => http.post(`/questions/${id}/answers`, data),
  listAnswers: (id) => http.get(`/questions/${id}/answers`),
  acceptAnswer: (id, data) => http.post(`/questions/${id}/accept`, data),

  // 互动
  like: (data) => http.post('/interactions/like', data),
  unlike: (data) => http.post('/interactions/unlike', data),
  collect: (data) => http.post('/interactions/collect', data),
  uncollect: (data) => http.post('/interactions/uncollect', data),
  view: (data) => http.post('/interactions/view', data),
  listComments: (params) => http.get('/comments', { params }),
  createComment: (data) => http.post('/comments', data),

  // 社交
  followUser: (id) => http.post(`/users/${id}/follow`),
  unfollowUser: (id) => http.delete(`/users/${id}/follow`),
  followTag: (id) => http.post(`/tags/${id}/follow`),
  unfollowTag: (id) => http.delete(`/tags/${id}/follow`),

  // 搜索 / 信息流 / 通知
  search: (params) => http.get('/search', { params }),
  getFeed: (params) => http.get('/feed', { params }),
  listNotifications: (params) => http.get('/notifications', { params }),
  unreadCount: () => http.get('/notifications/unread-count'),
  markRead: (id) => http.post(`/notifications/${id}/read`),
  markAllRead: () => http.post('/notifications/read-all'),

  // 后台管理
  adminListUsers: (params) => http.get('/admin/users', { params }),
  adminBanUser: (id) => http.post(`/admin/users/${id}/ban`),
  adminListReports: (params) => http.get('/admin/reports', { params }),
  adminHandleReport: (id, data) => http.post(`/admin/reports/${id}/handle`, data),
  adminListTags: () => http.get('/admin/tags'),
  adminListSensitiveWords: () => http.get('/admin/sensitive-words'),
  adminCreateSensitiveWord: (data) => http.post('/admin/sensitive-words', data),
  adminListAuditLogs: (params) => http.get('/admin/audit-logs', { params }),
  adminStatistics: () => http.get('/admin/statistics')
}
