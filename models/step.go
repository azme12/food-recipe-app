package models

type Step struct {
    ID          int    `json:"id"`
    RecipeID    int    `json:"recipe_id"`
    StepNumber  int    `json:"step_number"`
    Description string `json:"description"`
}
