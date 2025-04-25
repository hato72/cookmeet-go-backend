package usecase

// 全ての料理履歴を取得するGetAllCuisines、指定したIDに一致する料理を取得するGetCuisineByID、
// 料理を削除するDeleteCuisine、料理を追加するAddCuisine、料理を更新するSetCuisineを実装している
// それぞれcuisine_repositoryのメソッドを呼び出している

import (
	"backend/model"
	"backend/repository"
	"backend/validator"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ICuisineUsecase interface {
	GetAllCuisines(UserID uint) ([]model.CuisineResponse, error)
	GetCuisineByID(UserID uint, cuisineID uint) (model.CuisineResponse, error)
	// CreateCuisine(cuisine model.Cuisine) (model.CuisineResponse, error)
	// UpdateCuisine(cuisine model.Cuisine, UserID uint, cuisineID uint) (model.CuisineResponse, error)
	DeleteCuisine(UserID uint, cuisineID uint) error
	AddCuisine(cuisine model.Cuisine, iconFile *string, url string, title string) (model.CuisineResponse, error)
	SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, UserID uint, cuisineID uint) (model.CuisineResponse, error)
}

type cuisineUsecase struct {
	cr repository.ICuisineRepository
	cv validator.ICuisineValidator
}

func NewCuisineUsecase(tr repository.ICuisineRepository, tv validator.ICuisineValidator) ICuisineUsecase { // コンストラクタ
	return &cuisineUsecase{tr, tv}
}

func (cu *cuisineUsecase) GetAllCuisines(UserID uint) ([]model.CuisineResponse, error) {
	cuisines := []model.Cuisine{}
	if err := cu.cr.GetAllCuisines(&cuisines, UserID); err != nil {
		return nil, err
	}
	resCuisines := []model.CuisineResponse{}
	for _, v := range cuisines {
		t := model.CuisineResponse{
			ID:        v.ID,
			Title:     v.Title,
			IconURL:   v.IconURL,
			URL:       v.URL,
			Comment:   v.Comment,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			UserID:    v.UserID,
		}
		resCuisines = append(resCuisines, t)
	}
	return resCuisines, nil
}

func (cu *cuisineUsecase) GetCuisineByID(UserID uint, cuisineID uint) (model.CuisineResponse, error) {
	cuisine := model.Cuisine{}
	if err := cu.cr.GetCuisineByID(&cuisine, UserID, cuisineID); err != nil {
		return model.CuisineResponse{}, err
	}
	rescuisine := model.CuisineResponse{
		ID:        cuisine.ID,
		Title:     cuisine.Title,
		IconURL:   cuisine.IconURL,
		URL:       cuisine.URL,
		Comment:   cuisine.Comment,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		UserID:    cuisine.UserID,
	}
	return rescuisine, nil
}

// func (cu *cuisineUsecase) CreateCuisine(cuisine model.Cuisine) (model.CuisineResponse, error) {
// 	if err := cu.cv.CuisineValidate(cuisine); err != nil {
// 		return model.CuisineResponse{}, err
// 	}
// 	if err := cu.cr.CreateCuisine(&cuisine); err != nil {
// 		return model.CuisineResponse{}, err
// 	}
// 	rescuisine := model.CuisineResponse{
// 		ID:        cuisine.ID,
// 		Title:     cuisine.Title,
// 		IconURL:   cuisine.IconURL,
// 		URL:       cuisine.URL,
// 		CreatedAt: cuisine.CreatedAt,
// 		UpdatedAt: cuisine.UpdatedAt,
// 		UserID:    cuisine.UserID,
// 	}
// 	//log.Print(rescuisine)
// 	return rescuisine, nil
// }

// func (cu *cuisineUsecase) UpdateCuisine(cuisine model.Cuisine, UserID uint, cuisineID uint) (model.CuisineResponse, error) {
// 	if err := cu.cr.UpdateCuisine(&cuisine, UserID, cuisineID); err != nil {
// 		return model.CuisineResponse{}, err
// 	}
// 	// if err := cu.cr.AddURL(&cuisine, UserID, cuisineID); err != nil {
// 	// 	return model.CuisineResponse{}, err
// 	// }
// 	rescuisine := model.CuisineResponse{
// 		ID:        cuisine.ID,
// 		Title:     cuisine.Title,
// 		IconURL:   cuisine.IconURL,
// 		URL:       cuisine.URL,
// 		CreatedAt: cuisine.CreatedAt,
// 		UpdatedAt: cuisine.UpdatedAt,
// 		UserID:    cuisine.UserID,
// 	}
// 	return rescuisine, nil
// }

