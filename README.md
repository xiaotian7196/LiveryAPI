# Livery API

使用 Go 编写的涂装列表 API，默认使用 SQLite 数据库。服务启动时会自动创建 `liveries` 表。

## 启动

```bash
go mod tidy
go run .
```

服务对外访问地址为 `https://api.xn--3ckue.jp/`。Go 程序默认监听本机 `8080` 端口，建议通过 Nginx/Caddy 反向代理并配置 HTTPS。数据库文件默认为当前目录下的 `liveries.db`。

可通过环境变量配置：

- `PORT`：监听端口，默认 `8080`
- `DATABASE_URL`：SQLite 连接地址，默认 `file:liveries.db`

API GET 请求限制：同一个客户端 IP 在 60 秒内最多请求 30 次 `/api/` 接口。超过限制时返回 HTTP `429 Too Many Requests`，并通过 `Retry-After` 返回建议等待秒数。服务部署在 Cloudflare 后会优先使用 `CF-Connecting-IP` 识别真实访客 IP。`/health` 健康检查不受此限制。

## API

### 获取涂装列表

```http
GET /api/v1/liveries?page=1&page_size=20&search=737&game=msfs
```

返回示例：

```json
{
  "items": [
    {
      "id": 1,
      "name": "示例航空公司涂装",
      "game": "msfs",
      "aircraft": "Boeing 737",
      "download_url_cn": "https://cn.example.com/liveries/1.zip",
      "download_url_global": "https://global.example.com/liveries/1.zip"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 1,
  "total_pages": 1
}
```

### 获取单个涂装

```http
GET /api/v1/liveries/{id}
```

### 健康检查

```http
GET /health
```

## 写入测试数据

服务首次启动后，可以使用 SQLite 客户端写入数据：

```sql
INSERT INTO liveries (name, description, game, aircraft, image_url, download_url_cn, download_url_global, version)
VALUES ('示例涂装', '涂装说明', 'msfs', 'A320',
        'https://example.com/preview.jpg',
        'https://cn.example.com/livery.zip',
        'https://global.example.com/livery.zip', '1.0.0');
```

## 前端请求示例

```js
const response = await fetch('https://api.xn--3ckue.jp/api/v1/liveries?page=1&page_size=20');
const data = await response.json();
console.log(data.items[0].download_url_cn);
console.log(data.items[0].download_url_global);
```

## 从网站导入涂装

项目包含网站导入工具，会将 `https://blog.koten.top/` 中整理出的涂装写入本地 SQLite 数据库：

```bash
go run ./cmd/importer
```

工具支持重复执行，不会重复插入相同的涂装记录。网站只有中国下载地址的记录会将 `download_url_global` 留空。

## Debian 安装包

项目提供 Debian/systemd 打包配置。推荐在 Debian 或 Ubuntu 上执行：

```bash
chmod +x scripts/build-deb.sh
./scripts/build-deb.sh
sudo dpkg -i dist/liveryapi_1.0.0_amd64.deb
```

安装后服务会自动启动，数据库保存在：

```text
/var/lib/liveryapi/liveries.db
```

常用命令：

```bash
sudo systemctl status liveryapi
sudo systemctl restart liveryapi
sudo journalctl -u liveryapi -f
```

对外接口地址：

```text
https://api.xn--3ckue.jp/
```

如果使用 Nginx，可以将域名转发到本机的 Go 服务：

```nginx
server {
    listen 443 ssl http2;
    server_name api.xn--3ckue.jp;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

如果只需要生成 Linux 二进制文件，也可以执行：

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o liveryapi .
```
