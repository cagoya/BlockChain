<template>
  <MarketNav />
    <a-layout style="height: calc(100vh - 64px); background-color: #f0f2f5;">
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
        
        <!-- 搜索用户区域 -->
        <div class="search-user-container">
          <a-input-group compact>
            <a-input
              v-model:value="searchUserId"
              placeholder="输入用户ID搜索"
              style="width: calc(100% - 80px);"
              @press-enter="searchUser"
            />
            <a-button type="primary" @click="searchUser" :loading="searchingUser">
              搜索
            </a-button>
          </a-input-group>
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
                  <a-avatar 
                    :size="48" 
                    :src="item.otherUserAvatar"
                    style="background-color: #40a9ff; font-size: 24px; color: #fff; line-height: 48px;"
                  >
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
                  <a-avatar 
                    size="large"
                    :src="msg.senderAvatar"
                    style="background-color: #40a9ff; font-size: 20px; color: #fff;"
                  >
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
                  <a-avatar 
                    size="large"
                    :src="user.avatarURL"
                    style="background-color: #1890ff; font-size: 20px; color: #fff;"
                  >
                    我
                  </a-avatar>
                </div>
              </div>
            </div>
          </div>
  
          <!-- 消息输入区域 -->
          <div class="input-area-container" ref="inputContainer">
            <div class="resize-handle" @mousedown="startResize"></div>
            <a-textarea
              v-model:value="newMessageContent"
              placeholder="输入消息..."
              :auto-size="{ minRows: 1, maxRows: 10 }"
              @press-enter="sendMessage"
              class="message-textarea"
            />
            <div class="send-button-container">
              <a-button type="primary" @click="sendMessage" class="send-button">发送</a-button>
            </div>
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
  import { chatApi, accountApi, getImageURL, backendURL } from '../api/index';
  
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
    otherUserAvatar: string;
  }
  
  interface Message {
    id: number;
    senderId: number;
    recipientId: number;
    content: string;
    timeStamp: string;
    hasRead: boolean;
    senderAvatar?: string;
  }

  interface UserInfo {
    username: string;
    avatarURL: string;
    id: number;
  }
  
  // 状态定义
  const user = ref<UserInfo>({
    username: '',
    avatarURL: '',
    id: 0
  });
  const chatSessions = ref<ChatSession[]>([]);
  const selectedSession = ref<ChatSession | null>(null);
  const messages = ref<Message[]>([]);
  const newMessageContent = ref<string>('');
  const ws = ref<WebSocket | null>(null);
  const messageContainer = ref<HTMLDivElement | null>(null);
  const inputContainer = ref<HTMLDivElement | null>(null);
  
  // 加载状态
  const loading = ref<boolean>(true);
  const sessionsLoading = ref<boolean>(false);
  const messagesLoading = ref<boolean>(false);
  const wsConnected = ref<boolean>(false);
  
  // 搜索用户相关状态
  const searchUserId = ref<string>('');
  const searchingUser = ref<boolean>(false);
  
  // 输入框调整高度相关状态
  const isResizing = ref<boolean>(false);
  const startY = ref<number>(0);
  const startHeight = ref<number>(0);
  // 从localStorage中加载用户信息
  const loadUserInfo = () => {
    const userInfoString = localStorage.getItem('userInfo');
    if (userInfoString) {
      try {
        const parsedUserInfo = JSON.parse(userInfoString);
        user.value = {
          username: parsedUserInfo.username || user.value.username,
          avatarURL: getImageURL(parsedUserInfo.avatarURL) || user.value.avatarURL,
          id: parsedUserInfo.id || 0
        };
      } catch (e) {
        console.error('解析用户信息失败', e);
      }
    }
  };
  
  // 用户信息映射，用于显示用户名和头像
  const userNamesMap = ref<Record<number, string>>({});
  const userAvatarsMap = ref<Record<number, string>>({});
  
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

  // 开始调整输入框高度
  const startResize = (e: MouseEvent) => {
    e.preventDefault();
    isResizing.value = true;
    startY.value = e.clientY;
    if (inputContainer.value) {
      startHeight.value = inputContainer.value.offsetHeight;
    }
    
    document.addEventListener('mousemove', handleResize);
    document.addEventListener('mouseup', stopResize);
    document.body.style.cursor = 'ns-resize';
    document.body.style.userSelect = 'none';
  };

  // 处理拖动调整
  const handleResize = (e: MouseEvent) => {
    if (!isResizing.value || !inputContainer.value) return;
    
    const deltaY = startY.value - e.clientY; // 向上拖动为正值
    const newHeight = startHeight.value + deltaY;
    const maxHeight = window.innerHeight * 0.25; // 最大高度为页面高度的四分之一
    const minHeight = 60; // 最小高度
    
    const clampedHeight = Math.max(minHeight, Math.min(newHeight, maxHeight));
    inputContainer.value.style.height = `${clampedHeight}px`;
  };

  // 停止调整
  const stopResize = () => {
    isResizing.value = false;
    document.removeEventListener('mousemove', handleResize);
    document.removeEventListener('mouseup', stopResize);
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
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

  // 获取用户头像（带缓存）
  const getUserAvatar = async (userId: number): Promise<string> => {
    // 如果缓存中已有，直接返回
    if (userAvatarsMap.value[userId]) {
      return userAvatarsMap.value[userId];
    }

    try {
      const response = await accountApi.getAvatar(userId);
      if (response.data.code === 200 && response.data.data) {
        // 使用 getImageURL 构建完整的头像URL
        const avatarURL = getImageURL(response.data.data);
        userAvatarsMap.value[userId] = avatarURL;
        return avatarURL;
      }
    } catch (error) {
      console.error('获取用户头像失败:', error);
    }
    
    // 如果获取失败，返回空字符串，使用默认头像
    return '';
  };

  // 搜索用户并开始聊天
  const searchUser = async () => {
    if (!searchUserId.value.trim()) {
      message.warning('请输入用户ID');
      return;
    }

    const userId = parseInt(searchUserId.value.trim());
    if (isNaN(userId)) {
      message.error('请输入有效的用户ID');
      return;
    }

    if (userId === user.value.id) {
      message.warning('不能与自己聊天');
      return;
    }

    searchingUser.value = true;
    try {
      // 检查用户是否存在
      const userName = await getUserInfo(userId);
      if (!userName) {
        message.error('用户不存在');
        return;
      }

      // 检查是否已经存在聊天会话
      const existingSession = chatSessions.value.find(session => session.otherUserId === userId);
      if (existingSession) {
        // 如果会话已存在，直接选择该会话
        await selectUser(existingSession);
        message.success(`已切换到与 ${userName} 的聊天`);
      } else {
        // 创建新的聊天会话
        const userAvatar = await getUserAvatar(userId);
        const newSession: ChatSession = {
          id: 0, // 新会话暂时没有ID
          senderId: user.value.id,
          recipientId: userId,
          lastMessage: '',
          lastActivity: new Date().toISOString(),
          unreadCount: 0,
          otherUserId: userId,
          otherUserName: userName,
          otherUserAvatar: userAvatar
        };

        // 添加到会话列表顶部
        chatSessions.value.unshift(newSession);
        
        // 选择新会话
        await selectUser(newSession);
        message.success(`开始与 ${userName} 聊天`);
      }

      // 清空搜索框
      searchUserId.value = '';
    } catch (error) {
      message.error('搜索用户失败');
      console.error('Error searching user:', error);
    } finally {
      searchingUser.value = false;
    }
  };

  // 获取聊天会话列表
  const fetchChatSessions = async () => {
    sessionsLoading.value = true;
    try {
      const response = await chatApi.getChatSessions();
      if (response.data.code === 200) {
        const sessions = [];
        if (response.data.data === null) {
          chatSessions.value = [];
          return;
        }
        for (const s of response.data.data) {
          const otherUserId = s.senderId === user.value.id ? s.recipientId : s.senderId;
          const otherUserName = await getUserInfo(otherUserId);
          const otherUserAvatar = await getUserAvatar(otherUserId);
          
          sessions.push({
            ...s,
            otherUserId,
            otherUserName,
            otherUserAvatar,
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
        const messagesData = response.data.data;
        
        // 为每条消息添加发送者头像
        for (const msg of messagesData) {
          if (msg.senderId !== user.value.id) {
            msg.senderAvatar = await getUserAvatar(msg.senderId);
          }
        }
        
        messages.value = messagesData;
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
      const messageWithAvatar = {
        ...msgToSend,
        senderAvatar: user.value.avatarURL
      };
      messages.value.push(messageWithAvatar as Message);
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

        // 为接收到的消息添加发送者头像
        if (receivedMessage.senderId !== user.value.id) {
          getUserAvatar(receivedMessage.senderId).then(avatarURL => {
            receivedMessage.senderAvatar = avatarURL;
          });
        }

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

  .search-user-container {
    padding: 16px;
    background-color: #fff;
    border-bottom: 1px solid #e8e8e8;
  }

  .search-user-container :deep(.ant-input) {
    border-radius: 6px;
    border: 1px solid #d9d9d9;
    transition: all 0.3s ease;
  }

  .search-user-container :deep(.ant-input:hover),
  .search-user-container :deep(.ant-input:focus) {
    border-color: #1890ff;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
  }

  .search-user-container :deep(.ant-btn-primary) {
    background-color: #1890ff;
    border-color: #1890ff;
    border-radius: 6px;
    transition: all 0.3s ease;
  }

  .search-user-container :deep(.ant-btn-primary:hover) {
    background-color: #40a9ff;
    border-color: #40a9ff;
    transform: translateY(-1px);
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
    flex-direction: row;
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
    position: relative;
    min-height: 60px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .resize-handle {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background-color: #d9d9d9;
    cursor: ns-resize;
    transition: background-color 0.2s;
  }

  .resize-handle:hover {
    background-color: #1890ff;
  }

  .message-textarea {
    flex: 1;
    resize: none;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    padding: 8px 12px;
    font-size: 14px;
    line-height: 1.5;
    transition: border-color 0.3s;
  }

  .message-textarea:focus {
    border-color: #1890ff;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
    outline: none;
  }

  .send-button-container {
    display: flex;
    justify-content: flex-end;
  }

  .send-button {
    background-color: #1890ff;
    border-color: #1890ff;
    border-radius: 6px;
    transition: all 0.3s ease;
  }

  .send-button:hover {
    background-color: #40a9ff;
    border-color: #40a9ff;
    transform: translateY(-1px);
  }
  
  .empty-chat-window {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
  }
  </style>
  