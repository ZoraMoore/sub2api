//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModelMarketplaceService_ListPublicFiltersPrivateGroupsAndKeepsPricing(t *testing.T) {
	inputPrice := 0.000001
	outputPrice := 0.000002
	channels := []Channel{{
		ID:          10,
		Name:        "OpenAI Channel",
		Description: "public openai channel",
		Status:      StatusActive,
		GroupIDs:    []int64{1, 2, 3},
		ModelPricing: []ChannelModelPricing{{
			ID:          100,
			Platform:    PlatformOpenAI,
			Models:      []string{"gpt-5.5"},
			BillingMode: BillingModeToken,
			InputPrice:  &inputPrice,
			OutputPrice: &outputPrice,
		}},
	}}
	groupRepo := &stubGroupRepoForAvailable{activeGroups: []Group{
		{
			ID:                 1,
			Name:               "Public OpenAI",
			Description:        "visible group",
			Platform:           PlatformOpenAI,
			RateMultiplier:     0.5,
			Status:             StatusActive,
			ActiveAccountCount: 2,
		},
		{
			ID:                 2,
			Name:               "Exclusive OpenAI",
			Platform:           PlatformOpenAI,
			RateMultiplier:     0.7,
			IsExclusive:        true,
			Status:             StatusActive,
			ActiveAccountCount: 2,
		},
		{
			ID:                 3,
			Name:               "Empty OpenAI",
			Platform:           PlatformOpenAI,
			RateMultiplier:     0.9,
			Status:             StatusActive,
			ActiveAccountCount: 0,
		},
	}}
	channelSvc := newAvailableChannelService(channels, groupRepo)
	svc := NewModelMarketplaceService(channelSvc, groupRepo, nil)

	groups, err := svc.ListPublic(context.Background())
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, int64(1), groups[0].ID)
	require.Equal(t, "Public OpenAI", groups[0].Name)
	require.Equal(t, "visible group", groups[0].Description)
	require.Equal(t, PlatformOpenAI, groups[0].Platform)
	require.Equal(t, 0.5, groups[0].RateMultiplier)
	require.Equal(t, 1, groups[0].ModelCount)
	require.Len(t, groups[0].Models, 1)
	require.Equal(t, "gpt-5.5", groups[0].Models[0].ID)
	require.Equal(t, "gpt-5.5", groups[0].Models[0].DisplayName)
	require.Equal(t, string(BillingModeToken), groups[0].Models[0].Pricing.PricingMode)
	require.Equal(t, "priced", groups[0].Models[0].Pricing.PriceStatus)
	require.Equal(t, inputPrice, groups[0].Models[0].Pricing.InputPricePerToken)
	require.Equal(t, outputPrice, groups[0].Models[0].Pricing.OutputPricePerToken)
}
