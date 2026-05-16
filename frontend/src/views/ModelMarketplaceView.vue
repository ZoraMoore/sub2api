<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="overflow-hidden rounded-3xl border border-slate-200 bg-slate-950 text-white shadow-xl dark:border-dark-700">
        <div class="relative p-6 sm:p-8 lg:p-10">
          <div class="absolute inset-0 opacity-40">
            <div class="absolute -left-24 top-0 h-72 w-72 rounded-full bg-cyan-400 blur-3xl"></div>
            <div class="absolute right-0 top-12 h-72 w-72 rounded-full bg-emerald-400 blur-3xl"></div>
            <div class="absolute bottom-0 left-1/3 h-48 w-48 rounded-full bg-blue-500 blur-3xl"></div>
          </div>

          <div class="relative grid gap-8 lg:grid-cols-[1.3fr_0.7fr] lg:items-end">
            <div class="space-y-4">
              <div class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-3 py-1 text-xs font-medium text-cyan-100 backdrop-blur">
                <Icon name="grid" size="xs" />
                {{ t('marketplace.badge') }}
              </div>
              <div class="space-y-3">
                <h1 class="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl">
                  {{ t('marketplace.title') }}
                </h1>
                <p class="max-w-3xl text-sm leading-6 text-slate-200 sm:text-base">
                  {{ t('marketplace.description') }}
                </p>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-2">
              <div v-for="item in summaryCards" :key="item.label" class="rounded-2xl border border-white/10 bg-white/10 p-4 backdrop-blur">
                <p class="text-xs text-slate-300">{{ item.label }}</p>
                <p class="mt-2 text-2xl font-semibold text-white">{{ item.value }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section v-if="!modelMarketplaceEnabled" class="rounded-2xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-200">
        {{ t('marketplace.disabledHint') }}
      </section>

      <section class="card p-4 sm:p-5">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div class="relative w-full lg:max-w-md">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500" />
            <input
              v-model="searchQuery"
              type="text"
              class="input pl-10"
              :placeholder="t('marketplace.searchPlaceholder')"
            />
          </div>
          <button type="button" class="btn btn-secondary inline-flex items-center gap-2" :disabled="loading" @click="loadMarketplace">
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            {{ t('marketplace.refresh') }}
          </button>
        </div>
      </section>

      <section v-if="loading && groups.length === 0" class="grid gap-4 lg:grid-cols-2">
        <div v-for="index in 4" :key="index" class="card animate-pulse p-5">
          <div class="h-5 w-40 rounded bg-gray-200 dark:bg-dark-700"></div>
          <div class="mt-4 h-4 w-2/3 rounded bg-gray-100 dark:bg-dark-800"></div>
          <div class="mt-6 grid gap-2 sm:grid-cols-2">
            <div class="h-16 rounded-xl bg-gray-100 dark:bg-dark-800"></div>
            <div class="h-16 rounded-xl bg-gray-100 dark:bg-dark-800"></div>
          </div>
        </div>
      </section>

      <section v-else-if="filteredGroups.length === 0" class="card p-10 text-center">
        <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-gray-100 text-gray-500 dark:bg-dark-800 dark:text-dark-400">
          <Icon name="grid" size="lg" />
        </div>
        <h2 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('marketplace.empty') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('marketplace.emptyHint') }}</p>
      </section>

      <section v-else class="grid gap-4 lg:grid-cols-2">
        <article v-for="group in filteredGroups" :key="group.id" class="card overflow-hidden">
          <div class="border-b border-gray-100 p-5 dark:border-dark-700">
            <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <h2 class="truncate text-lg font-semibold text-gray-900 dark:text-white">{{ group.name }}</h2>
                  <span class="rounded-full bg-primary-50 px-2.5 py-1 text-xs font-medium text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
                    {{ group.platform }}
                  </span>
                  <span v-if="group.display_brand" class="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-dark-300">
                    {{ group.display_brand }}
                  </span>
                </div>
                <p v-if="group.description" class="mt-2 line-clamp-2 text-sm text-gray-500 dark:text-gray-400">
                  {{ group.description }}
                </p>
              </div>

              <div class="flex flex-shrink-0 flex-col items-start gap-2 sm:items-end">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('marketplace.capacity') }}</span>
                <GroupCapacityBadge
                  v-if="group.capacity"
                  layout="horizontal"
                  :concurrency-used="group.capacity.concurrency_used"
                  :concurrency-max="group.capacity.concurrency_max"
                  :sessions-used="group.capacity.sessions_used"
                  :sessions-max="group.capacity.sessions_max"
                  :rpm-used="group.capacity.rpm_used"
                  :rpm-max="group.capacity.rpm_max"
                />
                <span v-else class="text-xs text-gray-400">{{ t('marketplace.noCapacity') }}</span>
              </div>
            </div>
          </div>

          <div class="space-y-3 p-5">
            <div v-for="model in group.models" :key="`${group.id}-${model.id}`" class="rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-900/50">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div class="min-w-0">
                  <p class="truncate font-medium text-gray-900 dark:text-white">{{ model.display_name || model.id }}</p>
                  <p class="mt-1 truncate font-mono text-xs text-gray-500 dark:text-gray-400">{{ model.id }}</p>
                </div>
                <span :class="priceStatusClass(model.pricing.price_status)">
                  {{ priceStatusLabel(model.pricing.price_status) }}
                </span>
              </div>

              <div class="mt-4 flex flex-wrap gap-2">
                <span v-for="price in pricingBadges(model.pricing)" :key="price.label" class="rounded-lg border border-gray-200 bg-white px-2.5 py-1 text-xs text-gray-700 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200">
                  <span class="text-gray-500 dark:text-dark-400">{{ price.label }}:</span>
                  <span class="ml-1 font-mono font-semibold">{{ price.value }}</span>
                </span>
              </div>
            </div>
          </div>
        </article>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import GroupCapacityBadge from '@/components/common/GroupCapacityBadge.vue'
