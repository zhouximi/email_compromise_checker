package data_model

type EmailCheckAPIRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailInfo struct {
	Email       string `json:"email"`       // user email address
	Compromised bool   `json:"compromised"` // whether the email is compromised
}

type LocalCacheConfig struct {
	NumCounters int64 `json:"NumCounters"`
	MaxCost     int64 `json:"MaxCost"`
	BufferItems int64 `json:"BufferItems"`
}
