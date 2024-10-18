package controller

//user_usecaseで実装したユーザー登録、ログイン、ログアウト、ユーザー情報更新メソッドを呼び出す
//このプログラムの処理が一番外側に位置し、routerでAPIとして呼び出される
//routerで受け取ったボディーの入力値をusecaseで使える形に変換する

import (
	"backend/model"
	"backend/usecase"
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

func (uc *UserController) SignUp(c echo.Context) error { //クライアントから受け取るリクエストボディの値を構造体の値に変換
	user := model.User{}
	if err := c.Bind(&user); err != nil { //リクエストをユーザーオブジェクトが指し示す先の値に格納
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes) //Createdのステータス、新しく作成したユーザーを返す
}

func (uc *UserController) Login(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
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
