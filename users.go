package main

import "log"

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

var users []User

func InitUsers() {

	// Create users, for now use in memory list, fixed and same for all instances
	users = append(users, User{Username: "User1", Password: "Pass1"})
	users = append(users, User{Username: "User2", Password: "Pass2"})
	users = append(users, User{Username: "User3", Password: "Pass3"})
	users = append(users, User{Username: "User4", Password: "Pass4"})
	log.Printf("%+v\n", users)

	// Each user's data stored in separate SimpleDB domain, create those
	for _, item := range users {
		CreateDomainSimpleDB(item.Username)
	}

	ListDomainsSimpleDB()
}

func ValidateUser(username, password string) string {

	for _, item := range users {
		if item.Username == username {
			if item.Password == password {
				return username
			}
			log.Println("Password mismatch for ", username)
			return ""
		}
	}

	log.Println("No such user", username)
	return ""
}
