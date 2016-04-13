package gbdx

import "time"

const GBDX_HTTP_TIMEOUT = 60 * time.Second

type GBDXConfiguration struct {
	Username       string
	Password       string
	ClientID       string
	ClientPassword string
}
