package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	expected := `{"alive": true}`

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "they should be equal")
	assert.Equal(t, expected, rr.Body.String(), "they should be equal")

	//Convey("Создаем реквест для /health-check", t, func() {
	//
	//	Convey("Когда запрашиваем данные", func() {
	//
	//
	//		Convey("Должен вернуться правильный статус код", func() {
	//
	//		})
	//
	//		Convey("Должны вернуться верные данные в формате json", func() {
	//
	//		})
	//
	//	})
	//})

	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	//}

	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}

	//type args struct {
	//	w http.ResponseWriter
	//	r *http.Request
	//}
	//tests := []struct {
	//	name string
	//	args args
	//}{
	//	// TODыO: Add test cases.
	//
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//	})
	//}
}
