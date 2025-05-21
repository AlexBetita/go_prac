package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/AlexBetita/go_prac/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepo struct {
	users map[string]*models.User
}

// AddProvider implements repositories.UserRepository.
func (m *mockRepo) AddProvider(ctx context.Context, id primitive.ObjectID, provider string) error {
	panic("unimplemented")
}

func (m *mockRepo) Create(ctx context.Context, user *models.User) error {
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = user.CreatedAt
	m.users[user.Email] = user
	return nil
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

func TestRegister(t *testing.T) {
	repo := &mockRepo{}
	svc := services.NewAuthService(repo, "secret")

	user, err := svc.Register(context.Background(), "test@example.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)

	_, err = svc.Register(context.Background(), "test@example.com", "another")
	assert.Error(t, err)
	assert.Equal(t, "email already in use", err.Error())
}

func TestPasswordHash(t *testing.T) {
	hash, _ := utils.HashPassword("password")
	assert.True(t, utils.CheckPasswordHash("password", hash))
	assert.False(t, utils.CheckPasswordHash("wrongpass", hash))
}
