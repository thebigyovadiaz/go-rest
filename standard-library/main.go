package main

import (
	"net/http"
	"regexp"

	"github.com/thebigyovadiaz/go-rest/standard-library/recipes"
)

func main() {
	// Create Store and Recipe Handler
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)
	// Create multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipes", recipesHandler)
	mux.Handle("/recipes/", recipesHandler)

	// Run the serve
	http.ListenAndServe(":8080", mux)
}

var (
	RecipeRe       = regexp.MustCompile(`^/recipes/*$`)
	RecipeReWithID = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

// recipeStore represent a data storage containing recipes
type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

// RecipesHandler implements http.Handler and dispatch request to the store
type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my create recipe page"))
}

func (h *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my list recipes page"))
}

func (h *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my get recipe page"))
}

func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my update recipe page"))
}

func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my delete recipe page"))
}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && RecipeReWithID.MatchString(r.URL.Path):
		h.GetRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipes(w, r)
		return
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipe(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithID.MatchString(r.URL.Path):
		h.UpdateRecipe(w, r)
		return
	case r.Method == http.MethodDelete && RecipeRe.MatchString(r.URL.Path):
		h.DeleteRecipe(w, r)
		return
	default:
		return
	}
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
