package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	report "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReportUsecase struct{ mock.Mock }

func (m *MockReportUsecase) GetCategorySpending(userID int64, month, year int) ([]report.CategoryReport, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).([]report.CategoryReport), args.Error(1)
}

func TestReportHandler_GetCategorySpending(t *testing.T) {
	uc := new(MockReportUsecase)
	h := handler.NewReportHandler(uc)
	uc.On("GetCategorySpending", int64(1), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return([]report.CategoryReport{{CategoryID: 1, CategoryName: "Food", TotalAmount: 150}}, nil).Once()

	r := newRouter()
	r.GET("/reports/spending", h.GetCategorySpending)

	w := doJSON(r, http.MethodGet, "/reports/spending?month=5&year=2026", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Food")
	uc.AssertExpectations(t)
}

func TestReportHandler_GetCategorySpending_Error(t *testing.T) {
	uc := new(MockReportUsecase)
	h := handler.NewReportHandler(uc)
	uc.On("GetCategorySpending", int64(1), mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return([]report.CategoryReport{}, errors.New("db")).Once()

	r := newRouter()
	r.GET("/reports/spending", h.GetCategorySpending)

	w := doJSON(r, http.MethodGet, "/reports/spending", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
