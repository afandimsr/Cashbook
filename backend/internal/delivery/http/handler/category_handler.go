package handler

import (
	"net/http"
	"strconv"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/category"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase uc.Usecase
}

func NewCategoryHandler(usecase uc.Usecase) *CategoryHandler {
	return &CategoryHandler{
		usecase: usecase,
	}
}

// GetCategories godoc
// @Summary      List financial categories
// @Description  Retrieve all personalized spending and income categories available for the user.
// @Tags         Categories
// @Produce      json
// @Success      200 {object} response.SuccessCategoryResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	categories, err := h.usecase.GetAllByUserID(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", categories)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Define a new classification category with custom styling (color/icon) for transaction organization.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        body body category.Category true "Category payload"
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req category.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	req.UserID = c.MustGet("user_id").(int64)

	if err := h.usecase.Create(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "category created", nil)
}

// UpdateCategory godoc
// @Summary      Modify category details
// @Description  Update the properties of an existing category, such as name or visual identifiers.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Param        body body category.Category true "Category payload"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid id", err))
		return
	}

	var req category.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	if err := h.usecase.Update(id, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "category updated", nil)
}

// DeleteCategory godoc
// @Summary      Archive category
// @Description  Permanently remove a transaction category from the user's profile.
// @Tags         Categories
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid id", err))
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "category deleted", nil)
}
