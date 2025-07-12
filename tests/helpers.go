package tests

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-users-api/models"
	"go-users-api/repository"
	"go-users-api/services"
)

// setupTestRouter crea un router de prueba configurado
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// createTestUser crea un usuario de prueba
func createTestUser() *models.User {
	return &models.User{
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		Phone:     "+1234567890",
		Address:   "123 Test St, Test City",
		UUID:      "test-uuid-123",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// createTestUserRequest crea una request de prueba
func createTestUserRequest() models.CreateUserRequest {
	return models.CreateUserRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Age:     25,
		Phone:   "+1234567890",
		Address: "123 Test St, Test City",
	}
}

// createTestUpdateRequest crea una request de actualización de prueba
func createTestUpdateRequest() models.UpdateUserRequest {
	return models.UpdateUserRequest{
		Name:    "Updated Test User",
		Email:   "updated@example.com",
		Age:     26,
		Phone:   "+0987654321",
		Address: "456 Updated St, Updated City",
	}
}

// MockUserRepository implementa la interfaz UserRepositoryInterface para testing
type MockUserRepository struct {
	users  map[string]*models.User
	emails map[string]bool
}

func NewMockUserRepository() repository.UserRepositoryInterface {
	return &MockUserRepository{
		users:  make(map[string]*models.User),
		emails: make(map[string]bool),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.UUID] = user
	m.emails[user.Email] = true
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	for _, user := range m.users {
		if user.UUID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetAll(ctx context.Context, page, limit int64) ([]models.User, int64, error) {
	var users []models.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, int64(len(users)), nil
}

func (m *MockUserRepository) Update(ctx context.Context, id string, user *models.User) error {
	for uuid := range m.users {
		if uuid == id {
			m.users[uuid] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	for uuid := range m.users {
		if uuid == id {
			delete(m.users, uuid)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return m.emails[email], nil
}

func (m *MockUserRepository) GetByUUID(ctx context.Context, uuid string) (*models.User, error) {
	if user, exists := m.users[uuid]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

// MockUserService implementa la interfaz UserServiceInterface para testing
type MockUserService struct {
	users map[string]*models.User
}

func NewMockUserService() services.UserServiceInterface {
	return &MockUserService{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	user := models.NewUser(req)
	m.users[user.UUID] = user
	return user, nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, assert.AnError
}

func (m *MockUserService) GetUsers(ctx context.Context, pageStr, limitStr string) (*models.UsersResponse, error) {
	var users []models.UserResponse
	for _, user := range m.users {
		users = append(users, user.ToResponse())
	}
	return &models.UsersResponse{
		Users: users,
		Total: int64(len(users)),
	}, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		user.Update(req)
		return user, nil
	}
	return nil, assert.AnError
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return assert.AnError
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, assert.AnError
}

func (m *MockUserService) ValidateUserData(req models.CreateUserRequest) error {
	if req.Name == "" {
		return assert.AnError
	}
	if req.Email == "" {
		return assert.AnError
	}
	if req.Age < 1 || req.Age > 120 {
		return assert.AnError
	}
	return nil
}
