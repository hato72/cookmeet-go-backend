package controller

// GetAllCuisines: cuisine_usecaseの同メソッドを呼び出している
// GetCuisineByID:cuisine_usecaseの同メソッドを呼び出している
// DeleteCuisine:料理を削除している
// AddCuisine:cuisine_usecaseの同メソッドを呼び出している
// SetCuisine:cuisine_usecaseのgetAllcuisinesメソッドで料理を取得したのち、同メソッドを呼び出している
// このプログラムが一番外側であり、routerで呼び出される

import (
	"backend/model"
	"backend/usecase"
	"backend/utils"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ICuisineController interface {
	GetAllCuisines(c echo.Context) error
	GetCuisineByID(c echo.Context) error
	// CreateCuisine(c echo.Context) error
	// UpdateCuisine(c echo.Context) error
	DeleteCuisine(c echo.Context) error
	AddCuisine(c echo.Context) error
	// SetCuisine(c echo.Context) error
}

type cuisineController struct {
	cu usecase.ICuisineUsecase
}

func NewCuisineController(cu usecase.ICuisineUsecase) ICuisineController {
	return &cuisineController{cu}
}

func (cc *cuisineController) GetAllCuisines(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)    // コンテキストからjwtをデコードした値を読み込む
	claims := user.Claims.(jwt.MapClaims) // その中のデコードされたclaimsを取得
	UserID := claims["user_id"]           // claimsの中のUserIDを取得
	// log.Print(UserID)

	cuisineRes, err := cc.cu.GetAllCuisines(uint(UserID.(float64))) // 一度floatにしてからuintに型変換
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cuisineRes)
}

func (cc *cuisineController) GetCuisineByID(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["user_id"]

	id := c.Param("cuisineID")
	cuisineID, err := strconv.ParseUint(id, 10, 32) // Atoiの代わりにParseUintを使用
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid cuisine ID")
	}
	cuisineRes, err := cc.cu.GetCuisineByID(uint(UserID.(float64)), uint(cuisineID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cuisineRes)
}

func (cc *cuisineController) DeleteCuisine(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64)) // userIDに修正

	id := c.Param("cuisineID")
	cuisineID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid cuisine ID")
	}

	// 不要なBindを削除（DeleteCuisineではリクエストボディは不要）

	if delErr := cc.cu.DeleteCuisine(userID, uint(cuisineID)); delErr != nil {
		switch {
		case errors.Is(delErr, usecase.ErrCuisineNotFound):
			return c.JSON(http.StatusNotFound, delErr.Error())
		case errors.Is(delErr, usecase.ErrUnauthorized):
			return c.JSON(http.StatusForbidden, delErr.Error())
		default:
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (cc *cuisineController) AddCuisine(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	UserID := claims["user_id"]

	iconFile, err := c.FormFile("icon")
	title := c.FormValue("title")
	url := c.FormValue("url")
	comment := c.FormValue("comment") // コメントを取得

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	var imageURL string
	if iconFile != nil {
		// ファイルを読み込みbase64エンコード
		src, openErr := iconFile.Open()
		if openErr != nil {
			return c.JSON(http.StatusInternalServerError, openErr.Error())
		}
		defer src.Close()

		UserIDStr := strconv.FormatUint(uint64(UserID.(float64)), 10)

		// Cloud Storage にアップロード
		bucket := "cookmeet"
		objectName := "images/" + UserIDStr + "/" + uuid.New().String() + filepath.Ext(iconFile.Filename)

		imageURL, err = utils.UploadToCloudStorage(bucket, objectName, src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	cuisine := model.Cuisine{}
	cuisine.UserID = uint(UserID.(float64))
	cuisine.Title = title
	cuisine.URL = url
	cuisine.Comment = comment // コメントをセット
	// 画像がアップロードされた場合のみURLをセット
	if imageURL != "" {
		cuisine.IconURL = &imageURL // Cloud StorageのURLをセット
	}

	cuisineRes, err := cc.cu.AddCuisine(cuisine, &imageURL, url, title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cuisineRes)
}

// func (cc *cuisineController) SetCuisine(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	UserID := claims["user_id"]
// 	id := c.Param("cuisineID")
// 	cuisineID, err := strconv.ParseUint(id, 10, 32)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "Invalid cuisine ID")
// 	}

// 	url := c.FormValue("url")
// 	iconFile, err := c.FormFile("icon")
// 	title := c.FormValue("title")
// 	if err != nil {
// 		if err != http.ErrMissingFile {
// 			return c.JSON(http.StatusBadRequest, err.Error())
// 		}
// 	}

// 	cuisine := model.Cuisine{}
// 	cuisine.ID = uint(cuisineID)
// 	cuisine.UserID = uint(UserID.(float64))
// 	// cuisine.URL = url
// 	// cuisine.URL = url

// 	if bindErr := c.Bind(&cuisine); bindErr != nil {
// 		return c.JSON(http.StatusBadRequest, bindErr.Error())
// 	}

// 	cuisineRes, err := cc.cu.GetCuisineByID(uint(UserID.(float64)), uint(cuisineID))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	cuisine.CreatedAt = cuisineRes.CreatedAt
// 	cuisine.UpdatedAt = cuisineRes.UpdatedAt
// 	cuisine.Title = cuisineRes.Title
// 	cuisine.IconURL = cuisineRes.IconURL
// 	cuisine.URL = cuisineRes.URL

// 	newcuisineRes, err := cc.cu.SetCuisine(cuisine, iconFile, url, title, uint(UserID.(float64)), uint(cuisineID))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, newcuisineRes)
// }
