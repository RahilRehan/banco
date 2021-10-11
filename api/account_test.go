package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RahilRehan/banco/db/mocks"
	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/RahilRehan/banco/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := map[string]struct {
		accountID      int64
		expectedStatus int
		stub           func() *mocks.Store
	}{
		"Status OK": {
			accountID:      account.ID,
			expectedStatus: http.StatusOK,
			stub: func() *mocks.Store {
				mockStore := new(mocks.Store)
				mockStore.On("GetAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(*account, nil)
				return mockStore
			},
		},
		"Not Found": {
			accountID:      account.ID,
			expectedStatus: http.StatusNotFound,
			stub: func() *mocks.Store {
				mockStore := new(mocks.Store)
				mockStore.On("GetAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(db.Account{}, sql.ErrNoRows)
				return mockStore
			},
		},
		"Internal Server Error": {
			accountID:      account.ID,
			expectedStatus: http.StatusInternalServerError,
			stub: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("GetAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(db.Account{}, errors.New("internal server error"))
				return mocksStore
			},
		},
		"Bad request": {
			accountID:      0,
			expectedStatus: http.StatusBadRequest,
			stub: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("GetAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(db.Account{}, errors.New("bad request"))
				return mocksStore
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockStore := test.stub()
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", test.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server := NewServer(mockStore)
			server.router.ServeHTTP(recorder, req)

			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestCreateAccount(t *testing.T) {

	account := randomAccount()

	testCases := map[string]struct {
		body           gin.H
		expectedStatus int
		stubs          func() *mocks.Store
	}{
		"Status OK": {
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			expectedStatus: http.StatusCreated,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateAccountParams")).Return(*account, nil)
				return mocksStore
			},
		},
		"Bad Request": {
			body: gin.H{
				"currency": "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateAccountParams")).Return(*account, nil)
				return mocksStore
			},
		},
		"Internal server error": {
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			expectedStatus: http.StatusInternalServerError,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("CreateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.CreateAccountParams")).Return(db.Account{}, errors.New("internal error"))
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

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestListAccounts(t *testing.T) {

	n := 5
	accounts := make([]db.Account, 5)
	for i := 0; i < n; i++ {
		accounts[i] = *randomAccount()
	}

	type query struct {
		PageID   int
		PageSize int
	}

	testCases := map[string]struct {
		query          query
		expectedStatus int
		stubs          func() *mocks.Store
	}{
		"Status OK": {
			query:          query{PageID: 1, PageSize: n},
			expectedStatus: http.StatusOK,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(accounts, nil)
				return mocksStore
			},
		},
		"Bad Request": {
			query:          query{},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(accounts, nil)
				return mocksStore
			},
		},
		"Internal server error": {
			query:          query{PageID: 1, PageSize: n},
			expectedStatus: http.StatusInternalServerError,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(nil, errors.New("internal error"))
				return mocksStore
			},
		},
		"Invalid page id": {
			query:          query{PageID: 0, PageSize: n},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(accounts, nil)
				return mocksStore
			},
		},
		"Invalid page size": {
			query:          query{PageID: 1, PageSize: 50},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(accounts, nil)
				return mocksStore
			},
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {

			mockStore := test.stubs()

			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", test.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", test.query.PageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestUpdateAccount(t *testing.T) {

	account := randomAccount()
	extra := 100
	testCases := map[string]struct {
		body           gin.H
		expectedStatus int
		stubs          func() *mocks.Store
	}{
		"Status OK": {
			body: gin.H{
				"id":      account.ID,
				"balance": account.Balance + int64(extra),
			},
			expectedStatus: http.StatusOK,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("UpdateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.UpdateAccountParams")).Return(db.Account{ID: account.ID, Currency: account.Currency, Balance: account.Balance + int64(extra), Owner: account.Owner, CreatedAt: account.CreatedAt.Add(time.Second)}, nil)
				return mocksStore
			},
		},
		"Bad Request": {
			body: gin.H{
				"bad": "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("UpdateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.UpdateAccountParams")).Return(db.Account{}, nil)
				return mocksStore
			},
		},
		"Internal server error": {
			body: gin.H{
				"id":      account.ID,
				"balance": account.Balance + int64(extra),
			},
			expectedStatus: http.StatusInternalServerError,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("UpdateAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.UpdateAccountParams")).Return(db.Account{}, errors.New("internal error"))
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

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	account := randomAccount()

	testCases := map[string]struct {
		accountID      int64
		expectedStatus int
		stub           func() *mocks.Store
	}{
		"Status OK": {
			accountID:      account.ID,
			expectedStatus: http.StatusOK,
			stub: func() *mocks.Store {
				mockStore := new(mocks.Store)
				mockStore.On("DeleteAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(nil)
				return mockStore
			},
		},
		"Internal Server Error": {
			accountID:      account.ID,
			expectedStatus: http.StatusInternalServerError,
			stub: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("DeleteAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(errors.New("internal server error"))
				return mocksStore
			},
		},
		"Bad request": {
			accountID:      0,
			expectedStatus: http.StatusBadRequest,
			stub: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("DeleteAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(errors.New("bad request"))
				return mocksStore
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockStore := test.stub()
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", test.accountID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server := NewServer(mockStore)
			server.router.ServeHTTP(recorder, req)

			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func randomAccount() *db.Account {
	return &db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomString(6),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
