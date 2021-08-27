package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/users", getusers)

	e.GET("/users/:id", getUserById)

	e.POST("/users", addUser)

	e.PUT("/users/:id", editUser)

	e.DELETE("/users/:id", deleteUser)

	e.Start(":8000")
}

type User struct {
	Id     		int
	Name   		string	`json:name form:name`
	Email 		string	`json:email form:email`
	Password 	string	`json:password form: password`
}

var users = []User{}

// users = []User{
// 	{Id: 1, Name: "Ali Muhammad", Email: "m.ali@gmail.com", Password: "12wqwq"},
// 	{Id: 2, Name: "Fathima Azzahra", Email: "fathima.azzahra@gmail.com", Password: "11wqwq"},
// 	{Id: 3, Name: "Hatta Ahmad", Email: "hattaahmad@gmail.com", Password: "15tua"}}

func getusers(c echo.Context) error {
	res := filterUser(users, func(a User) bool{ return a.Id != -4})
	return c.JSON(http.StatusOK, res)
}

func getUserById(c echo.Context) error {
	targetId, _ := strconv.Atoi(c.Param("id"))
	id := findUserById(users, targetId)

	if id == -1 {
		return c.JSON(http.StatusBadRequest, struct {
			Status  string
			Message string
		}{Status: "Failed", Message: "Wrong User Id"})
	}
	return c.JSON(http.StatusOK, users[id])
}

func addUser(c echo.Context) error {
	newUser := User{}
	c.Bind(&newUser)

	if len(users) == 0 {
		newUser.Id = 1
	} else {
		newUser.Id = findLastUserId(users) + 1
	}

	users = append(users, newUser)
	
	return c.JSON(http.StatusOK, struct {
		Status string
		Message string
		User User
	}{Status: "succes", Message: "User has been created!", User: newUser})

}

func editUser(c echo.Context) error {
	targetId, _ := strconv.Atoi(c.Param("id"))
	id := findUserById(users, targetId)

	if id == -1 {
		return c.JSON(http.StatusBadRequest, struct {
			Status  string
			Message string
		}{Status: "Failed", Message: "Wrong User Id"})
	}

	oldId := users[id].Id
	targetUser := &users[id]
	c.Bind(targetUser)
	targetUser.Id = oldId

	return c.JSON(http.StatusOK, struct {
		Status string
		Message string
		User User
	}{Status: "succes", Message: "User has been updated!", User: *targetUser})
}

func deleteUser(c echo.Context) error {
	targetId, _ := strconv.Atoi(c.Param("id"))
	id := findUserById(users, targetId)

	if id == -1 {
		return c.JSON(http.StatusBadRequest, struct {
			Status  string
			Message string
		}{Status: "Failed", Message: "Wrong User Id"})
	}
	
	deleted := users[id]
	users[id] = User{Id: -4}

	return c.JSON(http.StatusOK, struct {
		Status string
		Message string
		User User
	}{Status: "succes", Message: "User has been deleted!", User: deleted})

}


func findUserById(arr []User, id int) int {
	res := -1
	for i, v := range arr {
		if v.Id == id {
			res = i
			return res
		}
	}
	return res
}

func filterUser(arr []User, filterFunc func(a User) bool ) []User {
	res := []User{}
	for _,v := range arr {
		if filterFunc(v) == true {
			res = append (res, v)
		}
	}
	return res
}

func findLastUserId(arr []User) int {
	res := 0
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i].Id != -4 {
			res = arr[i].Id
			break
		}
	}
	return res
}
