package api

import (
	"application/model"
	"application/service"
	"application/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler() (*ChatHandler, error) {
	chatService, err := service.NewChatService()
	if err != nil {
		return nil, err
	}
	return &ChatHandler{
		chatService: chatService,
	}, nil
}

// 发送消息，使用 websocket 发送消息
// 如果对方在线，则直接推送
// 如果对方不在线，则将消息先存入数据库
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	// 升级HTTP连接到WebSocket
	conn, err := model.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	defer conn.Close()

	// 将新连接添加到连接管理器
	connManager := model.GetConnManager()
	connManager.AddConn(userID.(int), conn)

	// 监听连接断开事件，并从管理器中移除
	defer func() {
		connManager.RemoveConn(userID.(int))
	}()

	// 设置读取超时和关闭处理
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 启动心跳检测
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	// 循环读取客户端发送的消息
	for {
		var msg model.Message
		// 读取JSON消息
		err := conn.ReadJSON(&msg)
		if err != nil {
			// 检查是否是连接关闭错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 连接异常关闭
				fmt.Printf("WebSocket连接异常关闭: %v\n", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				// 正常关闭
				fmt.Printf("WebSocket连接正常关闭\n")
			} else {
				// 其他错误（如超时、格式错误等）
				fmt.Printf("WebSocket读取错误: %v\n", err)
			}
			break
		}

		// 重置读取超时
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// 根据RecipientID获取目标连接
		receiverConn, ok := connManager.GetConn(msg.RecipientID)
		if ok {
			// 将消息转发给接收者
			err = receiverConn.WriteJSON(msg)
			if err != nil {
				fmt.Printf("转发消息失败: %v\n", err)
				// 如果转发失败，可能是接收者连接已断开，从管理器中移除
				connManager.RemoveConn(msg.RecipientID)
			}
		}

		// 写入数据库
		error := h.chatService.SendMessage(userID.(int), msg.RecipientID, msg.Content)
		if error != nil {
			fmt.Printf("保存消息到数据库失败: %v\n", error)
			// 发送错误消息给客户端
			conn.WriteJSON(gin.H{"error": "发送消息失败，请重试"})
		}
	}
}

// 获取聊天会话，即获取与当前用户有聊天记录的用户列表，按照最后一条消息的创建时间倒序排序
func (h *ChatHandler) GetChatSession(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ServerError(c, "用户信息获取失败")
		return
	}
	sessions, err := h.chatService.GetChatSession(userID.(int))
	if err != nil {
		utils.ServerError(c, "获取聊天会话失败")
		return
	}
	utils.Success(c, sessions)
}

// 获取两位用户之间的所有聊天记录
func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID1, err := strconv.Atoi(c.Query("userID1"))
	if err != nil {
		utils.BadRequest(c, "用户ID1不能为空")
		return
	}
	userID2, err := strconv.Atoi(c.Query("userID2"))
	if err != nil {
		utils.BadRequest(c, "用户ID2不能为空")
		return
	}
	messages, err := h.chatService.GetMessages(userID1, userID2)
	if err != nil {
		utils.ServerError(c, "获取消息失败")
		return
	}
	utils.Success(c, messages)
}
