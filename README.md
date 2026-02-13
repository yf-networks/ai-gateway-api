[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev/)

[English](README.md) | [简体中文](README_CN.md)

# AI Gateway API

AI Gateway API is the **control-plane core component** of AI Gateway. It provides APIs for creating, storing, and distributing policies/configurations. Built on top of the open-source [BFE](https://github.com/bfenetworks/bfe) project, it offers unified traffic gateway management for AI scenarios.

## Architecture Overview

![Architecture](/docs/zh_cn/assert/deploy_architecture_ai.png)

AI Gateway consists of the following core components:

| Component | Role | Description | Repository |
|---|---|---|---|
| **AI Gateway API** | Control plane | Exposes Open APIs to manage policies/configurations and distribute them | This repo |
| **Dashboard** | Admin console | Web UI for visual management | [yf-networks/ai-gateway-web](https://github.com/yf-networks/ai-gateway-web) |
| **BFE** | Data plane | Traffic forwarding and access control | [bfenetworks/bfe](https://github.com/bfenetworks/bfe) |
| **Conf Agent** | Config agent | Fetches the latest config and triggers BFE hot reload | [bfenetworks/conf-agent](https://github.com/bfenetworks/conf-agent) |
| **Service Controller** | Service discovery | Discovers and syncs Kubernetes backend services | [bfenetworks/service-controller](https://github.com/bfenetworks/service-controller) |

## Key Features

- **AI route management**: Route configuration for multiple AI model providers (OpenAI, DeepSeek, Anthropic, Google Gemini, etc.)
- **API key management**: Create/delete/validate API keys for AI services
- **Domain management**: Bind domains and configure routing rules
- **Certificate management**: Upload and manage TLS certificates
- **Cluster/sub-cluster management**: Manage backend service clusters
- **Traffic management**: Traffic allocation and scheduling
- **Dashboard integration**: Container image bundles the Web console UI
- **Config export**: Export configuration for BFE data plane and Conf Agent

## Quick Start

### Prerequisites

- Go 1.22 or later
- MySQL 8
- Redis 6.2

### Build from source

```bash
# Clone repo
git clone https://github.com/yf-networks/ai-gateway-api.git
cd ai-gateway-api

# Download deps and build
make

# Artifacts are in output/
```

### Initialize database

```bash
mysql -u{user} -p{password} < db_ddl.sql
```

### Configure

Edit `conf/ai_gateway_api.toml`. At minimum, update the database connection:

```toml
[Databases.bfe_db]
DBName               = "open_bfe"
Addr                 = "127.0.0.1:3306"
User                 = "{user}"
Passwd               = "{password}"
```

For detailed configuration parameters, see [配置文件说明](/docs/zh_cn/config_param.md).

### Start the service

```bash
./ai-gateway-api -c ./conf -sc ai_gateway_api.toml -l ./log
```

After startup:
- **API port**: `8183` (default)
- **Monitoring port**: `8284` (default)
- **Dashboard**: open `http://localhost:8183` (default username/password: `admin` / `admin`)

## Container Images & Kubernetes Deployment Example

See the [ai-gateway-demo](https://github.com/yf-networks/ai-gateway-demo/tree/main) repository:
- https://github.com/yf-networks/ai-gateway-demo

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow and guidelines.

## License

AI Gateway API is released under the [Apache License 2.0](LICENSE).
