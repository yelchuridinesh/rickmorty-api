package model

type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Origin  struct {
		Name string `json:"name"`
	} `json:"origin"`
}
