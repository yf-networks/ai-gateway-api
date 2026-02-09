# Docker 镜像构建与推送

## 前置条件

- 已安装 Docker（建议较新版本）
- 需要多架构推送时：Docker Buildx 可用（Docker Desktop 默认已带）
- 推送镜像前请先登录镜像仓库：例如 `docker login <registry>`

## 镜像内目录结构与挂载配置

镜像运行时的工作目录为：

- `/home/work/api-server`

该目录下的关键文件/目录如下：

- `/home/work/api-server/api-server`：服务二进制
- `/home/work/api-server/conf/`：配置目录（默认会 COPY 进镜像，可通过 volume 覆盖）
- `/home/work/api-server/static/`：静态资源目录（dashboard 解压后在这里）
- `/home/work/api-server/log/`：日志目录（启动参数 `-l ./log`）

### 挂载配置（推荐）

你可以通过挂载本地 `conf` 目录来覆盖镜像内默认配置。例如：

```bash
docker run -d \
	--name ai-gateway-api \
	-p 8183:8183 \
	-v $(pwd)/conf:/home/work/api-server/conf \
	ai-gateway-api:latest
```

## 本地构建（make docker）

Makefile 提供 `make docker` 目标用于本地构建。

### 默认行为

- 镜像名：`ai-gateway-api`
- 默认会生成两份本地 tag：
    - `ai-gateway-api:v<Version>`
    - `ai-gateway-api:latest`

### 常用参数

可通过环境变量覆盖：

- DASHBOARD_VERSION : 控制面版本，来自 `yf-networks/ai-gateway-web` 的 release 包

示例：

```bash
make docker DASHBOARD_VERSION=v0.0.1'
```

## 多架构构建并推送（make docker-push）

Makefile 提供 `make docker-push` 目标用于 **buildx 多架构** 构建并推送到远端仓库。

### 必填参数

- `REGISTRY`：镜像仓库前缀（必填），例如：
	- `ghcr.io/<org>`
	- `docker.io/<namespace>`
	- `registry.example.com/<project>`

推送的远端 tag：

- `<REGISTRY>/<BFE_IMAGE_NAME>:v<Version>`
- `<REGISTRY>/<BFE_IMAGE_NAME>:latest`

### 可选参数

- `PLATFORMS`：多架构平台（默认 `linux/amd64,linux/arm64`）

示例：

```bash
# 指定平台
make docker-push \
	REGISTRY=ghcr.io/cc14514 \
	PLATFORMS=linux/amd64
```
