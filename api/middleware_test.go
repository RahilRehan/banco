package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RahilRehan/banco/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(t *testing.T, req *http.Request, maker token.Maker, authorizationType string, username string, duration time.Duration) {
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestMiddleware(t *testing.T) {
	testCases := map[string]struct {
		setupAuth     func(t *testing.T, req *http.Request, maker token.Maker)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		"OK": {
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, "username", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		"No Auth": {
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		"Unsupported Auth": {
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, "unsupported", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		"Invalid Auth Format": {
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, "", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		"Expired token": {
			setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
				addAuth(t, req, maker, authorizationTypeBearer, "username", -time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})
				},
			)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			test.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(rec, req)
			test.checkResponse(t, rec)
		})
	}
}
