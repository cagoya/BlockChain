import axios from 'axios';

// 创建axios实例
const instance = axios.create({
  baseURL: '/api', // 基础API路径
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
}
};

// 导出所有API模块
export { accountApi, assetApi, walletApi };

// 默认导出包含所有API的对象
export default {
  account: accountApi,
  asset: assetApi,
  wallet: walletApi
};