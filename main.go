package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

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
	IsStarred    bool   `json:"isstarred"`
}

type docDetails struct {
	Fid          string `json:"fid"`
	Username     string `json:"username"`
	Filename     string `json:"filename"`
	Docstatus    bool   `json:"docstatus"`
	SaveDateTime string `json:"savedatetime"`
	SendBy       string `json:"sendby"`
	IsStarred    bool   `json:"isstarred"`
}

const (
	host     = "database-go-coders.c6xuaoufhj3r.ap-south-1.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "admin123"
	dbname   = "postgres"
)

var users = []userDetails{}
var ErrCiphertextTooShort = errors.New("ciphertext too short")

func main() {

	router := gin.Default()

	// Initialize routes

	router.Use(cors.Default())

	//dataController := Controller.NewDataController()

	//usersAPI
	router.GET("/userdata", GetAllData)
	router.POST("/data", AddUsers)
	router.GET("/data/:id", GetDataById)
	router.PUT("/resetpassword/:id", ResetPassword)

	//Notesmanager
	router.GET("/filedata/:id", GetFileDataById)
	router.PUT("/fileupdate/:id", UpdateFileById)
	router.PUT("/isstarred/:id", UpdatestarredById)
	router.POST("/filedata", AddFile)

	//FileManager

	router.POST("/fileupload", AddUploadedFile)
	router.GET("/filedownload/:filename", downloadFile)
	router.GET("/alldoc/:id", GetAllDoc)
	router.GET("/allhistorydoc/:id", GetAllHistoryDoc)
	router.PUT("/deletefile/:id", DeleteFile)
	router.PUT("/starredfile/:id", starredfile)

	router.Run("localhost:8000")
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

	sqlStatement := `INSERT INTO "TextEditorFileContent" (uuid,userid,filename,filecontent,savedate,filestatus,username,savedatetime,isstarred)
	VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9)`
	_, err = db.Exec(sqlStatement, newfile.Uuid, newfile.Userid, newfile.FileName, newfile.FileContent, newfile.SaveDate, newfile.FileStatus, newfile.UserName, newfile.SaveDateTime,newfile.IsStarred)

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
	rows, err := db.Query(`SELECT name,id,email FROM "UserDetailsTable"`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var userDetailsList []userDetails

	for rows.Next() {
		var userDetails userDetails
		err := rows.Scan(&userDetails.Name, &userDetails.Id, &userDetails.Email)
		if err != nil {
			panic(err)
		}
		userDetailsList = append(userDetailsList, userDetails)
	}

	c.IndentedJSON(http.StatusOK, userDetailsList)

}

func ResetPassword(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	var userDetailsList userDetails

	if err := c.BindJSON(&userDetailsList); err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := ID
	fmt.Println("userDetailsList", userDetailsList.Password, IDStr)
	sqlStatement := `UPDATE "UserDetailsTable" SET  "password"=$1 WHERE "id"=$2`
	_, err = db.Exec(sqlStatement, userDetailsList.Password, IDStr)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, userDetailsList)
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
		err := rows.Scan(&filesDetails.Uuid, &filesDetails.Userid, &filesDetails.FileName, &filesDetails.FileContent, &filesDetails.SaveDate, &filesDetails.FileStatus, &filesDetails.UserName, &filesDetails.SaveDateTime,&filesDetails.IsStarred)
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

func UpdatestarredById(c *gin.Context) {

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
	sqlStatement := `UPDATE "TextEditorFileContent" SET   "isstarred"=$1, "username"=$2 WHERE "uuid"=$3`
	_, err = db.Exec(sqlStatement,  updatedfile.IsStarred, updatedfile.UserName, IDStr)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, updatedfile)
}

//FileManager

