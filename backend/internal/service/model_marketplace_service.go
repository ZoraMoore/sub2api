package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// ModelMarketplaceGroup 是公开模型广场按分组聚合后的展示结构。
type ModelMarketplaceGroup struct {
	ID             int64
	Name           string
	Description    string
	Platform       string
	DisplayBrand   string
	SortOrder      int
	RateMultiplier float64
	Capacity       *GroupCapacitySummary
	ModelCount     int
	Models         []ModelMarketplaceModel
}

// ModelMarketplaceModel 是模型广场中单个模型的展示结构。
type ModelMarketplaceModel struct {
	ID          string
	DisplayName string
	Pricing     ModelMarketplacePricing
}

// ModelMarketplacePricing 是公开模型广场展示用的模型价格快照。
type ModelMarketplacePricing struct {
	PricingMode              string
	PriceStatus              string
	InputPricePerToken       float64
	OutputPricePerToken      float64
	CacheWritePricePerToken  float64
	CacheReadPricePerToken   float64
	ImageOutputPricePerToken float64
	PerRequestPrice          float64
	ContextIntervals         []ModelMarketplacePricingInterval
	ImagePrice1K             float64
	ImagePrice2K             float64
	ImagePrice4K             float64
}

// ModelMarketplacePricingInterval 是模型广场上下文/分层定价区间。
type ModelMarketplacePricingInterval struct {
	MinTokens                int
	MaxTokens                *int
	TierLabel                string
	InputPricePerToken       float64
	OutputPricePerToken      float64
	CacheWritePricePerToken  float64
	CacheReadPricePerToken   float64
	ImageOutputPricePerToken float64
	PerRequestPrice          float64
}

// ModelMarketplaceService 组合渠道模型定价、公开分组和容量快照，生成公开模型广场数据。
type ModelMarketplaceService struct {
	channelService  *ChannelService
	groupRepo       GroupRepository
	capacityService *GroupCapacityService
}

func NewModelMarketplaceService(
	channelService *ChannelService,
	groupRepo GroupRepository,
	capacityService *GroupCapacityService,
) *ModelMarketplaceService {
	return &ModelMarketplaceService{
		channelService:  channelService,
		groupRepo:       groupRepo,
		capacityService: capacityService,
	}
}

