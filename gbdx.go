package gbdx

// Config holds the various configuation items we need to interact with GBDX.
type Config struct {
	Username     string `mapstructure:"gbdx_username" toml:"gbdx_username"`
	Password     string `mapstructure:"gbdx_password" toml:"gbdx_password"`
	ClientID     string `mapstructure:"gbdx_client_id" toml:"gbdx_client_id"`
	ClientSecret string `mapstructure:"gbdx_client_secret" toml:"gbdx_client_secret"`
}
