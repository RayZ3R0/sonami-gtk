package tidalapi

import (
	"context"

	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/auth"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
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
