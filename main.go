package main

import (
	"chrisfoong/chaum-work-management-backend/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run() // listens on 0.0.0.0:8080 by default
}
