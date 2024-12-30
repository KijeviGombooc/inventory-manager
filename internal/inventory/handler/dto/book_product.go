package dto

type BookProduct struct {
	Product `json:",inline"`
	Author  string `json:"author"`
}
