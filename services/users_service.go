package services

import (
	"github.com/rezwanul-haque/Metadata-Service/domain/users"
	"github.com/rezwanul-haque/Metadata-Service/utils/errors"
	"github.com/rezwanul-haque/Metadata-Service/utils/helpers"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(userId string) (*users.User, *errors.RestErr)
	UpdateUser(user users.User) (*users.User, *errors.RestErr)
	SearchUser(domain string) (users.Users, *errors.RestErr)
	SearchUserByDomainAndIds(domain string, userIds []string) (*map[string]users.User, *errors.RestErr)
	FindByQueryDSL(queryDSL interface{}, domain string, queryParams users.QueryParamRequest) (users.UsersSearchResponse, *errors.RestErr)
}

func (u *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *usersService) GetUser(userId string) (*users.User, *errors.RestErr) {
	result := &users.User{UserId: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersService) UpdateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	current, getErr := UsersService.GetUser(user.UserId)
	if getErr != nil {
		return nil, getErr
	}
	if !helpers.ByteEmpty(user.Metadata) {
		current.Metadata = user.Metadata
	}

	if user.Status != "" {
		current.Status = user.Status
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (u *usersService) SearchUser(domain string) (users.Users, *errors.RestErr) {
	user := &users.User{}
	return user.Search(domain)
}

func (u *usersService) SearchUserByDomainAndIds(domain string, userIds []string) (*map[string]users.User, *errors.RestErr) {
	user := &users.User{}
	return user.SearchByDomainAndIds(domain, userIds)
}

func (u *usersService) FindByQueryDSL(queryDSL interface{}, domain string, queryParams users.QueryParamRequest) (users.UsersSearchResponse, *errors.RestErr) {
	users := users.UserSearchResponse{
		Domain: domain,
	}
	return users.SearchByMetadata(queryDSL, queryParams)
}
