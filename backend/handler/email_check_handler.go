package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/zhouximi/email_compromise_checker/data_model"
	"github.com/zhouximi/email_compromise_checker/middleware/cache"
	"github.com/zhouximi/email_compromise_checker/middleware/db"
	"github.com/zhouximi/email_compromise_checker/types"
)

var GlobalHandler IEmailCheckHandler

type IEmailCheckHandler interface {
	IsEmailCompromised(email string) (*data_model.EmailInfo, error)
}

type EmailCheckHandler struct {
	cache cache.ICache
	db    db.IDB
}

func NewEmailCheckHandler(cache cache.ICache, db db.IDB) *EmailCheckHandler {
	return &EmailCheckHandler{
		cache: cache,
		db:    db,
	}
}

func (h *EmailCheckHandler) IsEmailCompromised(email string) (*data_model.EmailInfo, error) {
	emailInfo := &data_model.EmailInfo{
		Email: email,
	}

	if !isValidEmail(email) {
		return nil, types.ErrInvalidEmailFormat
	}

	if cachedEmailInfo := h.checkEmailInfoFromCache(email); cachedEmailInfo != nil {
		emailInfo.Compromised = cachedEmailInfo.Compromised
		return emailInfo, nil
	}

	dbEmailInfo, err := h.checkEmailInfoFromDB(email)
	if err != nil {
		return nil, err
	} else {
		emailInfo.Compromised = dbEmailInfo.Compromised
	}

	h.setEmailInfoToCache(emailInfo)

	return emailInfo, nil
}

func (h *EmailCheckHandler) checkEmailInfoFromCache(email string) *data_model.EmailInfo {
	cacheKey := fmt.Sprintf("email:%s", email)

	cachedValue, err := h.cache.Get(cacheKey)
	if err != nil || cachedValue == nil {
		return nil
	}

	emailInfo, ok := cachedValue.(*data_model.EmailInfo)
	if !ok {
		return nil
	}

	return emailInfo
}

func (h *EmailCheckHandler) checkEmailInfoFromDB(email string) (*data_model.EmailInfo, error) {
	result, err := h.db.RunQuery(email, func(db *sql.DB, key string) (interface{}, error) {
		query := "SELECT email FROM compromised_emails WHERE email = ?"

		var found string
		err := db.QueryRow(query, key).Scan(&found)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Email not found, not compromised
				return &data_model.EmailInfo{
					Email:       key,
					Compromised: false,
				}, nil
			}
			// Unexpected DB error
			return nil, err
		}

		// Email found, is compromised
		return &data_model.EmailInfo{
			Email:       found,
			Compromised: true,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	emailInfo, ok := result.(*data_model.EmailInfo)
	if !ok {
		return nil, errors.New("invalid type assertion to *EmailInfo")
	}

	return emailInfo, nil
}

func (h *EmailCheckHandler) setEmailInfoToCache(emailInfo *data_model.EmailInfo) {
	cacheKey := fmt.Sprintf("email:%s", emailInfo.Email)
	if err := h.cache.Set(cacheKey, emailInfo); err != nil {
		// Log the error but do not return it, as this is a non-critical operation
		log.Println("[saveEmailInfoToCache]failed to save email info to cache:", err)
	}
}
