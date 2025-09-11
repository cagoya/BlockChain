package service

import (
	"errors"
	"fmt"
	"time"

	"application/model"

	"gorm.io/gorm"
)

type ChatService struct {
	db *gorm.DB
}

func NewChatService() (*ChatService, error) {
	db := model.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return &ChatService{db: db}, nil
}

func (s *ChatService) SendMessage(senderId int, recipientId int, content string) error {
	// 这里统一时间戳为当前时间
	timeStamp := time.Now()
	message := model.Message{
		SenderID:    senderId,
		RecipientID: recipientId,
		Content:     content,
		TimeStamp:   timeStamp,
	}
	err := s.db.Create(&message).Error
	if err != nil {
		return fmt.Errorf("发送消息失败：%v", err)
	}
	// 查询是否存在聊天会话,如果不存在，则创建新的会话
	// 存在则更新会话时间和内容
	var chatSession model.ChatSession
	err = s.db.Where("sender_id = ? AND recipient_id = ?", senderId, recipientId).First(&chatSession).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		chatSession = model.ChatSession{
			SenderID:     senderId,
			RecipientID:  recipientId,
			LastMessage:  content,
			LastActivity: timeStamp,
		}
		err = s.db.Create(&chatSession).Error
		if err != nil {
			return fmt.Errorf("创建聊天会话失败：%v", err)
		}
	} else {
		chatSession.LastActivity = timeStamp
		chatSession.LastMessage = content
		err = s.db.Save(&chatSession).Error
		if err != nil {
			return fmt.Errorf("更新聊天会话失败：%v", err)
		}
	}
	return nil
}

func (s *ChatService) GetChatSession(userID int) ([]model.ChatSession, error) {
	var sessions []model.ChatSession
	err := s.db.Where("sender_id = ? OR recipient_id = ?", userID, userID).Order("last_activity DESC").Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("获取聊天会话失败：%v", err)
	}
	// 使用一个map来进行去重，key为对方ID
	recentSessions := make(map[int]model.ChatSession)
	for _, session := range sessions {
		// 确定对话的另一方是谁
		var otherID int
		if session.SenderID == userID {
			otherID = session.RecipientID
		} else {
			otherID = session.SenderID
		}

		// 如果map中还没有这个对话的记录，就将当前记录存入
		// 因为我们是按时间倒序遍历，所以第一条遇到的记录就是最新的
		if _, ok := recentSessions[otherID]; !ok {
			recentSessions[otherID] = session
		}
	}

	// 5. 将map中的值转换成切片并返回
	var result []model.ChatSession
	for _, session := range recentSessions {
		result = append(result, session)
	}

	return result, nil
}

func (s *ChatService) GetMessages(userID1 int, userID2 int) ([]model.Message, error) {
	var messages []model.Message
	err := s.db.Where("sender_id = ? AND recipient_id = ? OR sender_id = ? AND recipient_id = ?", userID1, userID2, userID2, userID1).Order("time_stamp").Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("获取消息失败：%v", err)
	}
	return messages, nil
}

func (s *ChatService) ReadMessages(myID int, otherID int) error {
	err := s.db.Model(&model.Message{}).Where("recipient_id = ? AND sender_id = ?", myID, otherID).Update("has_read", true).Error
	if err != nil {
		return fmt.Errorf("标记消息为已读失败：%v", err)
	}
	return nil
}

func (s *ChatService) GetUnreadMessageCount(myID int, otherID int) (int, error) {
	var count int64
	err := s.db.Model(&model.Message{}).Where("recipient_id = ? AND sender_id = ? AND has_read = ?", myID, otherID, false).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取未读消息数量失败：%v", err)
	}
	return int(count), nil
}
