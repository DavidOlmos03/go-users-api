package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-users-api/models"
	"go-users-api/services"
)

// UserController maneja las peticiones HTTP para usuarios
type UserController struct {
	userService *services.UserService
}

// NewUserController crea una nueva instancia del controlador de usuarios
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} models.SuccessResponse{data=models.UserResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest

	// Validar datos de entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validar datos del usuario
	if err := c.userService.ValidateUserData(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Crear usuario
	user, err := c.userService.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}

		ctx.JSON(status, models.ErrorResponse{
			Error:   "Error creating user",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created successfully",
		Data:    user.ToResponse(),
	})
}

// GetUsers godoc
// @Summary Obtener lista de usuarios
// @Description Obtiene la lista paginada de todos los usuarios
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Número de página (default: 1)"
// @Param limit query int false "Límite de elementos por página (default: 10)"
// @Success 200 {object} models.UsersResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	// Obtener parámetros de paginación
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	// Obtener usuarios
	users, err := c.userService.GetUsers(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Error getting users",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Obtener usuario por ID
// @Description Obtiene un usuario específico por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} models.SuccessResponse{data=models.UserResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// Obtener usuario
	user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid user ID" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, models.ErrorResponse{
			Error:   "Error getting user",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User retrieved successfully",
		Data:    user.ToResponse(),
	})
}

// UpdateUser godoc
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param user body models.UpdateUserRequest true "Datos actualizados del usuario"
// @Success 200 {object} models.SuccessResponse{data=models.UserResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateUserRequest

	// Validar datos de entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Actualizar usuario
	user, err := c.userService.UpdateUser(ctx.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "user not found":
			status = http.StatusNotFound
		case "invalid user ID":
			status = http.StatusBadRequest
		case "email already exists":
			status = http.StatusConflict
		}

		ctx.JSON(status, models.ErrorResponse{
			Error:   "Error updating user",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User updated successfully",
		Data:    user.ToResponse(),
	})
}

// DeleteUser godoc
// @Summary Eliminar usuario
// @Description Elimina un usuario del sistema
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	// Eliminar usuario
	err := c.userService.DeleteUser(ctx.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid user ID" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, models.ErrorResponse{
			Error:   "Error deleting user",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User deleted successfully",
	})
} 

