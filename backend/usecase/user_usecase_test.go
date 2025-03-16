package usecase

import (
	"backend/model"
	"backend/validator"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository はUserRepositoryのモック
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(user *model.User, email string) error {
	args := m.Called(user, email)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserById(userId uint) (*model.User, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockUserValidator struct {
	mock.Mock
}

func (m *MockUserValidator) UserValidate(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignUp(t *testing.T) {
	// テストケース1: 正常なサインアップ（ユーザーが存在しない）
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockValidator := new(MockUserValidator)

		user := model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password",
		}

		// バリデーションが成功することを設定
		mockValidator.On("UserValidate", mock.AnythingOfType("model.User")).Return(nil)

		// GetUserByEmailがユーザーが見つからないエラーを返すように設定（正常：既存ユーザーがいない）
		mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), "test@example.com").Return(errors.New("user not found"))

		// CreateUserが成功することを設定
		mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*model.User)
			userArg.ID = 1 // IDをセット
		})

		usecase := NewUserUsecase(mockRepo, mockValidator)
		res, err := usecase.SignUp(user)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, user.Name, res.Name)
		assert.Equal(t, user.Email, res.Email)
		mockRepo.AssertExpectations(t)
		mockValidator.AssertExpectations(t)
	})

	// テストケース2: ユーザーがすでに存在する場合（エラー）
	t.Run("user already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockValidator := new(MockUserValidator)

		user := model.User{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "password",
		}

		// バリデーションが成功することを設定
		mockValidator.On("UserValidate", mock.AnythingOfType("model.User")).Return(nil)

		// GetUserByEmailがnilを返す（異常：ユーザーが既に存在する）
		mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), "existing@example.com").Return(nil)

		usecase := NewUserUsecase(mockRepo, mockValidator)
		_, err := usecase.SignUp(user)

		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
		mockRepo.AssertExpectations(t)
		mockValidator.AssertExpectations(t)
	})

	// テストケース3: バリデーションエラー
	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockValidator := new(MockUserValidator)

		user := model.User{
			Name:     "", // 名前が空
			Email:    "invalid@example.com",
			Password: "pw", // パスワードが短すぎる
		}

		validationErr := errors.New("validation error")
		mockValidator.On("UserValidate", mock.AnythingOfType("model.User")).Return(validationErr)

		usecase := NewUserUsecase(mockRepo, mockValidator)
		_, err := usecase.SignUp(user)

		assert.Error(t, err)
		assert.Equal(t, validationErr, err)
		mockValidator.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	// モックの準備
	mockRepo := new(MockUserRepository)
	validator := validator.NewUserValidator()
	usecase := NewUserUsecase(mockRepo, validator)

	// 正しいケース
	t.Run("valid login", func(t *testing.T) {
		user := model.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		// 正しいパスワードのハッシュを生成
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
		// モックの振る舞いを設定
		mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*model.User)
				arg.ID = 1
				arg.Email = user.Email
				arg.Password = string(hashedPassword)
			}).Return(nil).Once()

		tokenString, err := usecase.Login(user)
		assert.NoError(t, err, "unexpected error in valid login: %v", err)
		assert.NotEmpty(t, tokenString, "token must not be empty")
	})

	// 存在しないユーザーの場合
	t.Run("user not found", func(t *testing.T) {
		noexistuser := model.User{
			Email:    "noexist@example.com",
			Password: "password123",
		}
		// ユーザーが見つからないエラーを返す
		mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), noexistuser.Email).
			Return(errors.New("user not found")).Once()

		_, err := usecase.Login(noexistuser)
		assert.Error(t, err, "expected error for non-existent user")
		assert.Truef(t, errors.Is(err, ErrUserNotFound), "expected ErrUserNotFound, but got: %v", err)
	})

	// パスワードが間違っている場合
	t.Run("invalid password", func(t *testing.T) {
		misspassuser := model.User{
			Email:    "noexist@example.com",
			Password: "password12345",
		}
		// 正しいパスワードのハッシュを保持したユーザー情報を返す
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
		mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), misspassuser.Email).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*model.User)
				arg.ID = 2
				arg.Email = misspassuser.Email
				arg.Password = string(hashedPassword)
			}).Return(nil).Once()

		_, err := usecase.Login(misspassuser)
		assert.Error(t, err, "expected error for invalid password")
		assert.Truef(t, errors.Is(err, ErrInvalidPassword), "expected ErrInvalidPassword, but got: %v", err)
	})

	mockRepo.AssertExpectations(t)
}
