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

func (r *Recipe) DurationMinutes() float64 {
	return r.Duration / 60
}

func (r *Recipe) RatePerMinute() float64 {
	return 1 / r.DurationMinutes()
}

type Info struct {
	Items   []Item    `json:"items"`
	Recipes []*Recipe `json:"recipes"`
}
