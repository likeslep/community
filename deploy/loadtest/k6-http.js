// k6 负载测试脚本：HTTP 注册 + 登录 + 查询（Phase 14）。
// 用法：k6 run deploy/loadtest/k6-http.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // 爬升到 100 VU
    { duration: '1m', target: 100 },  // 保持
    { duration: '30s', target: 0 },   // 下降
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% 请求 < 500ms
    http_req_failed: ['rate<0.01'],   // 错误率 < 1%
  },
};

const BASE = 'http://localhost:8080/api/v1';

export default function () {
  // 注册
  const uniq = `${__VU}-${__ITER}-${Date.now()}`;
  const reg = http.post(`${BASE}/auth/register`, JSON.stringify({
    username: `user${uniq}`, email: `u${uniq}@example.com`, password: 'password123',
  }), { headers: { 'Content-Type': 'application/json' } });
  check(reg, { 'register 2xx': (r) => r.status >= 200 && r.status < 300 });

  // 登录
  const login = http.post(`${BASE}/auth/login`, JSON.stringify({
    username: `user${uniq}`, password: 'password123',
  }), { headers: { 'Content-Type': 'application/json' } });
  check(login, { 'login 2xx': (r) => r.status >= 200 && r.status < 300 });
}
