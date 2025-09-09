<template>
    <a-layout style="min-height: 100vh; background-color: #f0f2f5;">
      <!-- 左侧聊天会话列表 -->
      <a-layout-sider width="300" style="background-color: #fff; border-right: 1px solid #e8e8e8;">
        <div class="logo-container">
          <a-typography-title :level="4" style="color: #1890ff;">Web Chat</a-typography-title>
          <div class="connection-status">
            <a-badge 
              :status="wsConnected ? 'success' : 'error'" 
              :text="wsConnected ? '已连接' : '未连接'"
            />
          </div>
        </div>
        <div v-if="sessionsLoading" class="loading-container">
          <a-spin size="large" />
          <p>加载聊天会话中...</p>
        </div>
        <a-list v-else item-layout="horizontal" :data-source="chatSessions" :split="false">
          <template #renderItem="{ item }">
            <a-list-item
              :class="['chat-session-item', { 'active': selectedSession && selectedSession.otherUserId === item.otherUserId }]"
              @click="selectUser(item)"
            >
              <a-list-item-meta>
                <template #avatar>
                  <a-avatar :size="48" style="background-color: #40a9ff; font-size: 24px; color: #fff; line-height: 48px;">
                    {{ item.otherUserName.charAt(0).toUpperCase() }}
                  </a-avatar>
                </template>
                <template #title>
                  <a-row align="middle">
                    <a-col flex="1">
                      <span class="user-name">{{ item.otherUserName }}</span>
                    </a-col>
                    <a-col flex="none">
                      <a-badge :count="item.unreadCount" :overflow-count="99" :offset="[0, 10]" />
                    </a-col>
                  </a-row>
                </template>
                <template #description>
                  <span class="last-message">{{ item.lastMessage }}</span>
                </template>
              </a-list-item-meta>
            </a-list-item>
          </template>
        </a-list>
      </a-layout-sider>
  
      <!-- 右侧聊天窗口 -->
      <a-layout-content style="position: relative; overflow: hidden;">
        <div v-if="selectedSession" class="chat-window-container">
          <!-- 聊天窗口头部 -->
          <a-layout-header class="chat-header">
            <a-typography-title :level="5" style="color: #333; margin: 0;">{{ selectedSession.otherUserName }}</a-typography-title>
          </a-layout-header>
  
          <!-- 聊天消息列表 -->
          <div class="message-list-container" ref="messageContainer">
            <div v-if="messagesLoading" class="loading-container">
              <a-spin size="large" />
              <p>加载消息中...</p>
            </div>
            <div v-else class="message-list-inner">
              <div v-for="msg in messages" :key="msg.id" :class="['message-item', { 'sent': msg.senderId === user.id, 'received': msg.senderId !== user.id }]">
                <div class="avatar-wrapper" v-if="msg.senderId !== user.id">
                  <a-avatar style="background-color: #40a9ff; font-size: 20px; color: #fff;">
                    {{ selectedSession.otherUserName.charAt(0).toUpperCase() }}
                  </a-avatar>
                </div>
                <div class="message-bubble-wrapper">
                  <div class="message-content-wrapper">
                    <div class="message-content">
                      {{ msg.content }}
                    </div>
                  </div>
                  <div class="time-stamp">{{ formatTime(msg.timeStamp) }}</div>
                </div>
                <div class="avatar-wrapper" v-if="msg.senderId === user.id">
                  <a-avatar style="background-color: #1890ff; font-size: 20px; color: #fff;">
                    我
                  </a-avatar>
                </div>
              </div>
            </div>
          </div>
  
          <!-- 消息输入区域 -->
          <div class="input-area-container">
            <a-input-group compact style="display: flex;">
              <a-input
                v-model:value="newMessageContent"
                placeholder="输入消息..."
                @press-enter="sendMessage"
                style="flex: 1;"
              />
              <a-button type="primary" @click="sendMessage">发送</a-button>
            </a-input-group>
          </div>
        </div>
        <div v-else class="empty-chat-window">
          <a-empty description="选择一个聊天会话开始聊天" />
        </div>
      </a-layout-content>
    </a-layout>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, watch, nextTick } from 'vue';
  import { message } from 'ant-design-vue';
  import { chatApi, accountApi } from '../api/index';
  
  // 接口定义
  interface ChatSession {
    // 这部分和后端相同
    id: number;
    senderId: number;
    recipientId: number;
    lastMessage: string;
    lastActivity: string;
    // 这部分是前端增加的
    unreadCount: number;
    otherUserId: number;
    otherUserName: string;
  }
  
  interface Message {
    id: number;
    senderId: number;
    recipientId: number;
    content: string;
    timeStamp: string;
    hasRead: boolean;
  }

  interface UserInfo {
    username: string;
    avatarURL: string;
    id: number;
  }
  
  // 状态定义
  const user = ref<UserInfo>({
    username: '游客',
    avatarURL: 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
    id: 0
  });
  const chatSessions = ref<ChatSession[]>([]);
  const selectedSession = ref<ChatSession | null>(null);
  const messages = ref<Message[]>([]);
  const newMessageContent = ref<string>('');
  const ws = ref<WebSocket | null>(null);
  const messageContainer = ref<HTMLDivElement | null>(null);
  
  // 加载状态
  const loading = ref<boolean>(true);
  const sessionsLoading = ref<boolean>(false);
  const messagesLoading = ref<boolean>(false);
  const wsConnected = ref<boolean>(false);
  // 从localStorage中加载用户信息
  const loadUserInfo = () => {
    const userInfoString = localStorage.getItem('userInfo');
    if (userInfoString) {
      try {
        const parsedUserInfo = JSON.parse(userInfoString);
        user.value = {
          username: parsedUserInfo.username || user.value.username,
          avatarURL: parsedUserInfo.avatarURL || user.value.avatarURL,
          id: parsedUserInfo.id || 0
        };
      } catch (e) {
        console.error('解析用户信息失败', e);
      }
    }
  };
  
  // 用户信息映射，用于显示用户名
  const userNamesMap = ref<Record<number, string>>({});
  
  // 格式化时间
  const formatTime = (time: string): string => {
    const date = new Date(time);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };
  
  // 滚动到底部
  const scrollToBottom = () => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
    }
  };

  // 获取用户信息（带缓存）
  const getUserInfo = async (userId: number): Promise<string> => {
    // 如果缓存中已有，直接返回
    if (userNamesMap.value[userId]) {
      return userNamesMap.value[userId];
    }

    const response = await accountApi.getUserNameById(userId);
    if (response.data.code === 200) {
      userNamesMap.value[userId] = response.data.data;
      return response.data.data;
    }
    return '';
  };

  // 获取聊天会话列表
  const fetchChatSessions = async () => {
    sessionsLoading.value = true;
    try {
      const response = await chatApi.getChatSessions();
      if (response.data.code === 200) {
        const sessions = [];
        
        for (const s of response.data.data) {
          const otherUserId = s.senderId === user.value.id ? s.recipientId : s.senderId;
          const otherUserName = await getUserInfo(otherUserId);
          
          sessions.push({
            ...s,
            otherUserId,
            otherUserName,
            unreadCount: 0
          });
        }
        
        // 按最后活动时间排序
        sessions.sort((a: ChatSession, b: ChatSession) => new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime());
        chatSessions.value = sessions;
        
        // 为每个会话计算未读消息数
        for (const session of chatSessions.value) {
          await fetchUnreadCount(session);
        }
      } else {
        message.error('获取聊天会话失败');
      }
    } catch (error) {
      message.error('获取聊天会话失败');
      console.error('Error fetching chat sessions:', error);
    } finally {
      sessionsLoading.value = false;
    }
  };

  // 获取指定会话的未读消息数
  const fetchUnreadCount = async (session: ChatSession) => {
    try {
      const response = await chatApi.getUnreadMessageCount(session.otherUserId);
      if (response.data.code === 200) {
        session.unreadCount = response.data.data;
      }
    } catch (error) {
      console.error('Error fetching unread count:', error);
    }
  };
  
  // 获取消息列表
  const fetchMessages = async (otherId: number) => {
    messagesLoading.value = true;
    try {
      const response = await chatApi.getMessages(otherId);
      if (response.data.code === 200) {
        messages.value = response.data.data;
        nextTick(() => {
          scrollToBottom();
        });
        // 标记为已读
        if (response.data.data.length > 0) {
          await readMessages(otherId);
        }
      } else {
        message.error('获取消息失败');
      }
    } catch (error) {
      message.error('获取消息失败');
      console.error('Error fetching messages:', error);
    } finally {
      messagesLoading.value = false;
    }
  };
  
  // 标记消息为已读
  const readMessages = async (otherId: number) => {
    try {
      await chatApi.readMessages(otherId);
      // 更新本地未读消息数
      if (selectedSession.value) {
        selectedSession.value.unreadCount = 0;
      }
    } catch (error) {
      console.error('Error marking messages as read:', error);
    }
  };
  
  // 选择聊天用户
  const selectUser = async (session: ChatSession) => {
    selectedSession.value = session;
    await fetchMessages(session.otherUserId);
    // 消息获取后会自动标记为已读，这里重置未读数为0
    session.unreadCount = 0;
  };
  
  // 发送消息
  const sendMessage = async () => {
    if (!newMessageContent.value.trim() || !selectedSession.value) {
      message.warning('请输入消息内容');
      return;
    }

    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      message.error('WebSocket连接未建立，无法发送消息');
      return;
    }

    const msgToSend = {
      senderId: user.value.id,
      recipientId: selectedSession.value.otherUserId,
      content: newMessageContent.value.trim(),
      timeStamp: new Date().toISOString(),
      hasRead: false,
    };

    try {
      ws.value.send(JSON.stringify(msgToSend));
      
      // 立即在本地界面显示
      messages.value.push(msgToSend as Message);
      newMessageContent.value = '';
      nextTick(() => {
        scrollToBottom();
      });
      
      // 更新会话列表的最后一条消息
      selectedSession.value.lastMessage = msgToSend.content;
      selectedSession.value.lastActivity = msgToSend.timeStamp;
      
      // 重新排序会话列表
      chatSessions.value.sort((a, b) => new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime());
      
    } catch (error) {
      message.error('发送消息失败，请重试');
      console.error('Error sending message:', error);
    }
  };
  
