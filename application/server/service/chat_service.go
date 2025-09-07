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

func (s *ChatService) SendMessage(senderID int, recipientID int, content string) error {
	// 这里统一时间戳为当前时间
	timeStamp := time.Now()
	message := model.Message{
		SenderID:    senderID,
		RecipientID: recipientID,
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
	err = s.db.Where("sender_id = ? AND recipient_id = ?", senderID, recipientID).First(&chatSession).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		chatSession = model.ChatSession{
			SenderID:     senderID,
			RecipientID:  recipientID,
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
	return sessions, nil
}

func (s *ChatService) GetMessages(userID1 int, userID2 int) ([]model.Message, error) {
	var messages []model.Message
	err := s.db.Where("sender_id = ? AND recipient_id = ? OR sender_id = ? AND recipient_id = ?", userID1, userID2, userID2, userID1).Order("time_stamp").Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("获取消息失败：%v", err)
	}
	return messages, nil
}

func (s *ChatService) ReadMessages(messageIDs []int, userID int) error {
	err := s.db.Model(&model.Message{}).Where("id IN (?) AND recipient_id = ?", messageIDs, userID).Update("has_read", true).Error
	if err != nil {
		return fmt.Errorf("标记消息为已读失败：%v", err)
	}
	return nil
}
