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
   <img src="https://img.shields.io/badge/fabric chaincode go-v2.0.0-red"/>
   <img src="https://img.shields.io/badge/fabric contract api go-v2.0.0-red"/>
</div>


## 简介

基于 Hyperledger Fabric 框架开发的 NFT(Non-Fungible Token)交易平台，前端框架是 Vue(TS)，后端框架是 gin，使用 JWT 鉴权，该系统供三个组织使用和维护：

1. 平台运营方
2. NFT 创作者
3. 金融机构

支持以下主要功能：

1. 钱包：管理平台代币，支持查询余额和转账，金融机构可以铸币
2. NFT 上传和查询
3. NFT 交易：然后在市场中进行交易，交易分为以下两种：
   1. 拍卖
   2. 正常交易（允许议价）