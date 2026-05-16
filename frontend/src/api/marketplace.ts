import { apiClient } from './client'
import type { MarketplaceGroup, MarketplaceStats } from '@/types'

export async function getMarketplaceModels(options?: { signal?: AbortSignal }): Promise<MarketplaceGroup[]> {
  const { data } = await apiClient.get<MarketplaceGroup[]>('/marketplace/models', {
    signal: options?.signal,
  })
  return data
}

export async function getMarketplaceStats(options?: { signal?: AbortSignal }): Promise<MarketplaceStats> {
  const { data } = await apiClient.get<MarketplaceStats>('/marketplace/stats', {
    signal: options?.signal,
  })
  return data
}

export const marketplaceAPI = {
  getMarketplaceModels,
  getMarketplaceStats,
}

export default marketplaceAPI