// ListPublic 返回公开模型广场列表。
//
// 公开边界：只展示非专属、有活跃账号且至少存在可展示模型的分组；价格沿用
// ChannelService.ListAvailable 的展示定价结果，避免公开接口泄漏渠道内部字段。
func (s *ModelMarketplaceService) ListPublic(ctx context.Context) ([]ModelMarketplaceGroup, error) {
	if s == nil || s.channelService == nil || s.groupRepo == nil {
		return nil, fmt.Errorf("model marketplace dependencies are not initialized")
	}

	groups, err := s.groupRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active groups: %w", err)
	}
	publicGroups := make(map[int64]Group, len(groups))
	publicGroupIDs := make([]int64, 0, len(groups))
	for i := range groups {
		group := groups[i]
		if group.IsExclusive || group.ActiveAccountCount <= 0 {
			continue
		}
		publicGroups[group.ID] = group
		publicGroupIDs = append(publicGroupIDs, group.ID)
	}
	if len(publicGroups) == 0 {
		return []ModelMarketplaceGroup{}, nil
	}

	capacityMap := s.getPublicCapacityMap(ctx, publicGroupIDs)
	channels, err := s.channelService.ListAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("list available channels: %w", err)
	}

	byGroup := make(map[int64]*ModelMarketplaceGroup, len(publicGroups))
	seenModels := make(map[int64]map[string]struct{}, len(publicGroups))
	for _, channel := range channels {
		if channel.Status != StatusActive {
			continue
		}
		for _, ref := range channel.Groups {
			group, ok := publicGroups[ref.ID]
			if !ok {
				continue
			}
			entry := byGroup[group.ID]
			if entry == nil {
				entry = &ModelMarketplaceGroup{
					ID:             group.ID,
					Name:           group.Name,
					Description:    group.Description,
					Platform:       group.Platform,
					DisplayBrand:   marketplaceDisplayBrand(group),
					SortOrder:      group.SortOrder,
					RateMultiplier: group.RateMultiplier,
					Capacity:       marketplaceCapacity(capacityMap, group.ID),
					Models:         make([]ModelMarketplaceModel, 0),
				}
				byGroup[group.ID] = entry
				seenModels[group.ID] = make(map[string]struct{})
			}

			for _, model := range channel.SupportedModels {
				if model.Platform != group.Platform {
					continue
				}
				modelID := strings.TrimSpace(model.Name)
				if modelID == "" {
					continue
				}
				dedupKey := strings.ToLower(model.Platform + ":" + modelID)
				if _, exists := seenModels[group.ID][dedupKey]; exists {
					continue
				}
				seenModels[group.ID][dedupKey] = struct{}{}
				entry.Models = append(entry.Models, ModelMarketplaceModel{
					ID:          modelID,
					DisplayName: modelID,
					Pricing:     marketplacePricingFromChannel(model.Pricing, group),
				})
			}
		}
	}

	out := make([]ModelMarketplaceGroup, 0, len(byGroup))
	for _, entry := range byGroup {
		if len(entry.Models) == 0 {
			continue
		}
		sort.SliceStable(entry.Models, func(i, j int) bool {
			return strings.ToLower(entry.Models[i].ID) < strings.ToLower(entry.Models[j].ID)
		})
		entry.ModelCount = len(entry.Models)
		out = append(out, *entry)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].SortOrder != out[j].SortOrder {
			return out[i].SortOrder < out[j].SortOrder
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func (s *ModelMarketplaceService) getPublicCapacityMap(ctx context.Context, groupIDs []int64) map[int64]GroupCapacitySummary {
	if s.capacityService == nil || len(groupIDs) == 0 {
		return nil
	}
	capacityMap, err := s.capacityService.GetGroupCapacityByIDs(ctx, groupIDs)
	if err != nil {
		// 容量是模型广场辅助信息，失败时不阻断模型和价格展示。
		return nil
	}
	return capacityMap
}

func marketplaceCapacity(capacityMap map[int64]GroupCapacitySummary, groupID int64) *GroupCapacitySummary {
	if len(capacityMap) == 0 {
		return nil
	}
	capacity, ok := capacityMap[groupID]
	if !ok {
		return nil
	}
	return &capacity
}

func marketplaceDisplayBrand(group Group) string {
	if brand := strings.TrimSpace(group.Name); brand != "" {
		return brand
	}
	return group.Platform
}

func marketplacePricingFromChannel(pricing *ChannelModelPricing, group Group) ModelMarketplacePricing {
	out := ModelMarketplacePricing{
		PricingMode:  "unknown",
		PriceStatus:  "unpriced",
		ImagePrice1K: ptrValue(group.ImagePrice1K),
		ImagePrice2K: ptrValue(group.ImagePrice2K),
		ImagePrice4K: ptrValue(group.ImagePrice4K),
	}
	if pricing == nil {
		out.PriceStatus = marketplacePriceStatus(out)
		return out
	}

	if pricing.BillingMode == "" {
		out.PricingMode = string(BillingModeToken)
	} else {
		out.PricingMode = string(pricing.BillingMode)
	}
	out.InputPricePerToken = ptrValue(pricing.InputPrice)
	out.OutputPricePerToken = ptrValue(pricing.OutputPrice)
	out.CacheWritePricePerToken = ptrValue(pricing.CacheWritePrice)
	out.CacheReadPricePerToken = ptrValue(pricing.CacheReadPrice)
	out.ImageOutputPricePerToken = ptrValue(pricing.ImageOutputPrice)
	out.PerRequestPrice = ptrValue(pricing.PerRequestPrice)
	out.ContextIntervals = marketplaceIntervalsFromChannel(pricing.Intervals)
	applyImageIntervalPrices(&out, pricing.Intervals)
	out.PriceStatus = marketplacePriceStatus(out)
	return out
}

func marketplaceIntervalsFromChannel(intervals []PricingInterval) []ModelMarketplacePricingInterval {
	out := make([]ModelMarketplacePricingInterval, 0, len(intervals))
	for _, interval := range intervals {
		out = append(out, ModelMarketplacePricingInterval{
			MinTokens:               interval.MinTokens,
			MaxTokens:               interval.MaxTokens,
			TierLabel:               interval.TierLabel,
			InputPricePerToken:      ptrValue(interval.InputPrice),
			OutputPricePerToken:     ptrValue(interval.OutputPrice),
			CacheWritePricePerToken: ptrValue(interval.CacheWritePrice),
			CacheReadPricePerToken:  ptrValue(interval.CacheReadPrice),
			PerRequestPrice:         ptrValue(interval.PerRequestPrice),
		})
	}
	return out
}

func applyImageIntervalPrices(out *ModelMarketplacePricing, intervals []PricingInterval) {
	if out == nil {
		return
	}
	for _, interval := range intervals {
		price := ptrValue(interval.PerRequestPrice)
		if price <= 0 {
			continue
		}
		switch strings.ToUpper(strings.TrimSpace(interval.TierLabel)) {
		case "1K":
			out.ImagePrice1K = price
		case "2K":
			out.ImagePrice2K = price
		case "4K":
			out.ImagePrice4K = price
		}
	}
}

func marketplacePriceStatus(pricing ModelMarketplacePricing) string {
	if pricing.InputPricePerToken > 0 || pricing.OutputPricePerToken > 0 ||
		pricing.CacheWritePricePerToken > 0 || pricing.CacheReadPricePerToken > 0 ||
		pricing.ImageOutputPricePerToken > 0 || pricing.PerRequestPrice > 0 ||
		pricing.ImagePrice1K > 0 || pricing.ImagePrice2K > 0 || pricing.ImagePrice4K > 0 {
		return "priced"
	}
	for _, interval := range pricing.ContextIntervals {
		if interval.InputPricePerToken > 0 || interval.OutputPricePerToken > 0 ||
			interval.CacheWritePricePerToken > 0 || interval.CacheReadPricePerToken > 0 ||
			interval.ImageOutputPricePerToken > 0 || interval.PerRequestPrice > 0 {
			return "priced"
		}
	}
	return "unpriced"
}

func ptrValue(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
