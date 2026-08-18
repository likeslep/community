// k6 负载测试脚本：Feed 读取（Phase 14）。
// 用法：k6 run deploy/loadtest/k6-feed.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 200 },
    { duration: '1m', target: 200 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    http_req_failed: ['rate<0.01'],
  },
};

// 需要先通过登录获取 token，此处用占位 token 演示读取路径。
const BASE = 'http://localhost:8080/api/v1';
const TOKEN = __ENV.TOKEN || '';

export default function () {
  const res = http.get(`${BASE}/feed`, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });
  check(res, { 'feed 2xx': (r) => r.status >= 200 && r.status < 300 });
}
