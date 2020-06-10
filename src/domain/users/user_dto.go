package users

import (
	"encoding/json"
	"github.com/rezwanul-haque/Metadata-Service/src/utils/errors"
	"strings"
)

type User struct {
	Id       int64           `json:"id"`
	Domain   string          `json:"domain"`
	UserId   string          `json:"user_id"`
	Metadata json.RawMessage `json:"metadata"`
	Status   string          `json:"status"`
}

type Users []User

func (user *User) Validate() *errors.RestErr {
	user.Domain = strings.TrimSpace(strings.ToLower(user.Domain))
	user.UserId = strings.TrimSpace(strings.ToLower(user.UserId))
	if user.Domain == "" {
		return errors.NewBadRequestError("RLS-Referrer header is not present")
	}
	if user.UserId == "" {
		return errors.NewBadRequestError("user_id is missing or invalid")
	}
	return nil
}

type QueryParamRequest struct {
	UserIds string  `form:"user_ids,omitempty"`
	Status  string  `form:"status,omitempty"`
	Page    float64 `form:"page,omitempty"`
	Size    float64 `form:"size,omitempty"`
}

type UserSearchResponse struct {
	Domain   string          `json:"domain"`
	UserId   string          `json:"user_id"`
	Metadata json.RawMessage `json:"metadata"`
	Status   string          `json:"status"`
}

type UsersSearchResponse []UserSearchResponse
