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

func TestSignUp(t *testing.T) {
	// モックの準備
	mockRepo := new(MockUserRepository)
	validator := validator.NewUserValidator()
	usecase := NewUserUsecase(mockRepo, validator)

	// テストケース
	user := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// モックの振る舞いを設定
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	// テスト実行
	response, err := usecase.SignUp(user)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	// モックの準備
	mockRepo := new(MockUserRepository)
	validator := validator.NewUserValidator()
	usecase := NewUserUsecase(mockRepo, validator)

	// テストケース
	// password := "password123"
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// user := model.User{ //正しいケース
	// 	Email:    "test@example.com",
	// 	Password: password,
	// }

	// // モックの振る舞いを設定
	// mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).
	// 	Run(func(args mock.Arguments) {
	// 		arg := args.Get(0).(*model.User)
	// 		arg.ID = 1
	// 		arg.Email = user.Email
	// 		arg.Password = string(hashedPassword)
	// 	}).Return(nil)

	// // テスト実行
	// tokenString, err := usecase.Login(user)

	// // アサーション
	// assert.NoError(t, err)
	// assert.NotEmpty(t, tokenString)

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
