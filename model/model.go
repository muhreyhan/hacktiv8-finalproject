package model

type ToDo struct {
	ID             int    `json:"ID"`
	Title          string `json:"Title"`
	Desc           string `json:"Desc"`
	DueDate        string `json:"DueDate"`
	PersonInCharge int    `json:"PersonInCharge"`
	Status         int    `json:"Status"`
}

type Status struct {
	StatusID  int    `json:"StatusID"`
	StatusTxt string `json:"StatusTxt"`
}

type User struct {
	UserID int    `json:"userID"`
	Name   string `json:"Name"`
}