import { backendURL } from '../api/index';

// 连接WebSocket
const connectWebSocket = () => {
  // 构建WebSocket URL，包含认证token
  const wsUrl = `${backendURL.replace('http', 'ws')}/chat/ws?userId=${user.value.id}`;
  ws.value = new WebSocket(wsUrl);

    ws.value.onopen = () => {
      console.log('WebSocket连接已建立');
      wsConnected.value = true;
    };

    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        // 检查是否是错误消息
        if (data.error) {
          message.error(data.error);
          return;
        }

        const receivedMessage: Message = data;

        // 检查消息是否与当前选中的聊天对象相关
        if (selectedSession.value && (receivedMessage.senderId === selectedSession.value.otherUserId || receivedMessage.recipientId === selectedSession.value.otherUserId)) {
          messages.value.push(receivedMessage);
          nextTick(() => {
            scrollToBottom();
          });
          // 自动标记为已读
          if (receivedMessage.recipientId === user.value.id) {
            readMessages(receivedMessage.senderId);
          }
        } else if (receivedMessage.recipientId === user.value.id) {
          // 如果消息是发给我的，但不是当前聊天会话，则增加未读消息数
          const session = chatSessions.value.find(s => s.otherUserId === receivedMessage.senderId);
          if (session) {
            session.unreadCount++;
          }
        }
        
        // 更新会话列表的最后一条消息和时间
        const sessionToUpdate = chatSessions.value.find(s =>
          (s.otherUserId === receivedMessage.senderId || s.otherUserId === receivedMessage.recipientId)
        );
        if (sessionToUpdate) {
          sessionToUpdate.lastMessage = receivedMessage.content;
          sessionToUpdate.lastActivity = receivedMessage.timeStamp;
          chatSessions.value.sort((a, b) => new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime());
        }
      } catch (error) {
        console.error('解析WebSocket消息失败:', error);
      }
    };

    ws.value.onerror = (error) => {
      console.error('WebSocket连接错误:', error);
      message.error('WebSocket连接失败');
      wsConnected.value = false;
    };

    ws.value.onclose = (event) => {
      console.log('WebSocket连接已关闭:', event.code, event.reason);
      wsConnected.value = false;
      // 如果不是正常关闭，尝试重连
      if (event.code !== 1000) {
        console.log('WebSocket连接异常关闭，3秒后尝试重连...');
        setTimeout(() => {
          if (user.value.id) {
            connectWebSocket();
          }
        }, 3000);
      }
    };
  };
  
  onMounted(async () => {
    loading.value = true;
    try {
      loadUserInfo();
      await fetchChatSessions();
      connectWebSocket();
    } catch (error) {
      message.error('初始化聊天功能失败');
      console.error('Error initializing chat:', error);
    } finally {
      loading.value = false;
    }
  });
  
  // 监听 selectedSession 变化，如果变化了，滚动到底部
  watch(selectedSession, () => {
    nextTick(() => {
      scrollToBottom();
    });
  });
  
  </script>
  
  <style scoped>
  .logo-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 64px;
    padding: 0 16px;
    background-color: #fff;
    border-bottom: 1px solid #e8e8e8;
  }

  .connection-status {
    display: flex;
    align-items: center;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 40px 20px;
    color: #666;
  }

  .loading-container p {
    margin-top: 16px;
    margin-bottom: 0;
  }
  
  .chat-session-item {
    padding: 12px 16px;
    cursor: pointer;
    transition: all 0.3s;
    border-bottom: 1px solid #f0f0f0;
  }
  
  .chat-session-item:hover {
    background-color: #f5f5f5;
  }
  
  .chat-session-item.active {
    background-color: #e6f7ff;
    border-left: 4px solid #1890ff;
    padding-left: 12px;
  }
  
  .user-name {
    font-weight: 500;
    color: #333;
  }
  
  .last-message {
    font-size: 12px;
    color: #8c8c8c;
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 200px;
  }
  
  .chat-window-container {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  
  .chat-header {
    background-color: #fff;
    padding: 0 24px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #e8e8e8;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }
  
  .message-list-container {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    background-color: #f0f2f5;
    display: flex;
    flex-direction: column;
  }
  
  .message-list-inner {
    display: flex;
    flex-direction: column;
    width: 100%;
  }
  
  .message-item {
    display: flex;
    align-items: flex-end;
    margin-bottom: 15px;
    max-width: 70%;
  }
  
  .message-item.sent {
    align-self: flex-end;
    flex-direction: row-reverse;
  }
  
  .message-bubble-wrapper {
    display: flex;
    flex-direction: column;
  }
  
  .message-content-wrapper {
    display: flex;
    max-width: 100%;
  }
  
  .message-content {
    padding: 12px 16px;
    border-radius: 20px;
    line-height: 1.5;
    font-size: 14px;
    word-break: break-word;
  }
  
  .message-item.received .message-content {
    background-color: #fff;
    color: #333;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    border-bottom-left-radius: 2px;
  }
  
  .message-item.sent .message-content {
    background-color: #1890ff;
    color: #fff;
    border-bottom-right-radius: 2px;
  }
  
  .avatar-wrapper {
    margin: 0 10px;
  }
  
  .time-stamp {
    font-size: 10px;
    color: #999;
    margin-top: 4px;
    align-self: flex-start;
    margin-left: 10px;
  }
  
  .message-item.sent .time-stamp {
    align-self: flex-end;
    margin-right: 10px;
  }
  
  .input-area-container {
    padding: 16px 24px;
    background-color: #fff;
    border-top: 1px solid #e8e8e8;
  }
  
  .empty-chat-window {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
  }
  </style>
  