func AddUploadedFile(c *gin.Context) {

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to parse multipart form")
		return
	}

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
	tempFile, err := ioutil.TempFile("", "file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println("Error Retrieving the File")
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
	currentTime := time.Now().UTC().Format("20060102150405") // Format: YYYYMMDDHHmmss
	timestamp, err := strconv.Atoi(currentTime)
	fid := timestamp
	username := c.PostForm("username")
	savedatetime := c.PostForm("savedatetime")
	sendby := c.PostForm("sendBy")

	_, err = db.Exec(`INSERT INTO "UploadedFile" (fid,filename, filepath, filesize, mimetype,username,sendby,savedatetime,filecontent) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,E'\\x' || $9::bytea)`,
		fid, handler.Filename, tempFile.Name(), len(fileBytes), handler.Header.Get("Content-Type"), username, sendby, savedatetime, fileBytes)
	if err != nil {
		log.Fatal(err)
	}
	// return a success message to the client
	c.String(http.StatusOK, "Successfully Uploaded File")

}

func downloadFile(c *gin.Context) {
	filename := c.Param("filename")

	// open a database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// retrieve the file contents from the database
	var fileContent []byte

	err = db.QueryRow(`SELECT filecontent FROM "UploadedFile" WHERE filename=$1`, filename).Scan(&fileContent)

	if err != nil {
		log.Fatal(err)
	}

	// set the response headers
	c.Header("Content-Description", "Document")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", http.DetectContentType(fileContent))
	c.Header("Content-Length", strconv.Itoa(len(fileContent)))

	// write the file contents to the response
	c.Data(http.StatusOK, http.DetectContentType(fileContent), fileContent)
}

func GetAllDoc(c *gin.Context) {

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

	rows, err := db.Query(`SELECT fid,username,filename,docstatus,savedatetime,isstarred ,sendby FROM "UploadedFile" WHERE username = ` + IDStr)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var docDetailsList []docDetails

	for rows.Next() {
		var docDetails docDetails
		err := rows.Scan(&docDetails.Fid, &docDetails.Username, &docDetails.Filename, &docDetails.Docstatus, &docDetails.SaveDateTime,&docDetails.IsStarred, &docDetails.SendBy)
		if err != nil {
			log.Println(err)
		}
		if docDetails.Docstatus == true {
			docDetailsList = append(docDetailsList, docDetails)
		}

	}

	c.IndentedJSON(http.StatusOK, docDetailsList)

}

func GetAllHistoryDoc(c *gin.Context) {

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

	rows, err := db.Query(`SELECT fid,username,filename,docstatus,savedatetime ,sendby FROM "UploadedFile" WHERE sendby = ` + IDStr)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var docDetailsList []docDetails

	for rows.Next() {
		var docDetails docDetails
		err := rows.Scan(&docDetails.Fid, &docDetails.Username, &docDetails.Filename, &docDetails.Docstatus, &docDetails.SaveDateTime, &docDetails.SendBy)
		if err != nil {
			log.Println(err)
		}
		if docDetails.Docstatus == true {
			docDetailsList = append(docDetailsList, docDetails)
		}

	}

	c.IndentedJSON(http.StatusOK, docDetailsList)

}

func DeleteFile(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	var updateddoc docDetails

	if err := c.BindJSON(&updateddoc); err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := ID
	fmt.Println("updatedfile", updateddoc.Docstatus, updateddoc.SaveDateTime, IDStr)
	sqlStatement := `UPDATE "UploadedFile" SET  "docstatus"=$1, "username"=$2, "savedatetime"=$3 WHERE "fid"=$4`
	_, err = db.Exec(sqlStatement, updateddoc.Docstatus, updateddoc.Username, updateddoc.SaveDateTime, IDStr)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, updateddoc)

}

func starredfile(c *gin.Context) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	var updateddoc docDetails

	if err := c.BindJSON(&updateddoc); err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	ID := c.Param("id")
	log.Println("ID", ID)

	// Return the retrieved data in the HTTP response
	IDStr := ID
	sqlStatement := `UPDATE "UploadedFile" SET  "isstarred"=$1, "username"=$2 WHERE "fid"=$3`
	_, err = db.Exec(sqlStatement, updateddoc.IsStarred, updateddoc.Username, IDStr)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, updateddoc)

}
