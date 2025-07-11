package services

import (
	"context"
	"errors"
	"strconv"

	"go-users-api/models"
	"go-users-api/repository"
)

// UserService maneja la lógica de negocio para usuarios
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	// Verificar si el email ya existe
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Crear nuevo usuario
	user := models.NewUser(req)

	// Guardar en la base de datos
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers obtiene todos los usuarios con paginación
func (s *UserService) GetUsers(ctx context.Context, pageStr, limitStr string) (*models.UsersResponse, error) {
	// Parsear parámetros de paginación
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Obtener usuarios de la base de datos
	users, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// Convertir a respuesta
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return &models.UsersResponse{
		Users: userResponses,
		Total: total,
	}, nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error) {
	// Obtener usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verificar si el nuevo email ya existe (si se está actualizando)
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
	}

	// Actualizar campos del usuario
	user.Update(req)

	// Guardar cambios en la base de datos
	err = s.userRepo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Verificar que el usuario existe
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Eliminar usuario
	return s.userRepo.Delete(ctx, id)
}

// GetUserByEmail obtiene un usuario por su email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateUserData valida los datos del usuario
func (s *UserService) ValidateUserData(req models.CreateUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Age < 1 || req.Age > 120 {
		return errors.New("age must be between 1 and 120")
	}

	return nil
} 

