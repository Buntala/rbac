package main

import (
	"log"
	"net/http"

	"github.com/Buntala/rbac/controllers"
	"github.com/Buntala/rbac/models"
)

var (
    users_db        []models.User
    role_db         []models.Role
    permission_db   []models.Permission
)

func main(){
    permission_controller := controllers.NewPermissionController(&permission_db)
    role_controller := controllers.NewRoleController(&role_db, &permission_db)
    user_controller := controllers.NewUserController(&users_db, &role_db)

    mux := http.NewServeMux()
    mux.HandleFunc("/permissions", permission_controller.UrlHandler)
    mux.HandleFunc("/permissions/", permission_controller.UrlWithIdHandler)

    mux.HandleFunc("/users", user_controller.UrlHandler)
    mux.HandleFunc("/users/", user_controller.UrlWithIdHandler)

    mux.HandleFunc("/roles", role_controller.UrlHandler)
    mux.HandleFunc("/roles/", role_controller.UrlWithIdHandler)

    log.Println("Server is Running")
    http.ListenAndServe(":8080", mux)
}

