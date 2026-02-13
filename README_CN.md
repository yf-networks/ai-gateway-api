[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev/)

[English](README.md) | 简体中文

# AI Gateway API

AI Gateway API 是 AI Gateway 的**控制面核心组件**，负责策略/配置的录入、存储和下发接口。基于 [BFE](https://github.com/bfenetworks/bfe) 开源项目，为 AI 场景提供统一的流量网关管理能力。

## 架构概述

![架构](/docs/zh_cn/assert/deploy_architecture_ai.png)

AI Gateway 包含如下核心组件：

| 组件 | 角色 | 说明 | 仓库 |
|------|------|------|------|
| **AI Gateway API** | 控制面 | 对外提供 Open API 接口，完成策略/配置的变更、存储和下发 | 本仓库 |
| **Dashboard** | 管理控制台 | Web 可视化管理界面 | [yf-networks/ai-gateway-web](https://github.com/yf-networks/ai-gateway-web) |
| **BFE** | 数据面 | 负责流量转发与接入控制 | [bfenetworks/bfe](https://github.com/bfenetworks/bfe) |
| **Conf Agent** | 配置代理 | 获取最新配置并触发 BFE 热加载 | [bfenetworks/conf-agent](https://github.com/bfenetworks/conf-agent) |
| **Service Controller** | 服务发现 | 发现并同步 Kubernetes 后端服务 | [bfenetworks/service-controller](https://github.com/bfenetworks/service-controller) |

## 主要功能

- **AI 路由管理**：支持多 AI 模型提供商（OpenAI、DeepSeek、Anthropic、Google Gemini 等）的路由配置
- **API Key 管理**：AI 服务的 API Key 创建、删除与校验
- **域名管理**：域名绑定与路由规则配置
- **证书管理**：TLS 证书的上传与管理
- **集群/子集群管理**：后端服务集群的配置管理
- **流量管理**：流量分配与调度
- **Dashboard 集成**：镜像自动打包 Web 管理界面
- **配置导出**：为 BFE 数据面和 Conf Agent 提供配置导出接口

## 快速开始

### 前置条件

- Go 1.22 或更高版本
- MySQL 8
- Redis 6.2

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/yf-networks/ai-gateway-api.git
cd ai-gateway-api

# 下载依赖并编译
make

# 编译产出在 output/ 目录
```

### 初始化数据库

```bash
mysql -u{user} -p{password} < db_ddl.sql
```

### 修改配置

编辑 `conf/ai_gateway_api.toml`，至少需要修改数据库连接信息：

```toml
[Databases.bfe_db]
DBName               = "open_bfe"
Addr                 = "127.0.0.1:3306"
User                 = "{user}"
Passwd               = "{password}"
```

详细配置参数说明参见 [配置文件说明](/docs/zh_cn/config_param.md)。

### 启动服务

```bash
./ai-gateway-api -c ./conf -sc ai_gateway_api.toml -l ./log
```

服务启动后：
- **API 服务端口**：`8183`（默认）
- **监控端口**：`8284`（默认）
- **Dashboard**：浏览器访问 `http://localhost:8183`（默认账号/密码：`admin` / `admin`）

## 容器镜像与 Kubernetes 部署示例

请参阅 [ai-gateway-demo](https://github.com/yf-networks/ai-gateway-demo/tree/main): 
- https://github.com/yf-networks/ai-gateway-demo

## 贡献

欢迎贡献代码！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解开发流程和规范。

## 许可证

AI Gateway API 基于 [Apache License 2.0](LICENSE) 发布。
