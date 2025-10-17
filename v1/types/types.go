package types

type ConfigFile struct {
	ServerBaseUrl string `json:"server_base_url"`
	ServerPort string `json:"server_port"`
	ServerAPIKey string `json:"server_api_key"`
	ServerCookieSecret string `json:"server_cookie_secret"`
	ServerCookieAdminSecretMessage string `json:"server_cookie_admin_secret_message"`
	ServerCookieSecretMessage string `json:"server_cookie_secret_message"`
	ServerLiveUrl string `json:"server_live_url"`
	LocalHostUrl string `json:"local_host_url"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
	TimeZone string `json:"time_zone"`
	BoltDBPath string `json:"bolt_db_path"`
	BoltDBEncryptionKey string `json:"bolt_db_encryption_key"`
	BoltDBEncryptionKeyClients string `json:"bolt_db_encryption_key_clients"`
}