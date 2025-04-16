package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockPostRepo struct {
	posts map[primitive.ObjectID]*models.Post
}

func (m *mockPostRepo) Create(ctx context.Context, post *models.Post) error {
	if m.posts == nil {
		m.posts = make(map[primitive.ObjectID]*models.Post)
	}
	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now().Unix()
	post.UpdatedAt = post.CreatedAt
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
	if post, ok := m.posts[id]; ok {
		return post, nil
	}
	return nil, errors.New("not found")
}

func TestCreatePost(t *testing.T) {
	repo := &mockPostRepo{}

	post := &models.Post{
		UserID:     primitive.NewObjectID(),
		Topic:      "Unit Testing in Go",
		Content:    "Here's why unit tests are powerful.",
		Summary:    "Testing Go applications",
		Message:    "make a post about testing",
		Slug:       "unit-testing-in-go",
		Tags:       []string{"Tech", "Education"},
		Keywords:   []string{"testing", "go", "unittest"},
		CreatedBy:  "Tester",
	}

	err := repo.Create(context.Background(), post)

	assert.NoError(t, err)
	assert.NotEmpty(t, post.ID)
	assert.Equal(t, "Unit Testing in Go", post.Topic)
}

func TestFindPostByID(t *testing.T) {
	repo := &mockPostRepo{}
	svc := services.NewPostService(repo, "test-jwt-secret")

	post := &models.Post{
		UserID:  primitive.NewObjectID(),
		Topic:   "Testing FindByID",
		Content: "We want to ensure this is retrievable.",
		Slug:    "testing-findbyid",
	}
	err := repo.Create(context.Background(), post)
	assert.NoError(t, err)

	found, err := svc.GetPostsByID(context.Background(), post.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, post.ID, found.ID)

	_, err = svc.GetPostsByID(context.Background(), primitive.NewObjectID().Hex())
	assert.Error(t, err)
}
