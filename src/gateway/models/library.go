package models

type Library struct {
	Id         int    `json:"id"`
	LibraryUid string `json:"libraryUid"`
	Name       string `json:"name"`
	City       string `json:"city"`
	Address    string `json:"address"`
}

type ShortLibrary struct {
	LibraryUid string `json:"libraryUid"`
	Name       string `json:"name"`
	City       string `json:"city"`
	Address    string `json:"address"`
}