func (cu *cuisineUsecase) DeleteCuisine(UserID uint, cuisineID uint) error {
	if err := cu.cr.DeleteCuisine(UserID, cuisineID); err != nil {
		return err
	}
	return nil
}

func (cu *cuisineUsecase) AddCuisine(cuisine model.Cuisine, iconFile *string, url string, title string) (model.CuisineResponse, error) {
	if iconFile != nil {
		cuisine.IconURL = iconFile
	}

	if url != "" {
		cuisine.URL = url
	}

	if title != "" {
		cuisine.Title = title
	}

	if err := cu.cv.CuisineValidate(cuisine); err != nil {
		return model.CuisineResponse{}, err
	}
	if err := cu.cr.CreateCuisine(&cuisine); err != nil {
		return model.CuisineResponse{}, err
	}
	rescuisine := model.CuisineResponse{
		ID:        cuisine.ID,
		Title:     cuisine.Title,
		IconURL:   cuisine.IconURL,
		URL:       cuisine.URL,
		Comment:   cuisine.Comment, // コメントを追加
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		UserID:    cuisine.UserID,
	}
	// log.Print(rescuisine)
	return rescuisine, nil
}

func (cu *cuisineUsecase) SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, UserID uint, cuisineID uint) (model.CuisineResponse, error) {
	// cuisine := model.Cuisine{}

	if iconFile != nil {
		src, err := iconFile.Open()
		if err != nil {
			return model.CuisineResponse{}, err
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			return model.CuisineResponse{}, err
		}

		hasher := sha256.New()
		hasher.Write(data)
		hashValue := hex.EncodeToString(hasher.Sum(nil))

		ext := filepath.Ext(iconFile.Filename)

		img_url := "cuisine_icons/" + hashValue + ext

		dst, err := os.Create("./cuisine_images/" + img_url)
		if err != nil {
			return model.CuisineResponse{}, err
		}

		defer dst.Close()

		if _, err := dst.Write(data); err != nil {
			return model.CuisineResponse{}, nil
		}

		cuisine.IconURL = &img_url
	}

	if url != "" {
		cuisine.URL = url
	}

	if title != "" {
		cuisine.Title = title
	}

	updatedCuisine := model.Cuisine{
		ID:        cuisine.ID,
		Title:     title,
		IconURL:   cuisine.IconURL,
		URL:       url,
		Comment:   cuisine.Comment,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		User:      cuisine.User,
		UserID:    cuisine.UserID,
	}
	// log.Print("cuisine", cuisine)
	// log.Print("updatedCuisine", updatedCuisine)

	if err := cu.cr.SettingCuisine(&updatedCuisine); err != nil {
		return model.CuisineResponse{}, err
	}

	rescuisine := model.CuisineResponse{
		ID:        updatedCuisine.ID,
		Title:     cuisine.Title,
		IconURL:   cuisine.IconURL,
		URL:       cuisine.URL,
		Comment:   cuisine.Comment,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: updatedCuisine.UpdatedAt,
		UserID:    updatedCuisine.UserID,
	}

	// log.Print("updatedCuisine")
	// log.Print("title", updatedCuisine.Title)
	// log.Print("url", updatedCuisine.URL)
	// log.Print("CreatedAt", updatedCuisine.CreatedAt)
	// log.Print("UpdatedAt", updatedCuisine.UpdatedAt)

	// log.Print("cuisine")
	// log.Print("title", cuisine.Title)
	// log.Print("url", cuisine.URL)
	// log.Print("CreatedAt", cuisine.CreatedAt)
	// log.Print("UpdatedAt", cuisine.UpdatedAt)

	// log.Print("rescuisine")
	// log.Print("title", rescuisine.Title)
	// log.Print("url", rescuisine.URL)
	// log.Print("CreatedAt", rescuisine.CreatedAt)
	// log.Print("UpdatedAt", rescuisine.UpdatedAt)

	return rescuisine, nil
}
