package models

type LibrariesPage struct {
	Page          int        `json:"page"`
	PageSize      int        `json:"pageSize"`
	TotalElements int        `json:"totalElements"`
	Items         *[]Library `json:"items"`
}

type BooksPage struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"pageSize"`
	TotalElements int     `json:"totalElements"`
	Items         *[]Book `json:"items"`
}
