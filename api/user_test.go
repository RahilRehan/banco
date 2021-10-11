package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RahilRehan/banco/db/mocks"
	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/RahilRehan/banco/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {

	password := "tester"
	hashPass, err := util.HashPassword(password)
	require.NoError(t, err)
	user := randomUser(password)
	dbUser := &db.User{
		Username:       user.Username,
		Email:          user.Email,
		FullName:       user.FullName,
		HashedPassword: hashPass,
	}

	testCases := map[string]struct {
		body           gin.H
		expectedStatus int
		stubs          func() *mocks.Store
	}{
		"Status OK": {
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  user.Password,
			},
			expectedStatus: http.StatusCreated,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateUserParams")).Return(db.User{}, nil)
				return mocksStore
			},
		},
		"Small password length": {
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  "small",
			},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateUserParams")).Return(*dbUser, nil)
				return mocksStore
			},
		},
		"Invalid username": {
			body: gin.H{
				"username":  "ra%#",
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  user.Password,
			},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateUserParams")).Return(*dbUser, nil)
				return mocksStore
			},
		},
		"Same username twice": {
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  user.Password,
			},
			expectedStatus: http.StatusInternalServerError,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateUserParams")).Return(db.User{}, errors.New("username twice"))
				return mocksStore
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			mockStore := test.stubs()

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			url := "/users/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func randomUser(password string) *createUserRequest {
	return &createUserRequest{
		Username: util.RandomOwner(),
		FullName: util.RandomString(10),
		Email:    util.RandomEmail(),
		Password: password,
	}
}
