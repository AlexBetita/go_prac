package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/routes"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *mockPostRepo) Search(ctx context.Context, query string, limit int64) ([]*models.Post, error) {
	return []*models.Post{}, nil
}

func (m *mockPostRepo) VectorSearch(ctx context.Context, vector []float32, limit int64) ([]*models.Post, error) {
	return []*models.Post{}, nil
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
	svc := services.NewPostService(repo, nil)

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


func TestSearchByIndex_Unit(t *testing.T) {
	repo := &mockPostRepo{}
	svc := services.NewPostService(repo, nil)

	results, err := svc.SearchPosts(context.Background(), "foo", 5)
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearchByVector_Unit(t *testing.T) {
	repo := &mockPostRepo{}
	svc := services.NewPostService(repo, nil)

	results, err := svc.SearchPostsByVector(context.Background(), "foo", 5)
	assert.NoError(t, err)
	assert.Empty(t, results)
}


func setupIntegrationServer(t *testing.T) *httptest.Server {

	mongoURI := os.Getenv("MONGO_URI")
	dbName   := os.Getenv("DB_NAME")
	oaKey    := os.Getenv("OPENAI_API_KEY")

	cfg := &config.Config{
		MongoURI:  mongoURI,
		DBName:    dbName,
		OpenAIKey: oaKey,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("mongo.Connect failed: %v", err)
	}

	oaClient := openai.NewClient(oaKey)
	handler := routes.NewRouter(cfg, client, oaClient)
	return httptest.NewServer(handler)
}

func TestSearchByIndex_Integration(t *testing.T) {
	ts := setupIntegrationServer(t)
	defer ts.Close()

	cases := []struct {
		q        string
		limit    int
		expected int
	}{
		{"Go Will Change Your Life", 5, 5},
		{"duck debugging", 5, 1},
		{"mongodb fridge", 5, 1},
		{"api therapy", 5, 3},
		{"code needs therapy", 5, 4},
		{"frameworks pizza", 5, 1},
		{"ai poetry", 5, 4},
		{"nap before scaling", 5, 1},
		{"graphql adventure", 5, 1},
		{"funny scripting tips", 5, 1},
	}

	for _, tc := range cases {
		url := ts.URL + "/api/posts/search?q=" +
			strings.ReplaceAll(tc.q, " ", "%20") +
			"&limit=" + strconv.Itoa(tc.limit)

		res, err := http.Get(url)
		assert.NoError(t, err)
		defer res.Body.Close()

		var posts []*models.Post
		err = json.NewDecoder(res.Body).Decode(&posts)
		assert.NoError(t, err)
		assert.Len(t, posts, tc.expected, "q=%q", tc.q)
	}
}