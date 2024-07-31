// file: responder_test.go
package responder

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponder_ErrorResponses(t *testing.T) {
	responder := NewResponder()

	testCases := []struct {
		name           string
		errorFunc      func(http.ResponseWriter, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Unauthorized Error",
			errorFunc:      responder.ErrorUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "test error",
		},
		{name: "Bad Request Error",
			errorFunc:      responder.ErrorBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()

			// Вызываем тестируемую функцию
			tc.errorFunc(rr, errors.New(tc.expectedError))

			// Проверяем статус-код ответа
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// Проверяем содержимое ответа
			var response map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			if errorMsg, exists := response["error"]; !exists || errorMsg != tc.expectedError {
				t.Errorf("expected error message '%s', got '%s'", tc.expectedError, errorMsg)
			}
		})
	}
}
func TestResponder_ErrorForbidden(t *testing.T) {
	responder := NewResponder()
	rr := httptest.NewRecorder()
	testError := errors.New("forbidden error")

	responder.ErrorForbidden(rr, testError)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}

	// Проверяем содержимое ответа
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if errorMsg, exists := response["error"]; !exists || errorMsg != testError.Error() {
		t.Errorf("expected error message '%s', got '%s'", testError.Error(), errorMsg)
	}
}

func TestResponder_ErrorInternal(t *testing.T) {
	responder := NewResponder()
	rr := httptest.NewRecorder()
	testError := errors.New("internal server error")

	responder.ErrorInternal(rr, testError)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Проверяем содержимое ответа
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if errorMsg, exists := response["error"]; !exists || errorMsg != testError.Error() {
		t.Errorf("expected error message '%s', got '%s'", testError.Error(), errorMsg)
	}
}

func TestResponder_OutputJSON_Error(t *testing.T) {
	responder := NewResponder()
	rr := httptest.NewRecorder()

	// Создаем канал, который не может быть закодирован в JSON
	data := make(chan int)

	responder.OutputJSON(rr, data)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Проверяем содержимое ответа
	expectedBody := "Internal Server Error"
	actualBody := strings.TrimSpace(rr.Body.String())
	if actualBody != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, actualBody)
	}
}
