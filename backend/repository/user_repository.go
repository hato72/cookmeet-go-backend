package repository

// emailでのユーザー検索、ユーザーテーブルの作成、ユーザーテーブルの更新処理を実装

import (
	"backend/model"
	"fmt"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// プリペアドステートメントを無効化するセッションを使用
	if err := ur.db.Session(&gorm.Session{
		PrepareStmt: false,
	}).Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	// プリペアドステートメントを無効化したトランザクションを開始
	tx := ur.db.Session(&gorm.Session{
		PrepareStmt: false,
	}).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	return tx.Commit().Error
}

func (ur *userRepository) UpdateUser(user *model.User) error {
	// プリペアドステートメントを無効化したセッションを作成
	dbSession := ur.db.Session(&gorm.Session{PrepareStmt: false})

	if user.Email != "" {
		if err := dbSession.Model(user).Where("id = ?", user.ID).Update("email", user.Email).Error; err != nil {
			return err
		}
	}

	if user.Name != "" {
		if err := dbSession.Model(user).Where("id = ?", user.ID).Update("name", user.Name).Error; err != nil {
			return err
		}
	}

	if user.Password != "" {
		if err := dbSession.Model(user).Where("id = ?", user.ID).Update("password", user.Password).Error; err != nil {
			return err
		}
	}

	if user.IconURL != nil {
		if err := dbSession.Model(user).Where("id = ? ", user.ID).Update("icon_url", user.IconURL).Error; err != nil {
			return err
		}
	}

	return nil
}
