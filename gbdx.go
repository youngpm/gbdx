package gbdx

import (
	"time"
)

// GBDX_HTTP_TIMEOUT is the default timeout used by the context in http requests.
const GBDX_HTTP_TIMEOUT = 60 * time.Second

// endpoints holds where the various gbdx endpoints live.
var endpoints = struct {
	tokens          string
	browse          string
	browseJSON      string
	browseMetadata  string
	thumbnail       string
	orders          string
	ordersHeartbeat string
}{
	tokens:          "https://geobigdata.io/auth/v1/oauth/token/",
	browse:          "https://geobigdata.io/thumbnails/v1/browse/",
	browseJSON:      "https://geobigdata.io/thumbnails/v1/get/",
	browseMetadata:  "https://geobigdata.io/thumbnails/v1/metadata/",
	thumbnail:       "https://geobigdata.io/thumbnails/v1/thumbnail/",
	orders:          "https://geobigdata.io/orders/v2/order/",
	ordersHeartbeat: "https://geobigdata.io/orders/v2/heartbeat/",
}
