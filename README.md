# NFT 交易系统

<div align=center>
   <img src="https://img.shields.io/badge/golang-1.20-blue"/>
   <img src="https://img.shields.io/badge/gin-1.10-blue"/>
   <img src="https://img.shields.io/badge/gorm-1.30.2-blue"/>
   <img src="https://img.shields.io/badge/vue-3.5.13-green"/>
   <img src="https://img.shields.io/badge/ant design-3.2.20-green"/>
</div>
<div align=center>
   <img src="https://img.shields.io/badge/hyperledger fabric-2.5.10-yellow"/>
   <img src="https://img.shields.io/badge/fabric-gateway-1.7.0-red"/>
</div>


## 简介

### 技术栈与功能

基于 Hyperledger Fabric 框架开发的 NFT(Non-Fungible Token)交易平台，前端框架是 Vue(TS)，后端框架是 gin，使用 JWT 鉴权，该系统供三个组织使用和维护：

1. 平台运营方
2. NFT 创作者
3. 金融机构

支持以下主要功能：

1. 钱包：管理平台代币，支持查询余额、转账和预扣款记录，此外，金融机构可以铸币
2. NFT 上传和查询（仅NFT 创作者）
3. NFT 交易：然后在市场中进行交易，交易分为以下两种：
   1. 拍卖
   2. 正常交易（允许议价）
4. 聊天，与卖家议价或与其他买家协商

### 项目背景

在数字内容爆发式增长的时代，数字艺术品、收藏品、资产的确权与流通成为了一个巨大的挑战。传统的中心化平台存在单点故障、交易不透明、佣金高昂、以及作品来源难以追溯等问题，因此我们希望构建一个去中心化、可信、高效且功能丰富的NFT交易平台。

该平台不仅是一个简单的买卖市场，更是一个为数字权益提供全生命周期管理的生态系统。通过引入平台运营方、NFT创作者和金融机构三大核心组织，构建了一个权责清晰、相互制衡、共同维护的信任联盟，彻底解决数字交易中的信任、价值和流动性问题。

### 为什么适合使用 Hyperledger Fabric

- 许可性与身份管理 (MSP)：Hyperledger Fabric 是联盟链，三大组织的成员必须经过实名认证和授权才能加入网络。
- 性能与可扩展性：NFT 拍卖可能引发高频出价，市场浏览需要快速查询。公有链的低 TPS 和高 Gas 费无法满足商用需求。
- 数据隐私与保密性：虽然交易最终状态是公开的，但过程中的一些信息需要隐私保护。
- NFT 的根本要求：分布式账本记录了每一个NFT从创建、每一次所有权变更的全部历史，提供了一个唯一、不可篡改的真相来源，完美解决了数字资产的溯源和确权问题。

## 本地开发指南

### 本地开发环境要求

- Go 1.23+
- Node.js 20+
- npm 9+
- Docker
- Docker Compose

### 1. 拉取项目（或手动下载）

```bash
git clone git@github.com:cagoya/BlockChain.git
```

### 2. 设置脚本权限

```bash
cd fabric-realty
find . -name "*.sh" -exec chmod +x {} \;
```

### 3. 启动区块链网络

首先需要启动基础的区块链网络环境（注意是进入到 network 目录执行）：

```bash
# 启动区块链网络（仅启动网络，不启动应用）
cd network
./install.sh
```

### 4. 启动后端服务

后端服务需要在本地编译运行，这样可以实时修改代码：

```bash
# 进入后端目录
cd application/server

# 运行后端服务
go run main.go
```

后端服务默认运行在 8888 端口。

### 5. 启动前端服务

前端服务同样需要在本地编译运行：

```bash
# 进入前端目录
cd application/web

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端开发服务器默认运行在 5173 端口。

### 6. 访问前端服务

前端`application/web/api/index.tx`中可以设置后端的地址，本地开发设置成 localhost，如需局域网内访问可以设置成局域网 ip，前端访问 http://localhost:5173

## 注意事项

1. 后端代码修改后，需要手动重启 `go run main.go`
2. 前端代码修改后，Vite 会自动热更新，无需手动重启
3. 区块链网络的修改（如链码更新）需要重新部署区块链网络
