package constants

const (
	KeyCookieRefreshToken = "rt"
	KeyCookieAccessToken  = "at"
)

const (
	AccessExpiredAt  = 1 * 10            // 15 minutes in seconds
	RefreshExpiredAt = 30 * 24 * 60 * 60 // 30 days in seconds
	VerifyExpiredAt  = 10 * 60           // 10 minutes in seconds
	ForgotExpiredAt  = 15 * 60           // 15 minutes in seconds
)
