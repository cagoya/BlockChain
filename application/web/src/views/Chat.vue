<template>
    <a-layout style="min-height: 100vh; background-color: #f0f2f5;">
      <!-- 左侧聊天会话列表 -->
      <a-layout-sider width="300" style="background-color: #fff; border-right: 1px solid #e8e8e8;">
        <div class="logo-container">
          <a-typography-title :level="4" style="color: #1890ff;">Web Chat</a-typography-title>
        </div>
        <a-list item-layout="horizontal" :data-source="chatSessions" :split="false">
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
                      <!-- TODO: 后端接口未提供未读消息数，此为模拟数据，请根据实际接口调整 -->
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
            <div class="message-list-inner">
              <div v-for="msg in messages" :key="msg.id" :class="['message-item', { 'sent': msg.senderId === myId, 'received': msg.senderId !== myId }]">
                <div class="avatar-wrapper" v-if="msg.senderId !== myId">
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
                <div class="avatar-wrapper" v-if="msg.senderId === myId">
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
  
  // 定义后端 API 的基地址和 WebSocket 的基地址
  const BASE_API_URL = 'http://localhost:8888/api';
  const BASE_WS_URL = 'ws://localhost:8888/api';
  
  // 接口定义
  interface ChatSession {
    id: number;
    senderId: number;
    recipientId: number;
    lastMessage: string;
    lastActivity: string;
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
  
  // 状态
  const myId = ref<number>(1); // 当前用户ID，应通过登录接口获取
  const chatSessions = ref<ChatSession[]>([]);
  const selectedSession = ref<ChatSession | null>(null);
  const messages = ref<Message[]>([]);
  const newMessageContent = ref<string>('');
  const ws = ref<WebSocket | null>(null);
  const messageContainer = ref<HTMLDivElement | null>(null);
  
  // TODO: 假设存在一个用户信息存储，实际应通过后端接口获取
  // 为了确保页面可用，这里使用一个模拟映射，请替换为实际接口
  const userNamesMap = ref<Record<number, string>>({
    1: '我',
    2: '张三',
    3: '李四',
    4: '王五',
  });
  
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
  
  // 获取聊天会话列表
  const fetchChatSessions = async () => {
    try {
      const response = await fetch(`${BASE_API_URL}/chat/getChatSession`);
      const data = await response.json();
      if (data.success) {
        const sessions = data.data.map((s: any) => {
          const otherUserId = s.senderId === myId.value ? s.recipientId : s.senderId;
          const otherUserName = userNamesMap.value[otherUserId] || '未知用户';
          return {
            ...s,
            otherUserId,
            otherUserName,
            // TODO: 未读消息数需要后端在getChatSession接口中返回，目前为模拟
            unreadCount: Math.floor(Math.random() * 5),
          };
        }).sort((a: ChatSession, b: ChatSession) => new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime());
        chatSessions.value = sessions;
        message.success('聊天会话已更新');
      } else {
        message.error('获取聊天会话失败');
      }
    } catch (error) {
      message.error('请求失败');
      console.error('Error fetching chat sessions:', error);
    }
  };
  
  // 获取消息列表
  const fetchMessages = async (otherId: number) => {
    try {
      const response = await fetch(`${BASE_API_URL}/chat/getMessages?otherID=${otherId}`);
      const data = await response.json();
      if (data.success) {
        messages.value = data.data;
        nextTick(() => {
          scrollToBottom();
        });
        // 标记为已读
        const unreadMessageIds = data.data.filter((msg: Message) => !msg.hasRead && msg.recipientId === myId.value).map((msg: Message) => msg.id);
        if (unreadMessageIds.length > 0) {
          await readMessages(unreadMessageIds);
        }
      } else {
        message.error('获取消息失败');
      }
    } catch (error) {
      message.error('请求失败');
      console.error('Error fetching messages:', error);
    }
  };
  
  // 标记消息为已读
  const readMessages = async (messageIds: number[]) => {
    try {
      await fetch(`${BASE_API_URL}/chat/readMessages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(messageIds),
      });
      // 更新本地未读消息数
      if (selectedSession.value) {
        selectedSession.value.unreadCount = 0;
      }
    } catch (error) {
      console.error('Error marking messages as read:', error);
    }
  };
  
  // 选择聊天用户
  const selectUser = (session: ChatSession) => {
    selectedSession.value = session;
    fetchMessages(session.otherUserId);
    if (session.unreadCount > 0) {
      session.unreadCount = 0;
    }
  };
  
  // 发送消息
  const sendMessage = () => {
    if (!newMessageContent.value.trim() || !selectedSession.value) return;
  
    const msgToSend = {
      senderId: myId.value,
      recipientId: selectedSession.value.otherUserId,
      content: newMessageContent.value,
      timeStamp: new Date().toISOString(),
      hasRead: false,
    };
  
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(msgToSend));
    } else {
      // 如果WebSocket未连接，直接通过HTTP API发送
      // 实际应用中需要后端提供相应的HTTP发送接口
      console.warn('WebSocket is not connected, fallback to HTTP...');
    }
  
    // 立即在本地界面显示
    messages.value.push(msgToSend as Message);
    newMessageContent.value = '';
    nextTick(() => {
      scrollToBottom();
    });
  };
  
  // 连接WebSocket
  const connectWebSocket = () => {
    // 确保URL协议正确，通常为 ws:// 或 wss://
    const wsUrl = `${BASE_WS_URL}/chat/ws`;
    ws.value = new WebSocket(wsUrl);
  
    ws.value.onmessage = (event) => {
      const receivedMessage: Message = JSON.parse(event.data);
      // 检查消息是否与当前选中的聊天对象相关
      if (selectedSession.value && (receivedMessage.senderId === selectedSession.value.otherUserId || receivedMessage.recipientId === selectedSession.value.otherUserId)) {
        messages.value.push(receivedMessage);
        nextTick(() => {
          scrollToBottom();
        });
        // 自动标记为已读
        if (receivedMessage.recipientId === myId.value) {
          readMessages([receivedMessage.id]);
        }
      } else if (receivedMessage.recipientId === myId.value) {
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
    };
  
    ws.value.onclose = () => {
      console.log('WebSocket connection closed. Attempting to reconnect...');
      setTimeout(connectWebSocket, 3000);
    };
  };
  
  onMounted(() => {
    fetchChatSessions();
    connectWebSocket();
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
    justify-content: center;
    align-items: center;
    height: 64px;
    background-color: #fff;
    border-bottom: 1px solid #e8e8e8;
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
  