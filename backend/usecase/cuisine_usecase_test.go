package usecase

import (
	"backend/model"
	"backend/validator"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCuisineRepository はCuisineRepositoryのモック
type MockCuisineRepository struct {
	mock.Mock
}

func (m *MockCuisineRepository) GetAllCuisines(cuisines *[]model.Cuisine, userID uint) error {
	args := m.Called(cuisines, userID)
	if args.Get(0) != nil {
		*cuisines = args.Get(0).([]model.Cuisine)
	}
	return args.Error(1)
}

func (m *MockCuisineRepository) GetCuisineByID(cuisine *model.Cuisine, userID uint, cuisineID uint) error {
	args := m.Called(cuisine, userID, cuisineID)
	if args.Get(0) != nil {
		*cuisine = args.Get(0).(model.Cuisine)
	}
	return args.Error(1)
}

func (m *MockCuisineRepository) CreateCuisine(cuisine *model.Cuisine) error {
	args := m.Called(cuisine)
	return args.Error(0)
}

func (m *MockCuisineRepository) DeleteCuisine(userID uint, cuisineID uint) error {
	args := m.Called(userID, cuisineID)
	return args.Error(0)
}

func (m *MockCuisineRepository) SettingCuisine(cuisine *model.Cuisine) error {
	args := m.Called(cuisine)
	return args.Error(0)
}

func TestGetAllCuisines(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	UserID := uint(1)
	now := time.Now()
	mockCuisines := []model.Cuisine{
		{
			ID:        1,
			Title:     "Test Cuisine 1",
			URL:       "http://example.com/1",
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    UserID,
		},
		{
			ID:        2,
			Title:     "Test Cuisine 2",
			URL:       "http://example.com/2",
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    UserID,
		},
	}

	// モックの振る舞いを設定
	mockRepo.On("GetAllCuisines", mock.AnythingOfType("*[]model.Cuisine"), UserID).
		Run(func(args mock.Arguments) {
			cuisines := args.Get(0).(*[]model.Cuisine)
			*cuisines = mockCuisines
		}).
		Return(mockCuisines, nil)

	// テスト実行
	cuisines, err := usecase.GetAllCuisines(UserID)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, cuisines, 2)
	assert.Equal(t, mockCuisines[0].Title, cuisines[0].Title)
	assert.Equal(t, mockCuisines[0].URL, cuisines[0].URL)
	mockRepo.AssertExpectations(t)
}

func TestGetCuisineByID(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	UserID := uint(1)
	cuisineID := uint(1)
	now := time.Now()
	mockCuisine := model.Cuisine{
		ID:        cuisineID,
		Title:     "Test Cuisine",
		URL:       "http://example.com",
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    UserID,
	}

	// モックの振る舞いを設定
	mockRepo.On("GetCuisineByID", mock.AnythingOfType("*model.Cuisine"), UserID, cuisineID).
		Run(func(args mock.Arguments) {
			cuisine := args.Get(0).(*model.Cuisine)
			*cuisine = mockCuisine
		}).
		Return(mockCuisine, nil)

	// テスト実行
	cuisine, err := usecase.GetCuisineByID(UserID, cuisineID)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, mockCuisine.Title, cuisine.Title)
	assert.Equal(t, mockCuisine.URL, cuisine.URL)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCuisine(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	UserID := uint(1)
	cuisineID := uint(1)

	// モックの振る舞いを設定
	mockRepo.On("DeleteCuisine", UserID, cuisineID).Return(nil)

	// テスト実行
	err := usecase.DeleteCuisine(UserID, cuisineID)

	// アサーション
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddCuisine(t *testing.T) {
	// モックの準備
	mockRepo := new(MockCuisineRepository)
	validator := validator.NewCuisineValidator()
	usecase := NewCuisineUsecase(mockRepo, validator)

	cuisine := model.Cuisine{
		Title:  "Test Cuisine",
		URL:    "http://example.com",
		UserID: 1,
	}

	// モックの振る舞いを設定
	mockRepo.On("CreateCuisine", mock.AnythingOfType("*model.Cuisine")).Return(nil)

	// テスト実行
	response, err := usecase.AddCuisine(cuisine, nil, cuisine.URL, cuisine.Title)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, cuisine.Title, response.Title)
	assert.Equal(t, cuisine.URL, response.URL)
	mockRepo.AssertExpectations(t)
}
