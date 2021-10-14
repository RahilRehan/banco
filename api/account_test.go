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
	"github.com/RahilRehan/banco/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	user := randomUser("temp")
	account := randomAccount(user.Username)

	testCases := map[string]struct {
		accountID      int64
		expectedStatus int
		stub           func() *mocks.Store
		setupAuth      func(t *testing.T, req *http.Request, maker token.Maker)
	}{
		"Status OK": {
			accountID:      account.ID,
			expectedStatus: http.StatusOK,
			stub: func() *mocks.Store {
				mockStore := new(mocks.Store)
				mockStore.On("GetAccount", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).Return(*account, nil)
				return mockStore
			},
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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

			server := newTestServer(t, mockStore)
			test.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)

			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	user := randomUser("temp")
	account := randomAccount(user.Username)

	testCases := map[string]struct {
		body           gin.H
		expectedStatus int
		stubs          func() *mocks.Store
		setupAuth      func(t *testing.T, req *http.Request, maker token.Maker)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {

			mockStore := test.stubs()

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			test.setupAuth(t, request, server.tokenMaker)

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestListAccounts(t *testing.T) {
	user := randomUser("temp")
	n := 5
	accounts := make([]db.Account, 5)
	for i := 0; i < n; i++ {
		accounts[i] = *randomAccount(user.Username)
	}

	type query struct {
		PageID   int
		PageSize int
	}

	testCases := map[string]struct {
		query          query
		expectedStatus int
		stubs          func() *mocks.Store
		setupAuth      func(t *testing.T, req *http.Request, maker token.Maker)
	}{
		"Status OK": {
			query:          query{PageID: 1, PageSize: n},
			expectedStatus: http.StatusOK,
			stubs: func() *mocks.Store {
				mocksStore := new(mocks.Store)
				mocksStore.On("ListAccounts", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("db.ListAccountsParams")).Return(accounts, nil)
				return mocksStore
			},
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {

			mockStore := test.stubs()

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", test.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", test.query.PageSize))
			request.URL.RawQuery = q.Encode()

			test.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func TestUpdateAccount(t *testing.T) {

	user := randomUser("temp")
	account := randomAccount(user.Username)

	extra := 100
	testCases := map[string]struct {
		body           gin.H
		expectedStatus int
		stubs          func() *mocks.Store
		setupAuth      func(t *testing.T, req *http.Request, maker token.Maker)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
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
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, user.Username, time.Minute)
			},
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {

			mockStore := test.stubs()

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			url := "/accounts/"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			test.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

func randomAccount(username string) *db.Account {
	return &db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
