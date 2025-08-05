package handler

import (
	"github.com/zhouximi/email_compromise_checker/data_model"
	"github.com/zhouximi/email_compromise_checker/middleware/cache"
	"github.com/zhouximi/email_compromise_checker/middleware/db"
	"reflect"
	"testing"
)

func TestEmailCheckHandler_IsEmailCompromised(t *testing.T) {
	type fields struct {
		cache cache.ICache
		db    db.IDB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *data_model.EmailInfo
		wantErr bool
	}{
		{
			name: "invalid email format",
			fields: fields{
				cache: &cache.MockCache{},
				db:    &db.MockDB{},
			},
			args: args{
				email: "invalid-email",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "email not compromised, not in cache",
			fields: fields{
				cache: &cache.MockCache{},
				db:    &db.MockDB{CompromisedEmails: map[string]bool{}},
			},
			args: args{
				email: "safe@example.com",
			},
			want: &data_model.EmailInfo{
				Email:       "safe@example.com",
				Compromised: false,
			},
			wantErr: false,
		},
		{
			name: "email compromised, not in cache",
			fields: fields{
				cache: &cache.MockCache{},
				db:    &db.MockDB{CompromisedEmails: map[string]bool{"bad@example.com": true}},
			},
			args: args{
				email: "bad@example.com",
			},
			want: &data_model.EmailInfo{
				Email:       "bad@example.com",
				Compromised: true,
			},
			wantErr: false,
		},
		{
			name: "email info found in cache",
			fields: fields{
				cache: &cache.MockCache{Items: map[string]interface{}{"email:cached@example.com": &data_model.EmailInfo{Email: "cached@example.com", Compromised: true}}},
				db:    &db.MockDB{},
			},
			args: args{
				email: "cached@example.com",
			},
			want: &data_model.EmailInfo{
				Email:       "cached@example.com",
				Compromised: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EmailCheckHandler{
				cache: tt.fields.cache,
				db:    tt.fields.db,
			}
			got, err := h.IsEmailCompromised(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsEmailCompromised() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsEmailCompromised() got = %v, want %v", got, tt.want)
			}
		})
	}
}
