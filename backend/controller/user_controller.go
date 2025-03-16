package controller

import (
	"backend/model"
	"backend/usecase"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Update(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type UserController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &UserController{uu}
}

func (uc *UserController) SignUp(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil { // リクエストボディをバインド
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		// エラー処理を統一した形式に修正
		// ユーザーが既に存在する場合は409 Conflictを返す
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes) // Createdのステータス、新規作成したユーザーを返す
}

func (uc *UserController) Login(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		// errors.Isでエラーをチェックするように変更
		if errors.Is(err, usecase.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, err.Error()) // ユーザーが見つからない場合
		} else if errors.Is(err, usecase.ErrInvalidPassword) {
			return c.JSON(http.StatusUnauthorized, err.Error()) // パスワードが間違っている場合
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie) //httpレスポンスに含める
	return c.NoContent(http.StatusOK)
}

func (uc *UserController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie) //httpレスポンスに含める
	return c.NoContent(http.StatusOK)
}

func (uc *UserController) Update(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//log.Print(user)

	authUser := c.Get("user").(*jwt.Token)
	claims := authUser.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	//log.Print(authUser)
	//log.Print(userId)

	userID := uint(userId.(float64))
	user.ID = userID
	newEmail := c.FormValue("email")
	newName := c.FormValue("name")
	newPassword := c.FormValue("password")
	iconFile, err := c.FormFile("icon")

	//log.Print(userID, newEmail, newName, newPassword, iconFile)

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	userRes, err := uc.uu.Update(user, newEmail, newName, newPassword, iconFile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, userRes)
}

func (uc *UserController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{ //クライアントにcsrfトークンをレスポンス
		"csrf_token": token,
	})
}
