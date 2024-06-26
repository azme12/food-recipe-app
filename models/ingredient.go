package models

type Ingredient struct {
    ID       int    `json:"id"`
    RecipeID int    `json:"recipe_id"`
    Name     string `json:"name"`
    Quantity string `json:"quantity"`
}
