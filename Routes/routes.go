package routes

import (
	"My-APP-Go/Controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	dataController := Controller.NewDataController()
	//	log.Println("Inside Initroutes")

	//	v1 := r.Group("/api/v1")
	//	{
	r.GET("/data", dataController.GetAllData)
	// v1.POST("/data", dataController.CreateData)
	// v1.PUT("/data/:id", dataController.UpdateData)
	// v1.DELETE("/data/:id", dataController.DeleteData)
	//	}
}
