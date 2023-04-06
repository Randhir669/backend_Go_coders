package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type userDetails struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Id       string `json:"id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Country  string `json:"country"`
	Gender   string `json:"gender"`
}

type filesDetails struct {
	Uuid         int    `json:"uuid"`
	Userid       int    `json:"userid"`
	FileName     string `json:"filename"`
	FileContent  string `json:"filecontent"`
	SaveDate     string `json:"savedate"`
	FileStatus   bool   `json:"filestatus"`
	UserName     string `json:"username"`
	SaveDateTime string `json:"savedatetime"`
}

const (
	host     = "database-2.cwlf1t1daphi.eu-north-1.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "admin123"
	dbname   = "sampledb"
)

var users = []userDetails{}

func main() {

	router := gin.Default()

	// Initialize routes

	router.Use(cors.Default())

	//dataController := Controller.NewDataController()
	router.GET("/data", GetAllData)
	router.POST("/data", AddUsers)
	router.GET("/data/:id", GetDataById)
	router.GET("/filedata/:id", GetFileDataById)
	router.PUT("/fileupdate/:id", UpdateFileById)
	router.POST("/filedata", AddFile)
	router.POST("/fileupload", AddUploadedFile)
	router.Run("localhost:8000")

	fmt.Printf("Hello")

}

func AddFile(c *gin.Context) {

	var newfile filesDetails
	// Call BindJSON to bind the received JSON to

	if err := c.BindJSON(&newfile); err != nil {
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	fmt.Println("newfile.Uuid  newfile.Userid newfile.Userid", newfile.Uuid, newfile.Userid, newfile.Userid)

	sqlStatement := `INSERT INTO "TextEditorFileContent" (uuid,userid,filename,filecontent,savedate,filestatus,username,savedatetime)
	VALUES($1, $2,$3,$4,$5,$6,$7,$8)`
	_, err = db.Exec(sqlStatement, newfile.Uuid, newfile.Userid, newfile.FileName, newfile.FileContent, newfile.SaveDate, newfile.FileStatus, newfile.UserName, newfile.SaveDateTime)

	if err != nil {
		panic(err)
	}
	//	users = append(users, newfile)
	c.IndentedJSON(http.StatusCreated, newfile)

}

func AddUsers(c *gin.Context) {

	var newuser userDetails
	// Call BindJSON to bind the received JSON to

	if err := c.BindJSON(&newuser); err != nil {
		return
	}
	log.Println("newuser", newuser.Name)
	log.Println("newuser", newuser.Password)
	dbConnection(newuser, "post")
	users = append(users, newuser)
	c.IndentedJSON(http.StatusCreated, newuser)

}

func GetDataById(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := "'" + ID + "'"

	rows, err := db.Query(`SELECT * FROM "UserDetailsTable" WHERE Id = ` + IDStr)
	log.Println("rows", rows)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var userDetailsList []userDetails

	for rows.Next() {
		var userDetails userDetails
		err := rows.Scan(&userDetails.Name, &userDetails.Password, &userDetails.Id, &userDetails.Phone, &userDetails.Email, &userDetails.Address, &userDetails.Country, &userDetails.Gender)
		if err != nil {
			log.Println(err)
		}
		userDetailsList = append(userDetailsList, userDetails)
	}

	c.IndentedJSON(http.StatusOK, userDetailsList)

}

func GetAllData(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()

	//db := Database.dbConnection()
	rows, err := db.Query(`SELECT * FROM "UserDetailsTable"`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var userDetailsList []userDetails

	for rows.Next() {
		var userDetails userDetails
		err := rows.Scan(&userDetails.Name, &userDetails.Password, &userDetails.Id, &userDetails.Phone, &userDetails.Email, &userDetails.Address, &userDetails.Country, &userDetails.Gender)
		if err != nil {
			panic(err)
		}
		userDetailsList = append(userDetailsList, userDetails)
	}

	c.IndentedJSON(http.StatusOK, userDetailsList)

}

func dbConnection(newuser userDetails, requestTypes string) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	if requestTypes == "post" {

		sqlStatement := `INSERT INTO "UserDetailsTable" (id,name,password,address,email,gender,phone,country)
	VALUES($1, $2,$3,$4,$5,$6,$7,$8)`
		_, err = db.Exec(sqlStatement, newuser.Id, newuser.Name, newuser.Password, newuser.Address, newuser.Email, newuser.Gender, newuser.Phone, newuser.Country)

	}

	if err != nil {
		panic(err)
	}
	/*	sqlStatement := `INSERT INTO "UserDetailsTable" (name,password)`
		VALUES(newuser.Name, newuser.Password)
		_, err = db.Exec(sqlStatement)*/

}

func GetFileDataById(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := "'" + ID + "'"

	rows, err := db.Query(`SELECT * FROM "TextEditorFileContent" WHERE username = ` + IDStr)
	log.Println("rows", rows)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var fileDetailsList []filesDetails

	for rows.Next() {
		var filesDetails filesDetails
		err := rows.Scan(&filesDetails.Uuid, &filesDetails.Userid, &filesDetails.FileName, &filesDetails.FileContent, &filesDetails.SaveDate, &filesDetails.FileStatus, &filesDetails.UserName, &filesDetails.SaveDateTime)
		if err != nil {
			log.Println(err)
		}
		if filesDetails.FileStatus == true {
			fileDetailsList = append(fileDetailsList, filesDetails)
		}

	}

	c.IndentedJSON(http.StatusOK, fileDetailsList)

}

func UpdateFileById(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	var updatedfile filesDetails

	if err := c.BindJSON(&updatedfile); err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := ID
	fmt.Println("updatedfile", updatedfile.FileContent, updatedfile.SaveDateTime, IDStr)
	sqlStatement := `UPDATE "TextEditorFileContent" SET  "filecontent"=$1, "filestatus"=$2, "username"=$3, "savedatetime"=$4 WHERE "uuid"=$5`
	_, err = db.Exec(sqlStatement, updatedfile.FileContent, updatedfile.FileStatus, updatedfile.UserName, updatedfile.SaveDateTime, IDStr)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, updatedfile)
}

func AddUploadedFile(c *gin.Context) {

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to parse multipart form")
		return
	}

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.String(http.StatusBadRequest, "Failed to retrieve file from form data")
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//
	tempFile, err := ioutil.TempFile("", "Uploadedfile-")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	fmt.Println(tempFile)
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to read file contents")
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	c.String(http.StatusOK, "Successfully Uploaded File")
	// open a database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// insert a record in the database with the file's metadata
	fid := 7899077
	_, err = db.Exec("INSERT INTO UploadedFile (fid,filename, filepath, filesize, mimetype) VALUES ($1, $2, $3, $4,$5)",
		fid, handler.Filename, tempFile.Name(), len(fileBytes), handler.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	// return a success message to the client
	c.String(http.StatusOK, "Successfully Uploaded File")

}
