package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonPaus/exporter/internal/storages/memory"
	"github.com/stretchr/testify/require"
)

func Test_updateMetric(t *testing.T) {
	// type args struct {
	// 	res http.ResponseWriter
	// 	req *http.Request
	// }
	storage := memory.NewMemory()
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		req    string
		method string
		want   want
	}{
		{
			name:   "Test 200",
			req:    "/update/counter/testCounter/100",
			method: http.MethodPost,
			want: want{
				code:        200,
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong method",
			req:    "/update/counter/testCounter/100",
			method: http.MethodGet,
			want: want{
				code:        405,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(test.method, test.req, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { updateMetric(w, r, storage) })
			h(w, request)
			// fmt.Println("Response Status:", test.method)
			res := w.Result()
			defer res.Body.Close()
			require.Equal(t, test.want.code, res.StatusCode)
			// if test.want.code == http.StatusOK {
			// 	// получаем и проверяем тело запроса
			// 	defer res.Body.Close()
			// 	resBody, err := io.ReadAll(res.Body)
			// 	assert.JSONEq(t, test.want.response, string(resBody))

			// 	require.NoError(t, err)

			// }
			// updateMetric(tt.args.res, tt.args.req)
		})
	}
}
