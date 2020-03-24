package users

import (
	"github.com/gin-gonic/gin"
	"github.com/rezwanul-haque/Metadata-Service/domain/users"
	"github.com/rezwanul-haque/Metadata-Service/services"
	"github.com/rezwanul-haque/Metadata-Service/utils/errors"
	"github.com/rezwanul-haque/Metadata-Service/utils/helpers"
	"net/http"
	"strconv"
)

func getRlsReferrer(rlsReferrer string) (*string, *errors.RestErr) {
	if rlsReferrer == "" {
		return nil, errors.NewBadRequestError("RLS-Referrer header is not present")
	}
	return &rlsReferrer, nil
}

func getPageOrSize(pageOrSizeParam string) (int, *errors.RestErr) {
	param, err := strconv.Atoi(pageOrSizeParam)
	if err != nil {
		return 0, errors.NewBadRequestError("page or size should be a number")
	}
	return param, nil
}

func Create(c *gin.Context) {
	RlsReferrer, headerErr := getRlsReferrer(c.GetHeader("RLS-Referrer"))
	if headerErr != nil {
		c.JSON(headerErr.Status, headerErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Domain = *RlsReferrer

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetAllUsers(c *gin.Context) {
	RlsReferrer, headerErr := getRlsReferrer(c.GetHeader("RLS-Referrer"))
	if headerErr != nil {
		c.JSON(headerErr.Status, headerErr)
		return
	}
	page, pageErr := getPageOrSize(c.DefaultQuery("page", "1"))
	if pageErr != nil {
		c.JSON(pageErr.Status, pageErr)
		return
	}
	size, sizeErr := getPageOrSize(c.DefaultQuery("size", "100"))
	if sizeErr != nil {
		c.JSON(sizeErr.Status, sizeErr)
		return
	}

	users, getErr := services.UsersService.SearchUser(*RlsReferrer)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	start, end := helpers.Paginate(page, size, len(users))

	paginatedUsers := users[start:end]

	c.JSON(http.StatusOK, paginatedUsers)
}

func Update(c *gin.Context) {
	RlsReferrer, headerErr := getRlsReferrer(c.GetHeader("RLS-Referrer"))
	userId := c.Param("user_id")
	if headerErr != nil {
		c.JSON(headerErr.Status, headerErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Domain = *RlsReferrer
	user.UserId = userId

	result, updateErr := services.UsersService.UpdateUser(user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
	}
	c.JSON(http.StatusOK, result)
}
