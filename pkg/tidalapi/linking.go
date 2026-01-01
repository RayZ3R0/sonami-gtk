package tidalapi

import (
	"context"

	"codeberg.org/dergs/tidalwave/pkg/tidalapi/auth"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
)

func StartDeviceLinking(cb func(*auth.DeviceLinkingChallenge, context.CancelFunc)) (*auth.TokenResponse, error) {
	linking, err := auth.RequestDeviceLinkingChallenge(internal.ClientID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	cb(linking, cancel)
	return auth.PollDeviceLinkingStatus(ctx, internal.ClientID, internal.ClientSecret, linking.DeviceCode, linking.Interval)
}
