package usecase

//全ての料理履歴を取得するGetAllCuisines、指定したIDに一致する料理を取得するGetCuisineById、
//料理を削除するDeleteCuisine、料理を追加するAddCuisine、料理を更新するSetCuisineを実装している
//それぞれcuisine_repositoryのメソッドを呼び出している

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
	GetAllCuisines(userId uint) ([]model.CuisineResponse, error)
	GetCuisineById(userId uint, cuisineId uint) (model.CuisineResponse, error)
	//CreateCuisine(cuisine model.Cuisine) (model.CuisineResponse, error)
	//UpdateCuisine(cuisine model.Cuisine, userId uint, cuisineId uint) (model.CuisineResponse, error)
	DeleteCuisine(userId uint, cuisineId uint) error
	AddCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string) (model.CuisineResponse, error)
	SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, userId uint, cuisineId uint) (model.CuisineResponse, error)
}

type cuisineUsecase struct {
	cr repository.ICuisineRepository
	cv validator.ICuisineValidator
}

func NewCuisineUsecase(tr repository.ICuisineRepository, tv validator.ICuisineValidator) ICuisineUsecase { //コンストラクタ
	return &cuisineUsecase{tr, tv}
}

func (cu *cuisineUsecase) GetAllCuisines(userId uint) ([]model.CuisineResponse, error) {
	cuisines := []model.Cuisine{}
	if err := cu.cr.GetAllCuisines(&cuisines, userId); err != nil {
		return nil, err
	}
	resCuisines := []model.CuisineResponse{}
	for _, v := range cuisines {
		t := model.CuisineResponse{
			ID:        v.ID,
			Title:     v.Title,
			IconUrl:   v.IconUrl,
			URL:       v.URL,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			UserId:    v.UserId,
		}
		resCuisines = append(resCuisines, t)
	}
	return resCuisines, nil
}

func (cu *cuisineUsecase) GetCuisineById(userId uint, cuisineId uint) (model.CuisineResponse, error) {
	cuisine := model.Cuisine{}
	if err := cu.cr.GetCuisineById(&cuisine, userId, cuisineId); err != nil {
		return model.CuisineResponse{}, err
	}
	rescuisine := model.CuisineResponse{
		ID:        cuisine.ID,
		Title:     cuisine.Title,
		IconUrl:   cuisine.IconUrl,
		URL:       cuisine.URL,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		UserId:    cuisine.UserId,
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
// 		IconUrl:   cuisine.IconUrl,
// 		URL:       cuisine.URL,
// 		CreatedAt: cuisine.CreatedAt,
// 		UpdatedAt: cuisine.UpdatedAt,
// 		UserId:    cuisine.UserId,
// 	}
// 	//log.Print(rescuisine)
// 	return rescuisine, nil
// }

// func (cu *cuisineUsecase) UpdateCuisine(cuisine model.Cuisine, userId uint, cuisineId uint) (model.CuisineResponse, error) {
// 	if err := cu.cr.UpdateCuisine(&cuisine, userId, cuisineId); err != nil {
// 		return model.CuisineResponse{}, err
// 	}
// 	// if err := cu.cr.AddURL(&cuisine, userId, cuisineId); err != nil {
// 	// 	return model.CuisineResponse{}, err
// 	// }
// 	rescuisine := model.CuisineResponse{
// 		ID:        cuisine.ID,
// 		Title:     cuisine.Title,
// 		IconUrl:   cuisine.IconUrl,
// 		URL:       cuisine.URL,
// 		CreatedAt: cuisine.CreatedAt,
// 		UpdatedAt: cuisine.UpdatedAt,
// 		UserId:    cuisine.UserId,
// 	}
// 	return rescuisine, nil
// }

func (cu *cuisineUsecase) DeleteCuisine(userId uint, cuisineId uint) error {
	if err := cu.cr.DeleteCuisine(userId, cuisineId); err != nil {
		return err
	}
	return nil
}

func (cu *cuisineUsecase) AddCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string) (model.CuisineResponse, error) {

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

		cuisine.IconUrl = &img_url
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
		IconUrl:   cuisine.IconUrl,
		URL:       cuisine.URL,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		UserId:    cuisine.UserId,
	}
	//log.Print(rescuisine)
	return rescuisine, nil
}

func (cu *cuisineUsecase) SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, userId uint, cuisineId uint) (model.CuisineResponse, error) {
	//cuisine := model.Cuisine{}

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

		cuisine.IconUrl = &img_url
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
		IconUrl:   cuisine.IconUrl,
		URL:       url,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: cuisine.UpdatedAt,
		User:      cuisine.User,
		UserId:    cuisine.UserId,
	}
	//log.Print("cuisine", cuisine)
	//log.Print("updatedCuisine", updatedCuisine)

	if err := cu.cr.SettingCuisine(&updatedCuisine); err != nil {
		return model.CuisineResponse{}, err
	}

	rescuisine := model.CuisineResponse{
		ID:        updatedCuisine.ID,
		Title:     cuisine.Title,
		IconUrl:   cuisine.IconUrl,
		URL:       cuisine.URL,
		CreatedAt: cuisine.CreatedAt,
		UpdatedAt: updatedCuisine.UpdatedAt,
		UserId:    updatedCuisine.UserId,
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
