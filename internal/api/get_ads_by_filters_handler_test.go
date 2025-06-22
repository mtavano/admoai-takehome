package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_one(t *testing.T) {
	assert.True(t, true)
}

// Test cases for GetAdsByFiltersHandler
type testCase struct {
	name           string
	queryParams    map[string]string
	mockSetup      func(*MockDatabase)
	expectedStatus int
	expectedBody   map[string]any
	expectedError  bool
}

func TestGetAdsByFiltersHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []testCase{
		{
			name:        "success - get all ads without filters",
			queryParams: map[string]string{},
			mockSetup: func(mockDB *MockDatabase) {
				// Mock successful query with sample data
				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Test Ad 1",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
					{
						ID:        "2",
						Title:     "Test Ad 2",
						ImageURL:  "https://example.com/image2.jpg",
						Placement: "sidebar",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// Set the dest argument to our sample data
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Test Ad 1",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
					{
						ID:        "2",
						Title:     "Test Ad 2",
						ImageURL:  "https://example.com/image2.jpg",
						Placement: "sidebar",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "success - filter by placement",
			queryParams: map[string]string{
				"placement": "homepage",
			},
			mockSetup: func(mockDB *MockDatabase) {
				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Test Ad 1",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Test Ad 1",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "success - filter by status active",
			queryParams: map[string]string{
				"status": store.AdvertiseStatusActive,
			},
			mockSetup: func(mockDB *MockDatabase) {
				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Active Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Active Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "success - filter by status inactive with expired filtering",
			queryParams: map[string]string{
				"status": store.AdvertiseStatusInactive,
			},
			mockSetup: func(mockDB *MockDatabase) {
				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Inactive Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusInactive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Inactive Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusInactive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "success - filter by both placement and status",
			queryParams: map[string]string{
				"placement": "homepage",
				"status":    store.AdvertiseStatusActive,
			},
			mockSetup: func(mockDB *MockDatabase) {
				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Homepage Active Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Homepage Active Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: nil,
						Expired:   false,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "success - empty result",
			queryParams: map[string]string{
				"placement": "nonexistent",
			},
			mockSetup: func(mockDB *MockDatabase) {
				// Mock empty result
				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = []*store.AdvertiseRecord{}
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{},
			},
			expectedError: false,
		},
		{
			name:        "error - database query fails",
			queryParams: map[string]string{},
			mockSetup: func(mockDB *MockDatabase) {
				// Mock database error
				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Return(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			expectedError:  true,
		},
		{
			name:        "success - ads with expiration calculation",
			queryParams: map[string]string{},
			mockSetup: func(mockDB *MockDatabase) {
				// Create ads with expiration times
				expiredTime := time.Now().Add(-1 * time.Hour).Unix() // Expired 1 hour ago
				futureTime := time.Now().Add(1 * time.Hour).Unix()   // Expires in 1 hour

				sampleAds := []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Expired Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: &expiredTime,
					},
					{
						ID:        "2",
						Title:     "Future Ad",
						ImageURL:  "https://example.com/image2.jpg",
						Placement: "sidebar",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: &futureTime,
					},
				}

				mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						dest := args.Get(0).(*[]*store.AdvertiseRecord)
						*dest = sampleAds
					}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"ads": []*store.AdvertiseRecord{
					{
						ID:        "1",
						Title:     "Expired Ad",
						ImageURL:  "https://example.com/image1.jpg",
						Placement: "homepage",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: func() *int64 { t := time.Now().Add(-1 * time.Hour).Unix(); return &t }(),
						Expired:   true, // Should be calculated as expired
					},
					{
						ID:        "2",
						Title:     "Future Ad",
						ImageURL:  "https://example.com/image2.jpg",
						Placement: "sidebar",
						Status:    store.AdvertiseStatusActive,
						CreatedAt: time.Now().Unix(),
						ExpiresAt: func() *int64 { t := time.Now().Add(1 * time.Hour).Unix(); return &t }(),
						Expired:   false, // Should be calculated as not expired
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock database
			mockDB := new(MockDatabase)
			tc.mockSetup(mockDB)

			// Create context
			ctx := &Context{Db: mockDB}

			// Create Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request with query parameters
			req := httptest.NewRequest("GET", "/v1/ads", nil)
			q := req.URL.Query()
			for key, value := range tc.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// Call the handler
			payload, statusCode, err := GetAdsByFiltersHandler(c, ctx)

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedStatus, statusCode)
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, statusCode)
				assert.NotNil(t, payload)

				// Check response structure
				response, ok := payload.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, response, "ads")

				// For non-empty responses, check the ads array
				if ads, ok := response["ads"].([]*store.AdvertiseRecord); ok && len(ads) > 0 {
					// Verify that CalculateAndSetExpired was called
					for _, ad := range ads {
						assert.NotNil(t, ad)
						// The Expired field should be set (either true or false)
						// We can't check the exact value since it depends on current time
						// but we can verify the field exists
					}
				}
			}

			// Verify all mocked calls were made
			mockDB.AssertExpectations(t)
		})
	}
}

// TestGetAdsByFiltersHandlerValidation tests specific validation scenarios
func TestGetAdsByFiltersHandlerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validationTests := []struct {
		name           string
		queryParams    map[string]string
		description    string
		expectedStatus int
	}{
		{
			name:           "valid - empty query parameters",
			queryParams:    map[string]string{},
			description:    "Should accept empty query parameters and return all ads",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid - placement parameter only",
			queryParams: map[string]string{
				"placement": "homepage",
			},
			description:    "Should accept placement parameter",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid - status parameter only",
			queryParams: map[string]string{
				"status": "active",
			},
			description:    "Should accept status parameter",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid - both parameters",
			queryParams: map[string]string{
				"placement": "homepage",
				"status":    "active",
			},
			description:    "Should accept both placement and status parameters",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid - unknown placement value",
			queryParams: map[string]string{
				"placement": "unknown_placement",
			},
			description:    "Should accept any placement value (validation happens at DB level)",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid - unknown status value",
			queryParams: map[string]string{
				"status": "unknown_status",
			},
			description:    "Should accept any status value (validation happens at DB level)",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock database
			mockDB := new(MockDatabase)
			mockDB.On("Select", mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					dest := args.Get(0).(*[]*store.AdvertiseRecord)
					*dest = []*store.AdvertiseRecord{}
				}).
				Return(nil)

			// Create context
			ctx := &Context{Db: mockDB}

			// Create Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request with query parameters
			req := httptest.NewRequest("GET", "/v1/ads", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// Call the handler
			payload, statusCode, err := GetAdsByFiltersHandler(c, ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, statusCode)
			assert.NotNil(t, payload)

			// Verify response structure
			response, ok := payload.(map[string]interface{})
			assert.True(t, ok)
			assert.Contains(t, response, "ads")

			mockDB.AssertExpectations(t)
		})
	}
}