import { marketplaceAPI } from '@/api/marketplace'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import type { MarketplaceGroup, MarketplaceModelPricing, MarketplacePriceStatus } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const groups = ref<MarketplaceGroup[]>([])
const stats = ref({ today_tokens: 0, total_tokens: 0, total_users: 0 })
const loading = ref(false)
const searchQuery = ref('')

const modelMarketplaceEnabled = computed(() => appStore.cachedPublicSettings?.model_marketplace_enabled ?? true)

const filteredGroups = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return groups.value

  return groups.value
    .map((group) => {
      const groupHit = [group.name, group.description, group.platform, group.display_brand]
        .filter(Boolean)
        .some((value) => String(value).toLowerCase().includes(query))
      if (groupHit) return group

      const models = group.models.filter((model) =>
        [model.id, model.display_name].some((value) => value.toLowerCase().includes(query)),
      )
      return models.length > 0 ? { ...group, models, model_count: models.length } : null
    })
    .filter((group): group is MarketplaceGroup => group !== null)
})

const totalModels = computed(() => groups.value.reduce((sum, group) => sum + group.model_count, 0))

const summaryCards = computed(() => [
  { label: t('marketplace.totalGroups'), value: formatCompactNumber(groups.value.length) },
  { label: t('marketplace.totalModels'), value: formatCompactNumber(totalModels.value) },
  { label: t('marketplace.todayTokens'), value: formatCompactNumber(stats.value.today_tokens) },
  { label: t('marketplace.totalUsers'), value: formatCompactNumber(stats.value.total_users) },
])

async function loadMarketplace() {
  loading.value = true
  try {
    const [modelGroups, modelStats] = await Promise.all([
      marketplaceAPI.getMarketplaceModels(),
      marketplaceAPI.getMarketplaceStats(),
    ])
    groups.value = modelGroups
    stats.value = modelStats
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('marketplace.loadFailed')))
  } finally {
    loading.value = false
  }
}

function pricingBadges(pricing: MarketplaceModelPricing): Array<{ label: string; value: string }> {
  const badges: Array<{ label: string; value: string }> = []
  addPriceBadge(badges, t('marketplace.pricing.input'), pricing.input_price_per_token, t('marketplace.pricing.unitPerMillion'))
  addPriceBadge(badges, t('marketplace.pricing.output'), pricing.output_price_per_token, t('marketplace.pricing.unitPerMillion'))
  addPriceBadge(badges, t('marketplace.pricing.cacheWrite'), pricing.cache_write_price_per_token, t('marketplace.pricing.unitPerMillion'))
  addPriceBadge(badges, t('marketplace.pricing.cacheRead'), pricing.cache_read_price_per_token, t('marketplace.pricing.unitPerMillion'))
  addPriceBadge(badges, t('marketplace.pricing.imageOutput'), pricing.image_output_price_per_token, t('marketplace.pricing.unitPerMillion'))
  addPriceBadge(badges, t('marketplace.pricing.perRequest'), pricing.per_request_price, t('marketplace.pricing.unitPerRequest'))
  addPriceBadge(badges, t('marketplace.pricing.image1k'), pricing.image_price_1k, t('marketplace.pricing.unitPerRequest'))
  addPriceBadge(badges, t('marketplace.pricing.image2k'), pricing.image_price_2k, t('marketplace.pricing.unitPerRequest'))
  addPriceBadge(badges, t('marketplace.pricing.image4k'), pricing.image_price_4k, t('marketplace.pricing.unitPerRequest'))

  if (badges.length === 0) {
    badges.push({ label: t('marketplace.pricing.mode'), value: pricingModeLabel(pricing.pricing_mode) })
  }

  if (pricing.context_intervals?.length) {
    badges.push({ label: t('marketplace.pricing.intervals'), value: String(pricing.context_intervals.length) })
  }

  return badges
}

function addPriceBadge(badges: Array<{ label: string; value: string }>, label: string, value: number | undefined, unit: string) {
  if (value === undefined || value === null || Number(value) <= 0) return
  badges.push({ label, value: `${formatPrice(Number(value) * 1_000_000)} ${unit}` })
}

function priceStatusLabel(status: MarketplacePriceStatus | string): string {
  return status === 'priced' ? t('marketplace.priceStatus.priced') : t('marketplace.priceStatus.unpriced')
}

function priceStatusClass(status: MarketplacePriceStatus | string): string {
  return status === 'priced'
    ? 'inline-flex rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
    : 'inline-flex rounded-full bg-amber-50 px-2.5 py-1 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
}

function pricingModeLabel(mode: string): string {
  if (mode === 'token') return t('marketplace.pricing.token')
  if (mode === 'image') return t('marketplace.pricing.image')
  return t('marketplace.pricing.unknown')
}

function formatPrice(value: number): string {
  return new Intl.NumberFormat(undefined, { maximumFractionDigits: 6 }).format(value)
}

function formatCompactNumber(value: number): string {
  return new Intl.NumberFormat(undefined, { notation: 'compact', maximumFractionDigits: 1 }).format(value || 0)
}

onMounted(loadMarketplace)
</script>
