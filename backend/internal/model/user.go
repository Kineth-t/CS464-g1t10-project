package model

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    State   string `json:"state"`
    Country string `json:"country"`
    ZipCode string `json:"zip_code"`
}

type User struct {
    ID          int     `json:"id"`
    Username    string  `json:"username"`
    Password    string  `json:"-"`
    PhoneNumber string  `json:"phone_number"`
    Address     Address `json:"address"`
}