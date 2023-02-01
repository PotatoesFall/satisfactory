package game

type Item string

type Recipe struct {
	Name        string       `json:"name"`
	Ingredients map[Item]int `json:"ingredients"`
	Products    map[Item]int `json:"products"`
	Duration    float64      `json:"duration"`
	Machine     string       `json:"machine"`
	Power       float64      `json:"power"`
}

type Info struct {
	Items   []Item    `json:"items"`
	Recipes []*Recipe `json:"recipes"`
}
