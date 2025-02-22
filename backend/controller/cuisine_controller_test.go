package controller

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCuisineUsecase struct {
	mock.Mock
}

func (m *mockCuisineUsecase) GetAllCuisines(userId uint) ([]model.CuisineResponse, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.CuisineResponse), args.Error(1)
}

func (m *mockCuisineUsecase) GetCuisineById(userId uint, cuisineId uint) (model.CuisineResponse, error) {
	args := m.Called(userId, cuisineId)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

func (m *mockCuisineUsecase) DeleteCuisine(userId uint, cuisineId uint) error {
	args := m.Called(userId, cuisineId)
	return args.Error(0)
}

func (m *mockCuisineUsecase) AddCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string) (model.CuisineResponse, error) {
	args := m.Called(cuisine, iconFile, url, title)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

func (m *mockCuisineUsecase) SetCuisine(cuisine model.Cuisine, iconFile *multipart.FileHeader, url string, title string, userId uint, cuisineId uint) (model.CuisineResponse, error) {
	args := m.Called(cuisine, iconFile, url, title, userId, cuisineId)
	return args.Get(0).(model.CuisineResponse), args.Error(1)
}

func setupCuisineTest(t *testing.T) (*echo.Echo, *mockCuisineUsecase, ICuisineController) {
	e := echo.New()
	mockUsecase := new(mockCuisineUsecase)
	controller := NewCuisineController(mockUsecase)
	return e, mockUsecase, controller
}

func createJWTToken(userId float64) *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId
	return token
}

func TestGetAllCuisines(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userId       float64
		mockResponse []model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:   "正常な取得",
			userId: 1,
			mockResponse: []model.CuisineResponse{
				{
					ID:        1,
					Title:     "Test Cuisine 1",
					URL:       "https://example.com/1",
					UserId:    1,
					CreatedAt: time.Now(),
				},
				{
					ID:        2,
					Title:     "Test Cuisine 2",
					URL:       "https://example.com/2",
					UserId:    1,
					CreatedAt: time.Now(),
				},
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
		{
			name:         "データなし",
			userId:       2,
			mockResponse: []model.CuisineResponse{},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cuisines", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tc.userId))

			mockUsecase.On("GetAllCuisines", uint(tc.userId)).Return(tc.mockResponse, tc.mockError)

			err := controller.GetAllCuisines(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response []model.CuisineResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, len(tc.mockResponse), len(response))
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetCuisineById(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userId       float64
		cuisineId    string
		mockResponse model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:      "正常な取得",
			userId:    1,
			cuisineId: "1",
			mockResponse: model.CuisineResponse{
				ID:        1,
				Title:     "Test Cuisine",
				URL:       "https://example.com",
				UserId:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cuisines/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineId")
			c.SetParamValues(tc.cuisineId)
			c.Set("user", createJWTToken(tc.userId))

			mockUsecase.On("GetCuisineById", uint(tc.userId), uint(1)).Return(tc.mockResponse, tc.mockError)

			err := controller.GetCuisineById(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tc.mockResponse.ID, response.ID)
				assert.Equal(t, tc.mockResponse.Title, response.Title)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDeleteCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userId       float64
		cuisineId    string
		mockError    error
		expectStatus int
	}{
		{
			name:         "正常な削除",
			userId:       1,
			cuisineId:    "1",
			mockError:    nil,
			expectStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/cuisines/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineId")
			c.SetParamValues(tc.cuisineId)
			c.Set("user", createJWTToken(tc.userId))

			mockUsecase.On("DeleteCuisine", uint(tc.userId), uint(1)).Return(tc.mockError)

			err := controller.DeleteCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAddCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userId       float64
		title        string
		url          string
		mockResponse model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:   "正常な追加",
			userId: 1,
			title:  "New Cuisine",
			url:    "https://example.com/new",
			mockResponse: model.CuisineResponse{
				ID:        1,
				Title:     "New Cuisine",
				URL:       "https://example.com/new",
				UserId:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("title", tc.title)
			writer.WriteField("url", tc.url)
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/cuisines", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", createJWTToken(tc.userId))

			mockUsecase.On("AddCuisine",
				mock.AnythingOfType("model.Cuisine"),
				(*multipart.FileHeader)(nil),
				tc.url,
				tc.title,
			).Return(tc.mockResponse, tc.mockError)

			err := controller.AddCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tc.mockResponse.Title, response.Title)
				assert.Equal(t, tc.mockResponse.URL, response.URL)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSetCuisine(t *testing.T) {
	e, mockUsecase, controller := setupCuisineTest(t)

	testCases := []struct {
		name         string
		userId       float64
		cuisineId    string
		title        string
		url          string
		mockGetRes   model.CuisineResponse
		mockSetRes   model.CuisineResponse
		mockError    error
		expectStatus int
	}{
		{
			name:      "正常な更新",
			userId:    1,
			cuisineId: "1",
			title:     "Updated Cuisine",
			url:       "https://example.com/updated",
			mockGetRes: model.CuisineResponse{
				ID:        1,
				Title:     "Original Cuisine",
				URL:       "https://example.com/original",
				UserId:    1,
				CreatedAt: time.Now(),
			},
			mockSetRes: model.CuisineResponse{
				ID:        1,
				Title:     "Updated Cuisine",
				URL:       "https://example.com/updated",
				UserId:    1,
				CreatedAt: time.Now(),
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("title", tc.title)
			writer.WriteField("url", tc.url)
			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/cuisines/:id", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("cuisineId")
			c.SetParamValues(tc.cuisineId)
			c.Set("user", createJWTToken(tc.userId))

			mockUsecase.On("GetCuisineById", uint(tc.userId), uint(1)).Return(tc.mockGetRes, nil)
			mockUsecase.On("SetCuisine",
				mock.AnythingOfType("model.Cuisine"),
				(*multipart.FileHeader)(nil),
				tc.url,
				tc.title,
				uint(tc.userId),
				uint(1),
			).Return(tc.mockSetRes, tc.mockError)

			err := controller.SetCuisine(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusOK {
				var response model.CuisineResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tc.mockSetRes.Title, response.Title)
				assert.Equal(t, tc.mockSetRes.URL, response.URL)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
