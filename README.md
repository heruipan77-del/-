# SmartOps (Go + React)

一个最小可用的“智能运维存储平台”示例：

- Go 后端提供统一 API（登录设备、查询 StoragePool、查询 LUN）
- React 前端提供连接配置和查询页面

## 目录

- `backend/` Go API 服务
- `web/` React 前端

## 启动

### 1) 构建前端

```bash
cd web
npm install
npm run build
```

### 2) 启动后端（同时托管前端静态文件）

```bash
cd backend
go run ./cmd/server
```

默认访问: `http://localhost:8080`

## API

- `POST /api/storagepools`
- `POST /api/lun`

请求体示例：

```json
{
  "ip": "10.10.10.10",
  "port": 8088,
  "deviceId": "210235G7J20000xxxx",
  "username": "admin",
  "password": "xxx",
  "insecure": true,
  "lunId": "1"
}
```
