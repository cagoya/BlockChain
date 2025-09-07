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
    // 从本地存储获取token并添加到请求头
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
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
    return response.data;
  },
  (error) => {
    return Promise.reject(error.response?.data || '请求失败');
  }
);

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
      recipientID: recipientId,
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

export default walletApi;