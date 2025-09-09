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
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 添加静态文件服务
	r.Static("/public", "./public")

	apiGroup := r.Group("/api")

	// 注册路由
	accountHandler := api.NewAccountHandler()
	walletHandler := api.NewWalletHandler()
	assetHandler := api.NewAssetHandler()
	chatHandler, err := api.NewChatHandler()
	if err != nil {
		log.Fatalf("创建聊天处理程序失败：%v", err)
	}

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
	authAccount := apiGroup.Group("/account").Use(jwtMiddleware.Auth())
	{
		// 获取用户信息
		authAccount.GET("/profile", accountHandler.GetProfile)
		// 更新用户信息
		authAccount.PUT("/profile", accountHandler.UpdateProfile)
		// 获取头像
		authAccount.GET("/avatar", accountHandler.GetAvatar)
		// 更新头像
		authAccount.PUT("/avatar", accountHandler.UpdateAvatar)
		// 更新组织接口
		authAccount.PUT("/org", accountHandler.UpdateOrg)
	}

	// 钱包相关接口
	wallet := apiGroup.Group("/wallet").Use(jwtMiddleware.Auth())
	{
		wallet.POST("/create", walletHandler.CreateAccount)
		wallet.GET("/balance", walletHandler.GetBalance)
		wallet.POST("/transfer", walletHandler.Transfer)
		wallet.POST("/mintToken", walletHandler.MintToken)
		wallet.GET("/transferBySenderID", walletHandler.GetTransferBySenderID)
		wallet.GET("/transferByRecipientID", walletHandler.GetTransferByRecipientID)
		wallet.POST("/withHoldAccount", walletHandler.WithHoldAccount)
		wallet.GET("/getWithHoldingByAccountID", walletHandler.GetWithHoldingByAccountID)
		wallet.GET("/getWithHoldingByListingID", walletHandler.GetWithHoldingByListingID)
		wallet.POST("/clearWithHolding", walletHandler.ClearWithHolding)
	}

	// 资产相关接口
	asset := apiGroup.Group("/asset").Use(jwtMiddleware.Auth())
	{
		asset.POST("/create", assetHandler.CreateAsset)
		asset.GET("/getAssetByID", assetHandler.GetAssetByID)
		asset.GET("/getAssetByAuthorID", assetHandler.GetAssetByAuthorID)
		asset.GET("/getAssetByOwnerID", assetHandler.GetAssetByOwnerID)
		asset.POST("/transfer", assetHandler.TransferAsset)
	}

	// 聊天相关接口
	chat := apiGroup.Group("/chat").Use(jwtMiddleware.Auth())
	{
		chat.GET("/ws", chatHandler.SendMessage)
		chat.GET("/getChatSession", chatHandler.GetChatSession)
		chat.GET("/getMessages", chatHandler.GetMessages)
		chat.POST("/readMessages", chatHandler.ReadMessages)
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
