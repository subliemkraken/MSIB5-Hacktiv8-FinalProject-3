package main

import (
	"FinalProject3/config"
	"FinalProject3/controller"
	"FinalProject3/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.StartDB()

	//Users
	r.POST("/users/register", controller.RegisterUser)
	r.POST("/users/login", controller.LoginUser)
	r.PUT("/users/update-account", middleware.Authentication(), controller.UpdateUser)
	r.DELETE("/users/delete-account", middleware.Authentication(), controller.DeleteUser)

	//Categories
	r.GET("/categories", middleware.Authentication(), controller.GetAllCategory)
	r.POST("/categories", middleware.Authentication(), middleware.AdminAuthorization(), controller.CreateCategory)
	r.PATCH("/categories/:categoryID", middleware.Authentication(), middleware.AdminAuthorization(), controller.PatchCategory)
	r.DELETE("/categories/:categoryID", middleware.Authentication(), middleware.AdminAuthorization(), controller.DeleteCategory)

	//Tasks
	r.GET("/tasks", middleware.Authentication(), controller.GetAllTask)
	r.POST("/tasks", middleware.Authentication(), controller.CreateTask)
	r.PUT("/tasks/:taskID", middleware.Authentication(), controller.PutTask)
	r.PATCH("/tasks/update-status/:taskID", middleware.Authentication(), controller.PatchStat)
	r.PATCH("/tasks/update-category/:taskID", middleware.Authentication(), controller.PatchCatId)
	r.DELETE("/tasks/:taskID", middleware.Authentication(), controller.DeleteTask)

	r.Run()
}
