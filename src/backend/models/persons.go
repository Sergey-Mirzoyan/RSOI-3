package models

type Person struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type InputPerson struct {
	Name    *string
	Age     *int
	Address *string
	Work    *string
}
