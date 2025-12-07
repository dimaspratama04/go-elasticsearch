package entity

type Movies struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Director string `json:"director"`
	Released int    `json:"released"`
}
