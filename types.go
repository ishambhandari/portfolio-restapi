package main

type Work struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	Code_link   string `json:"code_link"`
	Live_link   string `json:"live_link"`
}

func NewWork(title string, description string, imageurl string, code_link string, live_link string) *Work {
	return &Work{
		Title:       title,
		Description: description,
		ImageUrl:    imageurl,
		Code_link:   code_link,
		Live_link:   live_link,
	}
}
