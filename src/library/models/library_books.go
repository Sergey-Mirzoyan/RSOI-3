package models

type LibraryBooks struct {
	BookUid        string `json:"bookUid"`
	LibraryUid     string `json:"libraryUid"`
	AvailableCount int    `json:"availableCount"`
}
