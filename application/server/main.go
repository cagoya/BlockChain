package main

import (
	"application/api"
	"application/config"
	"application/middleware"
	"application/model"
	"application/pkg/fabric"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败：%v", err)
	}

	// 初始化 Fabric 客户端
	if err := fabric.InitFabric(); err != nil {
		log.Fatalf("初始化Fabric客户端失败：%v", err)
	}

	// 初始化数据库
	if err := model.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败：%v", err)
	}

	// 创建 Gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	apiGroup := r.Group("/api")

	// 注册路由
	accountHandler := api.NewAccountHandler()
	realtyAgencyHandler := api.NewRealtyAgencyHandler()
	tradingPlatformHandler := api.NewTradingPlatformHandler()
	bankHandler := api.NewBankHandler()

	// 创建JWT中间件
	jwtMiddleware, err := middleware.NewJWTMiddleware()
	if err != nil {
		log.Fatalf("创建JWT中间件失败：%v", err)
	}

	// 账号相关接口（无需认证）
	account := apiGroup.Group("/account")
	{
		// 用户注册
		account.POST("/register", accountHandler.Register)
		// 用户登录
		account.POST("/login", accountHandler.Login)
		// 用户登出
		account.POST("/logout", accountHandler.Logout)
	}

	// 需要认证的账号接口
	authAccount := apiGroup.Group("/account").Use(jwtMiddleware.Auth(1))
	{
		// 获取用户信息
		authAccount.GET("/profile", accountHandler.GetProfile)
		// 更新用户信息
		authAccount.PUT("/profile", accountHandler.UpdateProfile)
		// 获取头像
		authAccount.GET("/avatar", accountHandler.GetAvatar)
		// 更新头像
		authAccount.PUT("/avatar", accountHandler.UpdateAvatar)
	}

	// 不动产登记机构的接口
	realty := apiGroup.Group("/realty-agency")
	{
		// 创建房产信息
		realty.POST("/realty/create", realtyAgencyHandler.CreateRealEstate)
		// 查询房产接口
		realty.GET("/realty/:id", realtyAgencyHandler.QueryRealEstate)
		realty.GET("/realty/list", realtyAgencyHandler.QueryRealEstateList)
		// 查询区块接口
		realty.GET("/block/list", realtyAgencyHandler.QueryBlockList)
	}

	// 交易平台的接口
	trading := apiGroup.Group("/trading-platform")
	{
		// 生成交易
		trading.POST("/transaction/create", tradingPlatformHandler.CreateTransaction)
		// 查询房产接口
		trading.GET("/realty/:id", tradingPlatformHandler.QueryRealEstate)
		// 查询交易接口
		trading.GET("/transaction/:txId", tradingPlatformHandler.QueryTransaction)
		trading.GET("/transaction/list", tradingPlatformHandler.QueryTransactionList)
		// 查询区块接口
		trading.GET("/block/list", tradingPlatformHandler.QueryBlockList)
	}

	// 银行的接口
	bank := apiGroup.Group("/bank")
	{
		// 完成交易
		bank.POST("/transaction/complete/:txId", bankHandler.CompleteTransaction)
		// 查询交易接口
		bank.GET("/transaction/:txId", bankHandler.QueryTransaction)
		bank.GET("/transaction/list", bankHandler.QueryTransactionList)
		// 查询区块接口
		bank.GET("/block/list", bankHandler.QueryBlockList)
	}

	// 打印路由信息
	printRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	log.Printf("服务器启动在端口: %d", config.GlobalConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败：%v", err)
	}
}

// printRoutes 打印所有注册的路由
func printRoutes(r *gin.Engine) {
	log.Println("=== 注册的路由信息 ===")

	routes := r.Routes()
	for _, route := range routes {
		log.Printf("%-8s %s", route.Method, route.Path)
	}

	log.Println("=====================")
}
