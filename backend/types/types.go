package types

import "errors"

var (
	ErrInvalidEmailFormat = errors.New("invalid_email_format")
	ErrReadConfigFile     = errors.New("invalid_config_path")
)
