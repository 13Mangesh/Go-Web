package main

import (
	"fmt"
	"net/http"
	"./models"
	"./utils"
	"./routes"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	fmt.Println("Starting the server at port 8080....")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
