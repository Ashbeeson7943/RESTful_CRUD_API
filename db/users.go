package database

type User struct {
	USERNAME string
	PASSWORD string
}

var Users = []User{
	{USERNAME: "test_1", PASSWORD: "1234"},
	{USERNAME: "Mickey", PASSWORD: "1234"},
	{USERNAME: "Minnie", PASSWORD: "1234"},
}
