package recipes

type Recipe struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
}

type Ingredient struct {
	Name string `json:"name"`
}
