package models

type Reservation struct {
	Id             int    `json:"id"`
	ReservationUid string `json:"reservationUid"`
	Username       string `json:"username"`
	BookUid        string `json:"bookUid"`
	LibraryUid     string `json:"libraryUid"`
	Status         string `json:"status"`
	StartDate      string `json:"startDate"`
	TillDate       string `json:"tillDate"`
}

type ReservationPatch struct {
	Status *string `json:"status,omitempty"`
}

type ReservationInfo struct {
	ReservationUid string        `json:"reservationUid"`
	Status         string        `json:"status"`
	StartDate      string        `json:"startDate"`
	TillDate       string        `json:"tillDate"`
	Book           *ShortBook    `json:"book"`
	Library        *ShortLibrary `json:"library"`
	Rating         *Rating       `json:"rating"`
}

type ReservationResult struct {
	ReservationUid string        `json:"reservationUid"`
	Status         string        `json:"status"`
	StartDate      string        `json:"startDate"`
	TillDate       string        `json:"tillDate"`
	Book           *ShortBook    `json:"book"`
	Library        *ShortLibrary `json:"library"`
	Rating         *Rating       `json:"rating"`
}
