package routes

import (
	middleware "chat-app/Auth"
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
	
	_ : middleware.RoleMiddleware("ADMIN")
	_ = middleware.RoleMiddleware("USER")
	both := middleware.RoleMiddleware("ADMIN","USER")

	userRepo := userrepository.NewUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	baseRoutes := r.route.PathPrefix("/api/v1").Subrouter()
	
	baseRoutes.Handle("/login",http.HandlerFunc(userHandler.Login)).Methods("POST")
	baseRoutes.Handle("/create-user",http.HandlerFunc(userHandler.CreateUser)).Methods("POST")
	baseRoutes.Handle("/get-user-by-id/{id}",both(http.HandlerFunc(userHandler.GetUserByID))).Methods("GET")
	baseRoutes.Handle("/get-user-by-phone_number/{phoneNumber}",both(http.HandlerFunc(userHandler.GetUserByPhoneNumber))).Methods("GET")
	baseRoutes.Handle("/update-user/{id}",both(http.HandlerFunc(userHandler.UpdateUser))).Methods("PUT")
	baseRoutes.Handle("/delete-user/{id}",both(http.HandlerFunc(userHandler.DeleteUser))).Methods("DELETE")
}
func (r *Router) Run(addr string, router *mux.Router) error {
	log.Println("Server running on port: ", addr)
	return http.ListenAndServe(addr, router)
}