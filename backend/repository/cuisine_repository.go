package repository

//GetAllCuisines:料理データベースの一覧から引数のユーザーidに一致する料理を取得する
//GetCuisineById:引数のユーザーidに一致する料理を取得し、その中でcuisineの主キーが引数で受け取ったcuisineIdに一致する料理を取得する
//CreateCuisine:料理を作成する
//DeleteCuisine:料理を削除する
//SettingCuisine:料理を更新する

import (
	"backend/model"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICuisineRepository interface {
	GetAllCuisines(cuisines *[]model.Cuisine, userId uint) error              //作成した料理の一覧を取得
	GetCuisineById(cuisine *model.Cuisine, userId uint, cuisineId uint) error //引数のcuisineIdに一致する料理を返す
	CreateCuisine(cuisine *model.Cuisine) error                               //料理の新規作成
	//UpdateCuisine(cuisine *model.Cuisine, userId uint, cuisineId uint) error  //料理の更新
	DeleteCuisine(userId uint, cuisineId uint) error //料理の削除
	SettingCuisine(cuisine *model.Cuisine) error
}

type cuisineRepository struct {
	db *gorm.DB
}

func NewCuisineRepository(db *gorm.DB) ICuisineRepository { //コンストラクタ
	return &cuisineRepository{db}
}

func (cr *cuisineRepository) GetAllCuisines(cuisines *[]model.Cuisine, userId uint) error {
	if err := cr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(cuisines).Error; err != nil { //料理の一覧から引数のユーザーidに一致する料理を取得する　その時、作成日時があたらしいものが末尾に来るようにする
		return err
	}
	return nil
}

func (cr *cuisineRepository) GetCuisineById(cuisine *model.Cuisine, userId uint, cuisineId uint) error {
	result := cr.db.Joins("User").Where("user_id=? AND cuisines.id=?", userId, cuisineId).First(cuisine)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (cr *cuisineRepository) CreateCuisine(cuisine *model.Cuisine) error {
	if cuisine.Title == "" {
		return fmt.Errorf("title is required")
	}
	return cr.db.Create(cuisine).Error
}

func (cr *cuisineRepository) DeleteCuisine(userId uint, cuisineId uint) error {
	result := cr.db.Where("id=? AND user_id=?", cuisineId, userId).Delete(&model.Cuisine{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 { //更新されたレコードの数を取得できる
		return fmt.Errorf("object does not exists")
	}
	return nil
}

func (cr *cuisineRepository) SettingCuisine(cuisine *model.Cuisine) error {
	if cuisine.IconURL != nil {
		icon_result := cr.db.Model(cuisine).Clauses(clause.Returning{}).Where("id=? AND user_id=?", cuisine.ID, cuisine.UserId).Update("icon_url", cuisine.IconURL)
		if icon_result.Error != nil {
			return icon_result.Error
		}
		if icon_result.RowsAffected < 1 {
			return fmt.Errorf("object does not exists")
		}
	}

	if cuisine.URL != "" {
		url_result := cr.db.Model(cuisine).Clauses(clause.Returning{}).Where("id=? AND user_id=?", cuisine.ID, cuisine.UserId).Update("url", cuisine.URL)
		if url_result.Error != nil {
			return url_result.Error
		}
		if url_result.RowsAffected < 1 {
			return fmt.Errorf("object does not exists")
		}
	}
	return nil
}
