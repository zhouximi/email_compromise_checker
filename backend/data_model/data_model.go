package data_model

import "time"

/***************res&resp****************************/
type EmailCheckAPIRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailCheckAPIResponse struct {
	Email       string `json:"email" binding:"required,email"`
	Conpromised bool   `json:"conpromised" binding:"required,conpromised"`
}

type CompromiseResult struct {
	Compromised bool
	// CompromiseLogs []*CompromiseLog  // todo: future feature
}

type CompromiseLog struct {
	CompromiseDate time.Time
	Source         string
}

type EmailInfo struct {
	Email       string `json:"email"`       // user email address
	Compromised bool   `json:"compromised"` // whether the email is compromised
}

type EmailRecord struct {
	Email string `json:"email"` // user email address
}

/**************middleware config*********************/
type LocalCacheConfig struct {
	NumCounters int64 `json:"NumCounters"`
	MaxCost     int64 `json:"MaxCost"`
	BufferItems int64 `json:"BufferItems"`
}

type MySQLConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
