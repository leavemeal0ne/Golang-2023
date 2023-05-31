package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {

	////db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go user=postgres password=1209")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	////dbRepo := getDatabase(&database.DbRepo{DB: db})

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"email":         "joeqeriwewqk@ukr.net",
				"password_hash": "we3123sad",
			},
			expectedCode: http.StatusOK,
		},
	}

	handler := http.HandlerFunc(SignUpJson)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			err := json.NewEncoder(b).Encode(tc.payload)
			if err != nil {
				return
			}
			req, _ := http.NewRequest("GET", "/api/user/signup", b)
			req.Header.Add("Content-Type", "application/json")
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
