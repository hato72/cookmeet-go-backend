package router

import (
	"backend/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, cc controller.ICuisineController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{ // corsのミドルウェア
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")}, // デプロイしたときに取得できるドメイン
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, // 許可するヘッダーの一覧
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"}, // 許可したいメソッド
		AllowCredentials: true,                                     // クッキーの送受信を可能にする
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{ // csrfのミドルウェア
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		// CookieSameSite: http.SameSiteDefaultMode, // postmanで確認のため
	}))
	// e.Start(":8080")
	// e.GET("/", func(c echo.Context) error {
	// 	return c.JSON(http.StatusOK, "hello world")
	// })

	e.POST("/signup", uc.SignUp) // エンドポイント追加
	e.POST("/login", uc.Login)
	e.POST("/logout", uc.Logout)
	// e.PUT("/update", uc.Update)
	// e.PUT("/update", uc.Update, echojwt.WithConfig(echojwt.Config{
	// 	SigningKey:  []byte(os.Getenv("SECRET")),
	// 	TokenLookup: "cookie:token",
	// }))
	e.GET("/csrf", uc.CsrfToken)

	u := e.Group("/update")
	u.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	u.PUT("", uc.Update)

	c := e.Group("/cuisines")
	c.Use(echojwt.WithConfig(echojwt.Config{ // エンドポイントにミドルウェアを追加
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	c.GET("", cc.GetAllCuisines)            // cuisinesのエンドポイントにリクエストがあった場合
	c.GET("/:cuisineID", cc.GetCuisineByID) // リクエストパラメーターにcuisineIDが入力された場合
	// c.POST("", cc.CreateCuisine)
	c.POST("", cc.AddCuisine) // cuisineテーブル追加
	// c.PUT("/:cuisineID", cc.UpdateCuisine) // titleしか更新されない
	c.DELETE("/:cuisineID", cc.DeleteCuisine)

	c.PUT("/:cuisineID", cc.SetCuisine) // cuisineの更新
	// c.PUT("/url/:cuisineID", cc.AddURL)
	return e
}
