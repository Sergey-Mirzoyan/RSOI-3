package models

type Book struct {
	Id        int    `json:"id"`
	BookUid   string `json:"bookUid"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Genre     string `json:"genre"`
	Condition string `json:"condition"`
}

type BookPatch struct {
	Name      *string `json:"name,omitempty"`
	Author    *string `json:"author,omitempty"`
	Genre     *string `json:"genre,omitempty"`
	Condition *string `json:"condition,omitempty"`
}