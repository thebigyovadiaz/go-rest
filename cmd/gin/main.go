package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/thebigyovadiaz/go-rest/pkg/recipes"
)

func main() {
	// Create Gin router
	router := gin.Default()

	// Instantiate recipe Handler and provide a data store implementation
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	// Register routes
	router.GET("/", homePage)
	router.POST("/recipes", recipesHandler.CreateRecipe)
	router.GET("/recipes", recipesHandler.ListRecipes)
	router.GET("/recipes/:id", recipesHandler.GetRecipe)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)

	// Start the server
	router.Run(":8020")
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h RecipesHandler) ListRecipes(c *gin.Context) {
	resources, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "resources": resources})
}

func (h RecipesHandler) GetRecipe(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	recipe, err := h.store.Get(id)
	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "recipe": recipe})
}

func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h RecipesHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.store.Remove(id); err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}
