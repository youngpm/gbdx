package gbdx

import (
	"time"
)

const GBDX_HTTP_TIMEOUT = 60 * time.Second

var endpoints = struct {
	tokens string
	browse string
}{
	tokens: "https://geobigdata.io/auth/v1/oauth/token/",
	browse: "https://geobigdata.io/thumbnails/v1/browse/",
}
