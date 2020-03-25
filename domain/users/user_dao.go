package users

import (
	"fmt"
	"github.com/rezwanul-haque/Metadata-Service/datasources/mysql/lms_db"
	"github.com/rezwanul-haque/Metadata-Service/logger"
	"github.com/rezwanul-haque/Metadata-Service/utils/errors"
	"strings"
)

const (
	queryInsertUser                  = "INSERT INTO user(domain, user_id, metadata, status) VALUES(LOWER(?), ?, ?, ?);"
	queryGetUser                     = "SELECT id, domain, user_id, metadata, status FROM user WHERE user_id=?;"
	queryUpdateUser                  = "UPDATE user SET metadata=JSON_MERGE(metadata, ?), status=? WHERE user_id=?;"
	queryFindUsersByDomain           = "SELECT id, domain, user_id, metadata, status FROM user WHERE domain=LOWER(?);"
	queryFindUsersByDomainAndUserIds = "SELECT id, domain, user_id, metadata, status FROM user WHERE domain=LOWER(?) AND user_id=?;"
)

func (user *User) Save() *errors.RestErr {
	stmt, err := lms_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.Domain, user.UserId, user.Metadata, user.Status)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Get() *errors.RestErr {
	stmt, err := lms_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.UserId)
	if getErr := result.Scan(&user.Id, &user.Domain, &user.UserId, &user.Metadata, &user.Status); getErr != nil {
		logger.Error("error when trying to get user", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := lms_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Metadata, user.Status, user.UserId)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Search(domain string) ([]User, *errors.RestErr) {
	stmt, err := lms_db.Client.Prepare(queryFindUsersByDomain)
	if err != nil {
		logger.Error("error when trying to prepare search user by domain statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(domain)
	if err != nil {
		logger.Error("error when trying to search users by domain", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Domain, &user.UserId, &user.Metadata, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching domain: %s", domain))
	}
	return results, nil
}

func (user *User) SearchByDomainAndIds(domain string, userIds []string) (*map[string]User, *errors.RestErr) {
	var results map[string]User

	results = make(map[string]User)

	for _, userId := range userIds {
		data, findErr := findUserById(domain, userId)
		if findErr != nil {
			return nil, findErr
		}
		results[userId] = *data
	}
	return &results, nil
}

func findUserById(domain string, userId string) (*User, *errors.RestErr) {
	var result User
	stmt, err := lms_db.Client.Prepare(queryFindUsersByDomainAndUserIds)
	if err != nil {
		logger.Error("error when trying to prepare search user by domain and userId statement", err)
		return nil, errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	row := stmt.QueryRow(domain, userId)
	if getErr := row.Scan(&result.Id, &result.Domain, &result.UserId, &result.Metadata, &result.Status); getErr != nil {
		if strings.Contains(getErr.Error(), "sql: no rows in result set") {
			return &result, nil
		}
		logger.Error("error when trying to find user", getErr)
		return nil, errors.NewInternalServerError("database error")
	}
	return &result, nil
}
