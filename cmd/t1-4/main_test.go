package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserViewHandler(t *testing.T) {
	users := make(map[string]User)
	u1 := User{
		ID:        "u1",
		FirstName: "Misha",
		LastName:  "Popov",
	}
	u2 := User{
		ID:        "u2",
		FirstName: "Sasha",
		LastName:  "Popov",
	}
	users["u1"] = u1
	users["u2"] = u2
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		req  string
		want want
	}{
		{
			name: "Test 200-1",
			req:  "/users?user_id=u1",
			want: want{
				code:        200,
				response:    `{"ID":"u1","FirstName":"Misha","LastName":"Popov"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Test 200-2",
			req:  "/users?user_id=u2",
			want: want{
				code:        200,
				response:    `{"ID":"u2","FirstName":"Sasha","LastName":"Popov"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Test 404",
			req:  "/users?user_id=u3",
			want: want{
				code:        404,
				response:    ``,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.req, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UserViewHandler(users))
			h(w, request)
			res := w.Result()

			require.Equal(t, test.want.code, res.StatusCode)
			if test.want.code == http.StatusOK {
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				assert.JSONEq(t, test.want.response, string(resBody))

				require.NoError(t, err)

			}

			// assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
	// Correct version
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		request := httptest.NewRequest(http.MethodPost, tt.request, nil)
	// 		w := httptest.NewRecorder()
	// 		h := http.HandlerFunc(UserViewHandler(tt.users))
	// 		h(w, request)

	// 		result := w.Result()

	// 		assert.Equal(t, tt.want.statusCode, result.StatusCode)
	// 		assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

	// 		userResult, err := ioutil.ReadAll(result.Body)
	// 		require.NoError(t, err)
	// 		err = result.Body.Close()
	// 		require.NoError(t, err)

	// 		var user User
	// 		err = json.Unmarshal(userResult, &user)
	// 		require.NoError(t, err)

	// 		assert.Equal(t, tt.want.user, user)
	// 	})
	// }
}
