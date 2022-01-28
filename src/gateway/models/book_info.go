package models

type BookReservationInfo struct {
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

type BookReturningInfo struct {
	Condition string `json:"condition"`
	Date      string `json:"date"`
}

type BookAmountInfo struct {
	BookUid        string `json:"bookUid"`
	LibraryUid     string `json:"libraryUid"`
	AvailableCount int    `json:"availableCount"`
}
