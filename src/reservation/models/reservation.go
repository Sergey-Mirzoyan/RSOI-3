package models

type InputReservation struct {
	Username   string `json:"username"`
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

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
