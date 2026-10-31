# Livery API

使用 Go 编写的涂装列表 API，默认使用 SQLite 数据库。服务启动时会自动创建 `liveries` 表。

## 启动

```bash
go mod tidy
go run .
```

默认监听 `http://localhost:8080`，数据库文件为当前目录下的 `liveries.db`。

可通过环境变量配置：

- `PORT`：监听端口，默认 `8080`
- `DATABASE_URL`：SQLite 连接地址，默认 `file:liveries.db`

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
const response = await fetch('http://localhost:8080/api/v1/liveries?page=1&page_size=20');
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
