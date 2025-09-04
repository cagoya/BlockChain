export const transformOrg = (org: number[]): string => {
    const orgMap: { [key: number]: string } = {
      1: '平台运营',
      2: '创作者',
      3: '金融机构',
    };
  
    const orgNames = org.map((item) => {
      return orgMap[item] || ''; // 使用映射表，如果找不到则返回空字符串
    });
  
    return orgNames.join(' '); // 使用 join 方法将数组元素用空格连接
};