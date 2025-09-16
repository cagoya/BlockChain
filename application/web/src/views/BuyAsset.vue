<template>
  <MarketNav />

  <div class="page">
    <!-- 头部 -->
    <header class="page-header">
      <div class="header-left">
        <h1>市场购买</h1>
        <p class="sub">浏览上架的 NFT 资产并立即购买</p>
      </div>
      <div class="header-right">
        <a-button type="default" size="small" @click="fetchListings" :loading="loading">
          刷新列表
        </a-button>
      </div>
    </header>

    <!-- 加载态 -->
    <div v-if="loading" class="loading">
      <a-skeleton :paragraph="{ rows: 4 }" active style="max-width: 680px; width: 100%" />
      <div class="skeleton-grid">
        <a-skeleton-image style="width: 100%; height: 180px; border-radius: 12px" />
        <a-skeleton-image style="width: 100%; height: 180px; border-radius: 12px" />
        <a-skeleton-image style="width: 100%; height: 180px; border-radius: 12px" />
      </div>
    </div>

    <!-- 列表 -->
    <div v-else-if="listings.length" class="grid">
      <a-card
        v-for="it in listings"
        :key="it.id"
        :hoverable="true"
        class="nft-card"
        @click="openBuy(it)"
      >
        <!-- 顶部图 -->
        <div class="cover">
          <img
            v-if="it.asset?.imageName"
            class="cover-img"
            :src="getImageURL(it.asset.imageName)"
            :alt="it.asset?.name || it.title || 'asset'"
          />
          <div v-else class="cover-placeholder">
            <PictureOutlined style="font-size:16px;margin-right:6px" />
            无图片
          </div>

          <!-- 价格角标 -->
          <div class="price-badge">
            <span class="price">{{ it.price }}</span>
            <span class="unit">代币</span>
          </div>
        </div>

        <!-- 主体信息 -->
        <div class="body">
          <div class="title-row">
            <!-- 标题区域 -->
            <a-tooltip :title="displayTitleAssetFirst(it)">
                <h3 class="title">{{ displayTitleAssetFirst(it) }}</h3>
            </a-tooltip>
            <a-tag color="blue" v-if="it.asset?.name">NFT</a-tag>
          </div>

          <div class="meta">
            <div class="meta-item">
              <span class="label">资产ID</span>
              <code class="value id">{{ it.assetId }}</code>
            </div>
            <div class="meta-item" v-if="it.deadline">
              <span class="label">截止</span>
              <span class="value">{{ formatTime(it.deadline) }}</span>
            </div>
            <div class="meta-item">
              <span class="label">创作者</span>
              <span class="value">{{ it.asset?.authorName || it.asset?.authorId || '未知' }}</span>
            </div>
          </div>
        </div>

        <!-- 行为区 -->
        <div class="actions">
          <a-button type="primary" block @click.stop="openBuy(it)">立即购买</a-button>
        </div>
      </a-card>
    </div>

    <!-- 空态 -->
    <div v-else class="empty">
      <a-empty description="暂无在售资产" />
    </div>

    <!-- 购买确认弹窗（用 visible，与项目其他页面保持一致） -->
    <a-modal
    :visible="confirm.open"
    @update:visible="(v) => (confirm.open = v)"
    title="确认购买"
    ok-text="确认"
    cancel-text="取消"
    :confirm-loading="buying"
    :ok-button-props="{ disabled: insufficient }"
    @ok="doPurchase"
    >
    <div class="confirm">
        <!-- 左侧封面 -->
        <div class="confirm-left">
        <div class="cover">
            <img
            v-if="confirm.asset?.imageName"
            class="cover-img"
            :src="getImageURL(confirm.asset.imageName)"
            :alt="confirm.title"
            />
            <div v-else class="cover-placeholder">无图片</div>
        </div>
        </div>

        <!-- 右侧信息 -->
        <div class="confirm-right">
        <div class="confirm-row">
            <span class="label">资产</span>
            <span class="value">{{ confirm.title }}</span>
        </div>
        <div class="confirm-row">
            <span class="label">价格</span>
            <span class="value price">{{ confirm.price }} 代币</span>
        </div>
        <div class="confirm-row">
            <span class="label">资产ID</span>
            <code class="value mono">{{ confirm.assetId }}</code>
        </div>
        <div class="confirm-row" v-if="confirm.asset?.authorName || confirm.asset?.authorId">
            <span class="label">创作者</span>
            <span class="value">{{ confirm.asset?.authorName || confirm.asset?.authorId }}</span>
        </div>

        <!-- 余额提示（可选） -->
        <a-alert
            v-if="insufficient"
            type="warning"
            show-icon
            style="margin-top: 8px"
            message="余额不足"
            description="请先铸币或减少出价金额。"
        />
        </div>
    </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import MarketNav from '../components/MarketNav.vue'
