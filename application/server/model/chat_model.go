package model

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 消息
type Message struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SenderID    int       `json:"senderId" gorm:"not null"`
	RecipientID int       `json:"recipientId" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	TimeStamp   time.Time `json:"timeStamp" gorm:"autoCreateTime"`
	HasRead     bool      `json:"hasRead" gorm:"default:false"`
}

type ChatSession struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SenderID     int       `json:"senderId" gorm:"not null"`
	RecipientID  int       `json:"recipientId" gorm:"not null"`
	LastMessage  string    `json:"lastMessage" gorm:"type:text;not null"`
	LastActivity time.Time `json:"lastActivity" gorm:"autoCreateTime"`
}

func (Message) TableName() string {
	return "messages"
}

func (ChatSession) TableName() string {
	return "chatSessions"
}

// 定义全局的 WebSocket 连接管理器
type ConnManager struct {
	// 使用 sync.Map 来保证并发安全
	connections sync.Map
}

// 单例模式
var connManager = &ConnManager{}

func GetConnManager() *ConnManager {
	return connManager
}

// 向管理器添加连接
func (m *ConnManager) AddConn(userID int, conn *websocket.Conn) {
	m.connections.Store(userID, conn)
}

// 从管理器获取连接
func (m *ConnManager) GetConn(userID int) (*websocket.Conn, bool) {
	conn, ok := m.connections.Load(userID)
	if !ok {
		return nil, false
	}
	return conn.(*websocket.Conn), true
}

// 从管理器移除连接
func (m *ConnManager) RemoveConn(userID int) {
	m.connections.Delete(userID)
}

// 检查连接是否有效
func (m *ConnManager) IsConnValid(userID int) bool {
	conn, ok := m.GetConn(userID)
	if !ok {
		return false
	}

	// 尝试发送 ping 消息来检测连接是否有效
	err := conn.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		// 连接无效，从管理器中移除
		m.RemoveConn(userID)
		return false
	}
	return true
}

// 获取所有在线用户ID
func (m *ConnManager) GetOnlineUsers() []int {
	var users []int
	m.connections.Range(func(key, value interface{}) bool {
		if userID, ok := key.(int); ok {
			users = append(users, userID)
		}
		return true
	})
	return users
}

// WebSocket 升级器
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许Vue默认端口和后端运行端口
		return r.Host == "localhost:5173" || r.Host == "localhost:8888"
	},
}

type SendMessageRequest struct {
	RecipientID int    `json:"recipient_id"`
	Content     string `json:"content"`
}
