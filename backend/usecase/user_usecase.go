package usecase

// サインアップ、ログイン、更新処理を実装
// サインアップでは、user_validatorを呼び出したのち、user_repositoryのユーザーテーブル作成メソッドを呼び出している
// ログインでは、user_repositoryのemailでのユーザー検索メソッドを呼び出したのち、jwtトークンの検証を行っている
// 更新処理では、更新情報があればデータの更新を行っている

import (
	"backend/model"
	"backend/repository"
	"backend/validator"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// エラー定義を追加
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
	Update(user model.User, newEmail string, newName string, newPassword string, iconFile *multipart.FileHeader) (model.UserResponse, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}

	// まず既存ユーザーがいないか確認
	existingUser := model.User{}
	if err := uu.ur.GetUserByEmail(&existingUser, user.Email); err == nil {
		// エラーがない場合はユーザーが見つかっている（既に存在する）
		return model.UserResponse{}, ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Name: user.Name, Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		// データベースエラーの場合も、重複に関するエラーかどうかをチェック
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique violation") {
			return model.UserResponse{}, ErrUserAlreadyExists
		}
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:      newUser.ID,
		Name:    newUser.Name,
		Email:   newUser.Email,
		IconURL: newUser.IconURL,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	storedUser := model.User{} // 空のユーザーオブジェクト
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", ErrUserNotFound
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)) // パスワードの検証
	if err != nil {
		// エラーをラップすることで、errors.Isでの判定が成功するようにする
		return "", fmt.Errorf("password mismatch: %w", ErrInvalidPassword)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(), // jwtの有効期限
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET"))) // jwtトークンの生成
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu *userUsecase) Update(user model.User, newEmail string, newName string, newPassword string, iconFile *multipart.FileHeader) (model.UserResponse, error) {

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return model.UserResponse{}, err
		}
		user.Password = string(hash)
	}

	if iconFile != nil {
		src, err := iconFile.Open()
		if err != nil {
			return model.UserResponse{}, err
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			return model.UserResponse{}, err
		}

		hasher := sha256.New()
		hasher.Write(data)
		hashValue := hex.EncodeToString(hasher.Sum(nil))

		ext := filepath.Ext(iconFile.Filename)

		IconURL := "icons/" + hashValue + ext

		dst, err := os.Create("./user_images/" + IconURL)
		if err != nil {
			return model.UserResponse{}, err
		}

		defer dst.Close()

		if _, err := dst.Write(data); err != nil {
			return model.UserResponse{}, nil
		}

		user.IconURL = &IconURL

	}

	updatedUser := model.User{
		ID:       user.ID,
		Name:     newName,
		Email:    newEmail,
		Password: newPassword,
		IconURL:  user.IconURL,
	}
	// log.Print("updateUser:", updatedUser)

	if err := uu.ur.UpdateUser(&updatedUser); err != nil {
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:      updatedUser.ID,
		Name:    updatedUser.Name,
		Email:   updatedUser.Email,
		IconURL: updatedUser.IconURL,
	}
	// log.Print("resUser:", resUser)

	return resUser, nil

}
