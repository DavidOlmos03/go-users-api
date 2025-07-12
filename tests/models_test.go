package tests

import (
	"testing"
	"time"

	"go-users-api/models"
)

func TestNewUser(t *testing.T) {
	req := models.CreateUserRequest{
		Name:    "John Doe",
		Email:   "john.doe@example.com",
		Age:     30,
		Phone:   "+1234567890",
		Address: "123 Main St, City, Country",
	}

	user := models.NewUser(req)

	// Verificar que los campos se asignaron correctamente
	if user.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, user.Name)
	}

	if user.Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, user.Email)
	}

	if user.Age != req.Age {
		t.Errorf("Expected Age %d, got %d", req.Age, user.Age)
	}

	if user.Phone != req.Phone {
		t.Errorf("Expected Phone %s, got %s", req.Phone, user.Phone)
	}

	if user.Address != req.Address {
		t.Errorf("Expected Address %s, got %s", req.Address, user.Address)
	}

	// Verificar que UUID se generó
	if user.UUID == "" {
		t.Error("Expected UUID to be generated")
	}

	// Verificar que timestamps se establecieron
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// Verificar que CreatedAt y UpdatedAt son iguales al crear
	if !user.CreatedAt.Equal(user.UpdatedAt) {
		t.Error("Expected CreatedAt and UpdatedAt to be equal when creating")
	}
}

func TestUserToResponse(t *testing.T) {
	user := &models.User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Phone:     "+1234567890",
		Address:   "123 Main St, City, Country",
		UUID:      "550e8400-e29b-41d4-a716-446655440000",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	response := user.ToResponse()

	// Verificar que la conversión mantiene los datos
	if response.Name != user.Name {
		t.Errorf("Expected Name %s, got %s", user.Name, response.Name)
	}

	if response.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, response.Email)
	}

	if response.Age != user.Age {
		t.Errorf("Expected Age %d, got %d", user.Age, response.Age)
	}

	if response.Phone != user.Phone {
		t.Errorf("Expected Phone %s, got %s", user.Phone, response.Phone)
	}

	if response.Address != user.Address {
		t.Errorf("Expected Address %s, got %s", user.Address, response.Address)
	}

	if response.UUID != user.UUID {
		t.Errorf("Expected UUID %s, got %s", user.UUID, response.UUID)
	}

	if !response.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", user.CreatedAt, response.CreatedAt)
	}

	if !response.UpdatedAt.Equal(user.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", user.UpdatedAt, response.UpdatedAt)
	}
}

func TestUserUpdate(t *testing.T) {
	user := &models.User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Phone:     "+1234567890",
		Address:   "123 Main St, City, Country",
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	originalUpdatedAt := user.UpdatedAt

	req := models.UpdateUserRequest{
		Name:    "John Updated",
		Email:   "john.updated@example.com",
		Age:     31,
		Phone:   "+0987654321",
		Address: "456 Oak Ave, Town, State",
	}

	user.Update(req)

	// Verificar que los campos se actualizaron
	if user.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, user.Name)
	}

	if user.Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, user.Email)
	}

	if user.Age != req.Age {
		t.Errorf("Expected Age %d, got %d", req.Age, user.Age)
	}

	if user.Phone != req.Phone {
		t.Errorf("Expected Phone %s, got %s", req.Phone, user.Phone)
	}

	if user.Address != req.Address {
		t.Errorf("Expected Address %s, got %s", req.Address, user.Address)
	}

	// Verificar que UpdatedAt se actualizó
	if user.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestUserUpdatePartial(t *testing.T) {
	user := &models.User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Phone:     "+1234567890",
		Address:   "123 Main St, City, Country",
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	originalEmail := user.Email
	originalAge := user.Age
	originalPhone := user.Phone
	originalAddress := user.Address

	req := models.UpdateUserRequest{
		Name: "John Updated",
		// Solo actualizar el nombre, otros campos vacíos
	}

	user.Update(req)

	// Verificar que solo se actualizó el nombre
	if user.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, user.Name)
	}

	if user.Email != originalEmail {
		t.Errorf("Expected Email to remain %s, got %s", originalEmail, user.Email)
	}

	if user.Age != originalAge {
		t.Errorf("Expected Age to remain %d, got %d", originalAge, user.Age)
	}

	if user.Phone != originalPhone {
		t.Errorf("Expected Phone to remain %s, got %s", originalPhone, user.Phone)
	}

	if user.Address != originalAddress {
		t.Errorf("Expected Address to remain %s, got %s", originalAddress, user.Address)
	}

	// Verificar que UpdatedAt se actualizó
	if user.UpdatedAt.Equal(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected UpdatedAt to be updated")
	}
}
