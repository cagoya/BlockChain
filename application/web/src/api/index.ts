import axios from 'axios';

// 本地开发可以选择 localhost:8888/api
// 如果需要多机调试，将 localhost 改为你的 ip 地址
// 这样在局域网内可以访问
//export const backendURL = 'http://10.162.199.212:8888/api';
export const backendURL = 'http://localhost:8888/api';

// 获取图片的完整 URL
export const getImageURL = (imageName: string) => {
  return `${backendURL.replace('/api', '')}/public/images/${imageName}`;
};

// 创建axios实例
const instance = axios.create({
  baseURL: backendURL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器添加认证信息
instance.interceptors.request.use(
  (config) => {
    // 检查是否为不需要认证的接口
    const isAuthRequest = config.url?.includes('/account/login') || 
                         config.url?.includes('/account/register');
    
    // 如果不是认证相关接口，则添加token
    if (!isAuthRequest) {
      const token = localStorage.getItem('userToken');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器处理错误
instance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error.response?.data || '请求失败');
  }
);

// 账户相关API
const accountApi = {
  /**
   * 用户登录
   * @param username 用户名
   * @param password 密码
   */
  login: (username: string, password: string) => {
    return instance.post('/account/login', {
      Username: username,
      Password: password
    });
  },

  /**
   * 用户注册
   * @param username 用户名
   * @param email 邮箱
   * @param password 密码
   * @param org 组织ID
   */
  register: (username: string, email: string, password: string, org: number) => {
    return instance.post('/account/register', {
      Username: username,
      Email: email,
      Password: password,
      Org: org
    });
  },

  /**
   * 用户登出
   */
  logout: () => {
    return instance.post('/account/logout');
  },

  /**
   * 更新用户个人资料
   * @param profileData 个人资料数据
   */
  updateProfile: (profileData: any) => {
    return instance.put('/account/profile', profileData, {
      headers: {
        'Content-Type': 'application/json'
      }
    });
  },

  /**
   * 更新用户头像
   * @param formData 头像文件数据
   */
  updateAvatar: (formData: FormData) => {
    return instance.put('/account/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },

  /**
   * 更新用户组织
   * @param orgData 组织数据
   */
  updateOrg: (orgData: any) => {
    return instance.put('/account/org', orgData, {
      headers: {
        'Content-Type': 'application/json'
      }
    });
  },

  /**
   * 获取用户个人资料
   */
  getProfile: () => {
    return instance.get('/account/profile');
  },

  /**
   * 获取用户名
   * @param userId 用户ID
   */
  getUserNameById: (userId: number) => {
    return instance.get(`/account/userName?id=${userId}`);
  },

  /**
   * 获取头像
   * @param userId 用户ID
   */
  getAvatar: (userId: number) => {
    return instance.get(`/account/avatar?id=${userId}`);
  }

};

// 资产相关API
const assetApi = {
  /**
   * 创建NFT资产
   * @param assetData 资产数据
   */
  create: (assetData: FormData) => {
    return instance.post('/asset/create', assetData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },

  /**
   * 根据Id获取资产
   * @param id 资产ID
   */
  getById: (id: string) => {
    return instance.get(`/asset/getAssetByID?id=${id}`);
  },

  /**
   * 根据作者ID获取资产
   * @param authorId 作者ID
   */
  getByAuthorId: (authorId: string) => {
    return instance.get(`/asset/getAssetByAuthorID?authorId=${authorId}`);
  },

  /**
   * 根据拥有者ID获取资产
   * @param ownerId 拥有者ID
   */
  getByOwnerId: (ownerId: string) => {
    return instance.get(`/asset/getAssetByOwnerID?ownerId=${ownerId}`);
  },

  /**
   * 获取资产状态
   * @param id 资产ID
   */
  getStatus: (id: string) => {
    return instance.get(`/asset/getStatus?id=${id}`);
  }

};

// 钱包相关API
const walletApi = {
  /**
   * 获取当前用户余额
   */
  getBalance: () => {
    return instance.get('/wallet/balance');
  },

  /**
   * 转账
   * @param recipientId 接收者ID
   * @param amount 转账金额
   */
  transfer: (recipientId: number, amount: number) => {
    return instance.post('/wallet/transfer', {
      recipientId: recipientId,
      amount: amount
    });
  },

  /**
   * 铸币（仅金融组织）
   * @param accountId 目标账户ID
   * @param amount 铸币金额
   */
  mintToken: (accountId: number, amount: number) => {
    return instance.post('/wallet/mintToken', {
      accountID: accountId,
      amount: amount
    });
  },

 /**
   * 获取转出记录（匹配后端/wallet/transferBySenderID）
   */
  getTransfersBySender: () => {
    return instance.get('/wallet/transferBySenderID');
  },

  /**
   * 获取转入记录（匹配后端/wallet/transferByRecipientID）
   */
  getTransfersByRecipient: () => {
    return instance.get('/wallet/transferByRecipientID');
  },

  /**
   * 获取当前用户的NFT资产（匹配后端/asset/getAssetByOwnerID）
   */
  getAssetsByOwner: () => {
    return instance.get('/asset/getAssetByOwnerID');
  },

  /**
   * 获取当前用户所属组织
   */
  getCurrentOrg: () => {
  return instance.get('/account/profile').then(user => {
    // 假设用户信息中组织字段为org
    return user.data.org;
    });
  },

 /**
   * 获取当前用户预扣款
   */
  getWithholdingsByAccount: () => {
    return instance.get('/wallet/getWithHoldingByAccountID');
  },
  getWithholdingsByListing: (listingId: string) => {
    return instance.get(`/wallet/getWithHoldingByListingID?listingID=${listingId}`);
  }
};

// 聊天相关API
const chatApi = {
  /**
   * 获取聊天会话列表
   */
  getChatSessions: () => {
    return instance.get('/chat/getChatSession');
  },

  /**
   * 获取与指定用户的消息记录
   * @param otherId 对方用户ID
   */
  getMessages: (otherId: number) => {
    return instance.get(`/chat/getMessages?otherID=${otherId}`);
  },

  /**
   * 标记消息为已读
   * @param otherId 对方用户ID
   */
  readMessages: (otherId: number) => {
    return instance.post(`/chat/readMessages?otherID=${otherId}`);
  },

  /**
   * 获取未读消息数量
   * @param otherId 对方用户ID
   */
  getUnreadMessageCount: (otherId: number) => {
    return instance.get(`/chat/getUnreadMessageCount?otherID=${otherId}`);
  }
};

// 市场相关API
const marketApi = {
  /**
   * 创建挂牌
   * @param listingData 挂牌数据 {assetId, title, price, deadline}
   */
  createListing: (listingData: any) => {
    return instance.post('/market/listing', listingData);
  },
  /**
   * 查询在售挂牌
   * 后端路由：GET /market/listings
   * @param params { page?: number; pageSize?: number }
   */
  list: (params?: { page?: number; pageSize?: number }) => {
    return instance.get('/market/listings', { params });
  },
  // 提交出价
  createOffer: (offerData: { listingId: number | string; offerPrice: number }) =>
    instance.post('/market/offer', offerData),

  // 卖家接受出价
  acceptOffer: (offerId: number | string) =>
    instance.post(`/market/offer/${offerId}/accept`),

  // 撤回出价
  cancelOffer: (offerId: number | string) =>
    instance.post(`/market/offer/${offerId}/cancel`),

  // 我的出价
  listMyOffers: (params?: any) => instance.get('/market/offers/mine', { params }),
  buyNow: (data: { listingId: number | string }) => instance.post('/market/buyNow', data)
};

// 拍卖相关API
const auctionApi = {
  /**
   * 创建拍卖品
   * @param auctionData 拍卖品数据 {assetId, title, reservePrice, startTime, deadline}
   */
  create: (auctionData: any) => {
    return instance.post('/auction/create', auctionData);
  },
  
  /**
   * 列出所有拍卖品
   */
  list: () => {
    return instance.get('/auction/list');
  },
  
  /**
   * 按照出售者获取拍卖品
   */
  getBySeller: () => {
    return instance.get('/auction/seller');
  },

  /**
   * 出价
   * @param bidData 出价数据 {lotId, bidPrice}
   */
  bid: (bidData: any) => {
    return instance.post('/auction/bid', bidData);
  },

  /**
   * 获取当前用户的出价
   * @param lotId 拍卖品ID
   */
  getBid: (lotId: number) => {
    return instance.get(`/auction/bid?lotID=${lotId}`);
  },

  /**
   * 获取拍卖结果
   * @param lotId 拍卖品ID
   */
  getResult: (lotId: number) => {
    return instance.get(`/auction/result?lotID=${lotId}`);
  }
};

// 导出所有API模块
export { accountApi, assetApi, walletApi, chatApi, auctionApi, marketApi };

// 默认导出包含所有API的对象
export default {
  account: accountApi,
  asset: assetApi,
  wallet: walletApi,
  chat: chatApi,
  auction: auctionApi
};