import { marketApi, assetApi, getImageURL,accountApi,walletApi  } from '../api'
import { PictureOutlined } from '@ant-design/icons-vue'
interface Asset {
  id: string
  name: string
  imageName: string
  timeStamp?: string
  authorId?: number
  authorName?: string
}

interface Listing {
  id: number | string
  assetId: string
  title: string
  price: number
  deadline?: string | null
  asset?: Asset | null
}

const loading = ref(false)
const buying = ref(false)
const listings = ref<Listing[]>([])
const balance = ref<number | null>(null)
const insufficient = computed(() => balance.value != null && confirm.value.price > (balance.value || 0))
const confirm = ref<{
  open: boolean
  id: number | string | null
  title: string
  price: number
  assetId: string | null
  asset: Asset | null
}>({
  open: false,
  id: null,
  title: '',
  price: 0,
  assetId: null,
  asset: null
})
async function refreshBalance() {
  try {
    const r = await walletApi.getBalance()
    const b = r.data?.balance ?? r.data?.data ?? r.data
    balance.value = Number(b) || 0
  } catch { /* 忽略错误 */ }
}
function normalizeItems(respData: any): Listing[] {
  const data = respData?.data ?? respData
  const items = Array.isArray(data) ? data : (data?.items ?? [])
  return Array.isArray(items) ? items : []
}
// 统一的标题取值：title -> Title -> asset.name
function displayTitleAssetFirst(l: any): string {
  // 先用资产名，再退回挂牌标题（兼容 Title/title）
  const assetName = l?.asset?.name
  if (assetName && String(assetName).trim() !== '') return assetName
  const t = l?.title ?? l?.Title
  return (t && String(t).trim() !== '') ? String(t) : '无标题'
}
async function fetchListings() {
  loading.value = true
  try {
    const resp = await marketApi.list({ page: 1, pageSize: 50 })
    let base: Listing[] = []
    if (resp.data?.code === 200) {
      base = normalizeItems(resp.data)
    } else {
      base = normalizeItems(resp.data)
      if (!base.length) message.error(resp.data?.message || '获取在售列表失败')
    }

    // 仅为显示图片补充 asset
    const enriched = await Promise.all(
      base.map(async (l) => {
        try {
          if (l.assetId) {
            const ar = await assetApi.getById(l.assetId)
            if (ar.data?.code === 200 && ar.data?.data) {
              l.asset = ar.data.data
              // 👇 新增：取作者名字
              if (l.asset?.authorId) {
                const u = await accountApi.getUserNameById(l.asset.authorId)
                if (u.data?.code === 200) {
                  l.asset.authorName = u.data.data // 例如 "张三"
                }
              }
            }
          }
        } catch { /* 忽略单条失败 */ }
        return l
      })
    )
    listings.value = enriched
  } catch {
    message.error('获取在售列表失败')
    listings.value = []
  } finally {
    loading.value = false
  }
}

function formatTime(s?: string | null) {
  if (!s) return ''
  const d = new Date(s)
  return isNaN(d.getTime()) ? s : d.toLocaleString('zh-CN', { hour12: false })
}

function openBuy(it: any) {
  confirm.value = {
    open: true,
    id: it.id,
    title: (it.asset?.name && String(it.asset.name).trim()) ? it.asset.name : (it.title || '无标题'),
    price: it.price,
    assetId: it.assetId,
    asset: it.asset || null
  }
  // 打开时顺便刷新余额，用于不足判断
  refreshBalance()
}

