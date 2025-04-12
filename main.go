package main

import (
	r "chat-app/routes"

	"github.com/gorilla/mux"
)

func main(){
	router := mux.NewRouter()
	newRouter := r.NewRouter(router)
	newRouter.RegisterRoute()
	newRouter.Run(":8080",router)
	


}