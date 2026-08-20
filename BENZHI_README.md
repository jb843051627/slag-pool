# 渣池水冷循环监控系统

## 构建

```bash
docker build -t benzhi/slag-pool:latest -f benzhi.Dockerfile .
```

或使用构建脚本：

```bash
bash build_benzhi_docker.sh slag-pool linux/amd64
```

## 运行

```bash
docker run -d -p 8080:8080 -v slag-pool-data:/app/data benzhi/slag-pool:latest
```

访问 http://localhost:8080 查看监控仪表盘。
