package configs

const (
	ProjectName = "meeseeks"

	RequestHeaderJWTKey = "Authorization"

	RequestHeaderRefershKey = "Authenticate"

	RequestHeaderJWTExpireTime = 60 * 60 * 2

	RequestOpenAPIHeaderAppID     = "Authorization-Appid"
	RequestOpenAPIHeaderSign      = "Authorization-Sign"
	RequestOpenAPIHeaderTimestamp = "Authorization-Timestamp"
	RequestOpenAPIHeaderNone      = "Authorization-None"
)
