package routes

import (
	userusecase "chat-app/application/usecases"
	database "chat-app/infrastructure/database"
	userrepository "chat-app/infrastructure/repository"
	handler "chat-app/transport"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct{
	route *mux.Router
}

func NewRouter(route *mux.Router) Router{
	return Router{route: route}
}

func (r *Router) RegisterRoute(){
	db,err :=database.NewDatabase()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	// err=database.Migrate(db)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println("Database migrated successfully")
	userRepo := userrepository.NewUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	baseRoutes := r.route.PathPrefix("/api/v1").Subrouter()

	baseRoutes.HandleFunc("/create-user",userHandler.CreateUser).Methods("POST")
	baseRoutes.HandleFunc("/get-user-by-id/{id}",userHandler.GetUserByID).Methods("GET")
	baseRoutes.HandleFunc("/get-user-by-phone_number/{phoneNumber}",userHandler.GetUserByPhoneNumber).Methods("GET")
	baseRoutes.HandleFunc("/update-user/{id}",userHandler.UpdateUser).Methods("PUT")
	baseRoutes.HandleFunc("/delete-user/{id}",userHandler.DeleteUser).Methods("DELETE")
}
func (r *Router) Run(addr string, router *mux.Router) error {
	log.Println("Server running on port: ", addr)
	return http.ListenAndServe(addr, router)
}