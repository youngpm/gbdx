package gbdx

import (
	"time"
)

const GBDX_HTTP_TIMEOUT = 60 * time.Second

var endpoints = struct {
	tokens string
}{
	tokens: "https://geobigdata.io/auth/v1/oauth/token/",
}
