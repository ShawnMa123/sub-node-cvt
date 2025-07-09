# sub-node-cvt

[English](./README.en.md) | **中文**

一个轻量、无状态、保护隐私的 Clash/Clash-Meta 订阅链接生成器。它允许用户通过一个简单的 Web 界面，将节点配置（YAML 格式）转换为功能强大的订阅链接，并支持通过 GitHub Gist 进行私密、安全的托管。

项目核心设计理念是**无状态与隐私优先**。当使用 URL 模式时，所有配置信息均编码在链接中；当使用 Gist 模式时，配置会安全地存储在用户自己的私密 Gist 中。后端服务在任何情况下都不存储用户数据。

![alt text](image-1.png)

---

## ✨ 核心特性

-   **双模式生成**:
    -   **URL 模式**: 将所有配置编码进一个超长 URL，简单直接，无需认证。
    -   **Gist 模式**: 通过 GitHub OAuth 授权，将生成的配置安全地存入用户自己的**私密 Gist**，生成一个简短、固定的订阅链接。配置更新只需修改 Gist 内容即可。
-   **隐私与安全**:
    -   后端服务完全无状态，不记录任何用户信息或配置。
    -   Gist 创建通过标准 OAuth 2.0 流程，应用仅获取创建 Gist 的最小权限。
-   **强大的配置能力**:
    -   支持 Clash Meta `proxies` YAML 格式的直接输入。
    -   支持中转链 (Relay / 链式代理) 配置。
    -   支持可插拔的规则集（如去广告、GFW 列表等）。
-   **轻量与高性能**:
    -   后端由 Go 编写，性能卓越，资源占用极低。
    -   前端使用 Vue 3 构建，界面响应迅速。
-   **多种部署方式**: 支持本地运行、传统服务器部署以及通过 Docker 进行容器化部署。

## 🛠️ 技术栈

-   **后端**: Go (Golang)
    -   Web 框架: `net/http`
    -   GitHub 集成: `golang.org/x/oauth2`, `github.com/google/go-github`
-   **前端**: Vue 3 (Composition API) & `js-yaml`
-   **部署**: Docker, Docker Compose, systemd

## 🚀 快速开始 (本地开发)

### 1. 环境要求

-   Go 1.23 或更高版本
-   一个 GitHub OAuth Application 用于 Gist 功能。
    -   在 [GitHub 开发者设置](https://github.com/settings/developers) 中创建一个 OAuth App。
    -   **Homepage URL**: `http://localhost:8080`
    -   **Authorization callback URL**: `http://localhost:8080/auth/github/callback`
    -   创建后，获取 `Client ID` 和 `Client Secret`。

### 2. 克隆仓库

```bash
git clone https://github.com/ShawnMa123/sub-node-cvt.git
cd sub-node-cvt
```

### 3. 配置环境变量

为了在本地运行，您需要设置环境变量。最简单的方式是在一行命令中完成：

```bash
# 将 "xxx" 替换为您自己的真实 ID 和 Secret
GITHUB_CLIENT_ID="xxx" GITHUB_CLIENT_SECRET="xxx" go run main.go
```
程序将在 `http://localhost:8080` 启动。

### 4. 如何使用

1.  打开 `http://localhost:8080`。
2.  **(可选) 登录 GitHub**: 点击 "使用 GitHub 登录" 以启用 Gist 功能。
3.  **粘贴节点**: 将你的 Clash `proxies` 列表粘贴到文本框中。
4.  **选择规则与配置中转**。
5.  **生成链接**:
    -   点击 **"生成预览链接"** 会得到一个包含所有信息的长 URL。
    -   如果已登录，点击 **"保存到私密 Gist"** 会将配置上传到您的 GitHub，并得到一个简短、稳定的 Gist raw URL 作为订阅链接。

## 部署方案

### 方案一: 使用 Docker 和 Docker Compose (推荐)

这是最简单、最便携的部署方式。

1.  **创建 `.env` 文件**:
    在项目根目录创建 `.env` 文件，并填入您的密钥：
    ```env
    GITHUB_CLIENT_ID="iv1.xxxxxxxxxxxxxxxx"
    GITHUB_CLIENT_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    ```
    **重要**: 将 `.env` 文件添加到你的 `.gitignore` 中！

2.  **启动服务**:
    确保您已安装 Docker 和 Docker Compose，然后在项目根目录运行：
    ```bash
    docker-compose up -d
    ```
    服务将在 `http://<您的服务器IP>:8080` 上运行。

### 方案二: 在 Linux 服务器上使用 `systemd`

这是一种稳定、可靠的传统部署方式。

1.  **编译应用**:
    为您的 Linux 服务器编译一个静态二进制文件：
    ```bash
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o sub-node-cvt main.go
    ```

2.  **上传文件**:
    将以下文件和目录上传到服务器的同一目录下（例如 `/opt/sub-node-cvt/`）：
    -   `sub-node-cvt` (编译好的二进制文件)
    -   `frontend/`
    -   `templates/`
    -   `rulesets/`

3.  **配置 `systemd` 服务**:
    -   在服务器上创建环境变量文件 `/etc/sub-node-cvt/config.env` 并填入您的密钥。
    -   创建 `systemd` 单元文件 `/etc/systemd/system/sub-node-cvt.service`，内容如下：
    ```ini
    [Unit]
    Description=Subscription Node Converter
    After=network.target

    [Service]
    User=your_username # 替换为您的用户名
    WorkingDirectory=/opt/sub-node-cvt # 替换为您的应用目录
    EnvironmentFile=/etc/sub-node-cvt/config.env
    ExecStart=/opt/sub-node-cvt/sub-node-cvt
    Restart=always

    [Install]
    WantedBy=multi-user.target
    ```

4.  **启动服务**:
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable --now sub-node-cvt
    ```

## 💡 未来展望

-   [ ] 支持更多类型的节点格式输入（如 V2RayN 分享链接）。
-   [ ] 允许用户更新已存在的 Gist 而不是每次都创建新的。
-   [ ] 适配 Cloudflare Workers 部署。
