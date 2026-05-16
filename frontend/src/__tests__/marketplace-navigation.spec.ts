import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { describe, expect, it } from 'vitest'

const srcRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..')
const routerSource = readFileSync(resolve(srcRoot, 'router/index.ts'), 'utf8')
const sidebarSource = readFileSync(resolve(srcRoot, 'components/layout/AppSidebar.vue'), 'utf8')
const marketplaceViewPath = resolve(srcRoot, 'views/ModelMarketplaceView.vue')

describe('模型广场导航入口', () => {
  it('公开 /models 并将旧可用渠道地址重定向到模型广场', () => {
    expect(routerSource).toContain("path: '/models'")
    expect(routerSource).toContain("name: 'ModelMarketplace'")
    expect(routerSource).toContain("component: () => import('@/views/ModelMarketplaceView.vue')")
    expect(routerSource).toContain("path: '/available-channels'")
    expect(routerSource).toContain("redirect: '/models'")
  })

  it('不再暴露用户渠道状态和管理员渠道监控路由', () => {
    expect(routerSource).not.toContain("path: '/monitor'")
    expect(routerSource).not.toContain("path: '/admin/channels/monitor'")
    expect(routerSource).not.toContain('ChannelStatusView.vue')
    expect(routerSource).not.toContain('ChannelMonitorView.vue')
  })

  it('侧边栏展示模型广场和模型定价入口', () => {
    expect(sidebarSource).toContain("path: '/models'")
    expect(sidebarSource).toContain("t('nav.modelMarketplace')")
    expect(sidebarSource).toContain("path: '/admin/channels/model-pricing'")
    expect(sidebarSource).toContain("t('nav.modelPricing')")
  })

  it('模型广场页面调用公开 marketplace API 并展示容量信息', () => {
    const marketplaceViewSource = readFileSync(marketplaceViewPath, 'utf8')

    expect(marketplaceViewSource).toContain('marketplaceAPI.getMarketplaceModels')
    expect(marketplaceViewSource).toContain('marketplaceAPI.getMarketplaceStats')
    expect(marketplaceViewSource).toContain('GroupCapacityBadge')
    expect(marketplaceViewSource).toContain('model_marketplace_enabled')
  })
})
