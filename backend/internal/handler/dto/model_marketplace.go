package dto

import "github.com/Wei-Shaw/sub2api/internal/service"

type ModelMarketplaceStats struct {
	TodayTokens int64 `json:"today_tokens"`
	TotalTokens int64 `json:"total_tokens"`
	TotalUsers  int64 `json:"total_users"`
}

type ModelMarketplacePricing struct {
	PricingMode              string                            `json:"pricing_mode"`
	PriceStatus              string                            `json:"price_status"`
	InputPricePerToken       float64                           `json:"input_price_per_token,omitempty"`
	OutputPricePerToken      float64                           `json:"output_price_per_token,omitempty"`
	CacheWritePricePerToken  float64                           `json:"cache_write_price_per_token,omitempty"`
	CacheReadPricePerToken   float64                           `json:"cache_read_price_per_token,omitempty"`
	ImageOutputPricePerToken float64                           `json:"image_output_price_per_token,omitempty"`
	PerRequestPrice          float64                           `json:"per_request_price,omitempty"`
	ContextIntervals         []ModelMarketplacePricingInterval `json:"context_intervals,omitempty"`
	ImagePrice1K             float64                           `json:"image_price_1k,omitempty"`
	ImagePrice2K             float64                           `json:"image_price_2k,omitempty"`
	ImagePrice4K             float64                           `json:"image_price_4k,omitempty"`
}

type ModelMarketplacePricingInterval struct {
	MinTokens                int     `json:"min_tokens"`
	MaxTokens                *int    `json:"max_tokens,omitempty"`
	TierLabel                string  `json:"tier_label,omitempty"`
	InputPricePerToken       float64 `json:"input_price_per_token,omitempty"`
	OutputPricePerToken      float64 `json:"output_price_per_token,omitempty"`
	CacheWritePricePerToken  float64 `json:"cache_write_price_per_token,omitempty"`
	CacheReadPricePerToken   float64 `json:"cache_read_price_per_token,omitempty"`
	ImageOutputPricePerToken float64 `json:"image_output_price_per_token,omitempty"`
	PerRequestPrice          float64 `json:"per_request_price,omitempty"`
}

type ModelMarketplaceModel struct {
	ID          string                  `json:"id"`
	DisplayName string                  `json:"display_name"`
	Pricing     ModelMarketplacePricing `json:"pricing"`
}

type ModelMarketplaceCapacity struct {
	ConcurrencyUsed int `json:"concurrency_used"`
	ConcurrencyMax  int `json:"concurrency_max"`
	SessionsUsed    int `json:"sessions_used"`
	SessionsMax     int `json:"sessions_max"`
	RPMUsed         int `json:"rpm_used"`
	RPMMax          int `json:"rpm_max"`
}

type ModelMarketplaceGroup struct {
	ID             int64                     `json:"id"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	Platform       string                    `json:"platform"`
	DisplayBrand   string                    `json:"display_brand"`
	SortOrder      int                       `json:"sort_order"`
	RateMultiplier float64                   `json:"rate_multiplier"`
	Capacity       *ModelMarketplaceCapacity `json:"capacity,omitempty"`
	ModelCount     int                       `json:"model_count"`
	Models         []ModelMarketplaceModel   `json:"models"`
}

func ModelMarketplaceGroupsFromService(groups []service.ModelMarketplaceGroup) []ModelMarketplaceGroup {
	out := make([]ModelMarketplaceGroup, 0, len(groups))
	for _, group := range groups {
		models := make([]ModelMarketplaceModel, 0, len(group.Models))
		for _, model := range group.Models {
			models = append(models, ModelMarketplaceModel{
				ID:          model.ID,
				DisplayName: model.DisplayName,
				Pricing:     modelMarketplacePricingFromService(model.Pricing),
			})
		}

		out = append(out, ModelMarketplaceGroup{
			ID:             group.ID,
			Name:           group.Name,
			Description:    group.Description,
			Platform:       group.Platform,
			DisplayBrand:   group.DisplayBrand,
			SortOrder:      group.SortOrder,
			RateMultiplier: group.RateMultiplier,
			Capacity:       modelMarketplaceCapacityFromService(group.Capacity),
			ModelCount:     group.ModelCount,
			Models:         models,
		})
	}
	return out
}

func modelMarketplaceCapacityFromService(capacity *service.GroupCapacitySummary) *ModelMarketplaceCapacity {
	if capacity == nil {
		return nil
	}
	return &ModelMarketplaceCapacity{
		ConcurrencyUsed: capacity.ConcurrencyUsed,
		ConcurrencyMax:  capacity.ConcurrencyMax,
		SessionsUsed:    capacity.SessionsUsed,
		SessionsMax:     capacity.SessionsMax,
		RPMUsed:         capacity.RPMUsed,
		RPMMax:          capacity.RPMMax,
	}
}

func modelMarketplacePricingFromService(pricing service.ModelMarketplacePricing) ModelMarketplacePricing {
	intervals := make([]ModelMarketplacePricingInterval, 0, len(pricing.ContextIntervals))
	for _, interval := range pricing.ContextIntervals {
		intervals = append(intervals, ModelMarketplacePricingInterval{
			MinTokens:                interval.MinTokens,
			MaxTokens:                interval.MaxTokens,
			TierLabel:                interval.TierLabel,
			InputPricePerToken:       interval.InputPricePerToken,
			OutputPricePerToken:      interval.OutputPricePerToken,
			CacheWritePricePerToken:  interval.CacheWritePricePerToken,
			CacheReadPricePerToken:   interval.CacheReadPricePerToken,
			ImageOutputPricePerToken: interval.ImageOutputPricePerToken,
			PerRequestPrice:          interval.PerRequestPrice,
		})
	}
	return ModelMarketplacePricing{
		PricingMode:              pricing.PricingMode,
		PriceStatus:              pricing.PriceStatus,
		InputPricePerToken:       pricing.InputPricePerToken,
		OutputPricePerToken:      pricing.OutputPricePerToken,
		CacheWritePricePerToken:  pricing.CacheWritePricePerToken,
		CacheReadPricePerToken:   pricing.CacheReadPricePerToken,
		ImageOutputPricePerToken: pricing.ImageOutputPricePerToken,
		PerRequestPrice:          pricing.PerRequestPrice,
		ContextIntervals:         intervals,
		ImagePrice1K:             pricing.ImagePrice1K,
		ImagePrice2K:             pricing.ImagePrice2K,
		ImagePrice4K:             pricing.ImagePrice4K,
	}
}
