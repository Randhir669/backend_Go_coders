package Controller

import (
	services "My-APP-Go/Services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataController struct {
	dataService *services.DataService
}

func NewDataController() *DataController {
	return &DataController{
		dataService: services.NewDataService(),
	}
}

func (c *DataController) GetAllData(ctx *gin.Context) {

	log.Println("Inside Controller")

	data, err := c.dataService.GetAllData()

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, data)
}

func (c *DataController) registerNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Do something with the data
		log.Println("data", data)
		fmt.Println(data)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data received"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