async function doPurchase() {
  if (confirm.value.id == null) return
  buying.value = true
  try {
    const r = await marketApi.buyNow({ listingId: Number(confirm.value.id) })
    // 打印所有关键信息，便于定位
    console.error('buyNow resp:', r.status, r.data)

    if (r.data?.code !== 200) {
      message.error(r.data?.message || '购买失败')
      return
    }
    message.success('购买成功')
    confirm.value.open = false
    await Promise.all([refreshBalance(), fetchListings()])
  }  catch (err: any) {
  console.error('buyNow error:', err)
  message.error(err?.message || '购买失败')
  } finally {
    buying.value = false
  }
}

onMounted(fetchListings)
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f8fb;
  padding: 16px 20px 28px;
  box-sizing: border-box;
}

/* 头部 */
.page-header {
  max-width: 1240px;
  margin: 0 auto 12px auto;
  padding: 14px 18px;
  background: #fff;
  border: 1px solid #eef0f6;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  box-shadow: 0 6px 24px rgba(20, 26, 56, 0.05);
}

.header-left h1 {
  font-size: 22px;
  color: #1a237e;
  margin: 0;
}
.header-left .sub {
  margin: 4px 0 0 0;
  color: #607d8b;
  font-size: 13px;
}

/* 加载骨架 */
.loading {
  max-width: 1240px;
  margin: 24px auto 0;
}
.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 14px;
  margin-top: 16px;
}

/* 列表网格 */
.grid {
  max-width: 1240px;
  margin: 16px auto 0 auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}

/* 卡片 */
.nft-card {
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid #eef0f6;
  box-shadow: 0 10px 30px rgba(44, 73, 255, 0.06);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}
.nft-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 16px 36px rgba(44, 73, 255, 0.14);
}

/* 封面图 */
.cover {
  position: relative;
  height: 180px;
  background: linear-gradient(135deg, #f5f7ff, #eef2ff);
  overflow: hidden;
}
.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform .35s ease;
}
.nft-card:hover .cover-img { transform: scale(1.03); }
.cover-placeholder {
  width: 100%;
  height: 100%;
  color: #97a0b3;
  font-size: 12px;
  display:flex; align-items:center; justify-content:center;
}

/* 价格角标 */
.price-badge {
  position: absolute;
  left: 12px;
  bottom: 12px;
  background: rgba(26,35,126,.92);
  color:#fff;
  border-radius: 10px;
  padding: 6px 10px;
  display:flex; align-items: baseline; gap: 6px;
  box-shadow: 0 6px 18px rgba(26,35,126,.25);
}
.price-badge .price { font-size: 18px; font-weight: 800; letter-spacing: .3px; }
.price-badge .unit { font-size: 12px; opacity: .9; }

/* 主体信息 */
.body {
  padding: 14px 14px 0 14px;
}
.title-row {
  display: flex; align-items: center; justify-content: space-between; gap: 8px;
}
.title {
  margin: 0;
  font-size: 16px;
  color: #1a237e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.meta {
  margin-top: 10px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 6px;
}
.meta-item { display: flex; gap: 8px; align-items: center; }
.meta-item .label { color: #6b7280; font-size: 12px; }
.meta-item .value { color: #374151; font-size: 12px; }
.meta-item .id {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace;
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 6px;
}

/* 行为区 */
.actions {
  padding: 12px 14px 14px 14px;
}

/* 空态 */
.empty {
  max-width: 1240px;
  margin: 40px auto;
  display:flex; justify-content:center;
}

/* 弹窗 */
.confirm { display: grid; grid-template-columns: 160px 1fr; gap: 14px; }
.confirm-left .cover { width: 100%; height: 140px; border-radius: 10px; overflow: hidden; background: #f3f4f6; display:flex; align-items:center; justify-content:center; }
.cover-img { width: 100%; height: 100%; object-fit: cover; }
.cover-placeholder { color:#9ca3af; font-size:12px; }
.confirm-right .confirm-row { display:flex; gap:12px; align-items:center; margin:6px 0; }
.confirm-right .label { width: 56px; color:#6b7280; font-size:12px; flex: none; }
.confirm-right .value { color:#111827; font-size:14px; }
.confirm-right .value.price { color:#ef4444; font-weight:700; }
.mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace; background:#f3f4f6; padding:2px 6px; border-radius:6px; }
</style>