# Community Web 前端

Vue3 + Vite + JavaScript + Element Plus，包含**用户客户端**与**后台管理界面**。

## 技术栈

- Vue 3 + Vite
- Vue Router + Pinia
- Element Plus
- Axios

## 启动

```bash
npm install
npm run dev    # 开发模式 http://localhost:5173
```

开发模式下 `/api` 请求代理到 gateway（`http://localhost:8080`），见 `vite.config.js`。

## 构建

```bash
npm run build   # 产物输出到 dist/
```

## 目录结构

```
src/
├── api/          axios 实例 + 接口封装
├── store/        Pinia（auth）
├── router/       路由
├── components/   共享组件（Navbar）
└── views/        页面
    ├── *.vue            客户端页面（登录/注册/首页/文章/问答/搜索/个人）
    └── admin/           后台管理（统计/用户/举报/标签/敏感词/审计）
```

## 页面路由

| 路由 | 说明 |
|------|------|
| `/login` `/register` | 登录 / 注册 |
| `/` | 首页（最新文章/问题） |
| `/articles/new` `/articles/:id` | 写文章 / 文章详情 |
| `/questions/new` `/questions/:id` | 提问 / 问题详情 |
| `/search` | 搜索 |
| `/profile` | 个人资料 |
| `/admin/*` | 后台管理（需 admin 角色） |
