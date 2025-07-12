package tests

import (
	"context"
	"testing"

	"go-users-api/models"
	"go-users-api/services"
)

func TestValidateUserData(t *testing.T) {
	service := services.NewUserService(nil)

	tests := []struct {
		name    string
		req     models.CreateUserRequest
		wantErr bool
	}{
		{
			name: "Valid user data",
			req: models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			req: models.CreateUserRequest{
				Name:  "",
				Email: "john.doe@example.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Empty email",
			req: models.CreateUserRequest{
				Name:  "John Doe",
				Email: "",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Age too young",
			req: models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Age:   0,
			},
			wantErr: true,
		},
		{
			name: "Age too old",
			req: models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Age:   121,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateUserData(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)

	req := models.CreateUserRequest{
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Age:     30,
		Phone:   "+1234567890",
		Address: "123 Main St, City, Country",
	}

	// Test crear usuario exitoso
	user, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, user.Name)
	}

	if user.Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, user.Email)
	}

	// Test crear usuario con email duplicado
	_, err = service.CreateUser(context.Background(), req)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "email already exists" {
		t.Errorf("Expected error 'email already exists', got %s", err.Error())
	}
}

func TestServiceGetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)

	// Crear un usuario primero
	req := models.CreateUserRequest{
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Age:     30,
		Phone:   "+1234567890",
		Address: "123 Main St, City, Country",
	}

	createdUser, _ := service.CreateUser(context.Background(), req)

	// Test obtener usuario existente
	user, err := service.GetUserByID(context.Background(), createdUser.UUID)
	if err != nil {
		t.Errorf("GetUserByID() error = %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, user.Name)
	}

	// Test obtener usuario inexistente
	_, err = service.GetUserByID(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected error 'user not found', got %s", err.Error())
	}
}

func TestServiceGetUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)

	// Crear algunos usuarios
	users := []models.CreateUserRequest{
		{Name: "John Doe", Email: "john.doe@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane.smith@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob.johnson@example.com", Age: 35},
	}

	for _, req := range users {
		service.CreateUser(context.Background(), req)
	}

	// Test obtener usuarios
	response, err := service.GetUsers(context.Background(), "1", "10")
	if err != nil {
		t.Errorf("GetUsers() error = %v", err)
	}

	if len(response.Users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response.Users))
	}

	if response.Total != 3 {
		t.Errorf("Expected total 3, got %d", response.Total)
	}
}

func TestServiceUpdateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)

	// Crear un usuario
	req := models.CreateUserRequest{
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Age:     30,
		Phone:   "+1234567890",
		Address: "123 Main St, City, Country",
	}

	createdUser, _ := service.CreateUser(context.Background(), req)

	// Test actualizar usuario
	updateReq := models.UpdateUserRequest{
		Name:    "John Updated",
		Email:   "john.updated@example.com",
		Age:     31,
		Phone:   "+0987654321",
		Address: "456 Oak Ave, Town, State",
	}

	updatedUser, err := service.UpdateUser(context.Background(), createdUser.UUID, updateReq)
	if err != nil {
		t.Errorf("UpdateUser() error = %v", err)
	}

	if updatedUser.Name != updateReq.Name {
		t.Errorf("Expected Name %s, got %s", updateReq.Name, updatedUser.Name)
	}

	if updatedUser.Email != updateReq.Email {
		t.Errorf("Expected Email %s, got %s", updateReq.Email, updatedUser.Email)
	}

	// Test actualizar usuario inexistente
	_, err = service.UpdateUser(context.Background(), "non-existent-id", updateReq)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected error 'user not found', got %s", err.Error())
	}
}

func TestServiceDeleteUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := services.NewUserService(mockRepo)

	// Crear un usuario
	req := models.CreateUserRequest{
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Age:     30,
		Phone:   "+1234567890",
		Address: "123 Main St, City, Country",
	}

	createdUser, _ := service.CreateUser(context.Background(), req)

	// Test eliminar usuario existente
	err := service.DeleteUser(context.Background(), createdUser.UUID)
	if err != nil {
		t.Errorf("DeleteUser() error = %v", err)
	}

	// Verificar que el usuario fue eliminado
	_, err = service.GetUserByID(context.Background(), createdUser.UUID)
	if err == nil {
		t.Error("Expected error for deleted user")
	}

	// Test eliminar usuario inexistente
	err = service.DeleteUser(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected error 'user not found', got %s", err.Error())
	}
}
