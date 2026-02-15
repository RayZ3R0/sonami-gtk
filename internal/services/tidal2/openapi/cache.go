package openapi

import "codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"

var (
	artistCache = make(map[string]*openapi.Artist)
)
