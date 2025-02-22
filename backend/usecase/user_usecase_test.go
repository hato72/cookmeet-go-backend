package usecase

import (
	"backend/model"
	"backend/validator"
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
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := model.User{
		Email:    "test@example.com",
		Password: password,
	}

	// モックの振る舞いを設定
	mockRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*model.User)
			arg.ID = 1
			arg.Email = user.Email
			arg.Password = string(hashedPassword)
		}).Return(nil)

	// テスト実行
	tokenString, err := usecase.Login(user)

	// アサーション
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	mockRepo.AssertExpectations(t)
}
