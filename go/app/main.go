package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"path"
	"strings"
	"encoding/json"
	"crypto/sha256"
	"strconv"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ImgDir = "images"
)
type Response struct {
	Message    string `json:"message"`
}


// Define structure
type ItemIndex struct {
	Items      []Item `json:"items"`
}
type Item struct {	
	Name       string `json:"name"`
	Category   string `json:"category"`
	Image_name string `json:"image_name"`
}

// GET "/"
func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

var itemindex ItemIndex
var item Item

// GET "/items"
func getItem(c echo.Context) error {
	// Open database
	DBopen, err := sql.Open("sqlite3", "./mercari.sqlite3")
	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer DBopen.Close()
	// Query items
	rows, err := db.Query("SELECT name, category, image_name FROM items")
	var getitem ItemIndex

	// Decode
	if err := json.NewDecoder(file).Decode(&getitem); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer file.Close()

	return c.JSON(http.StatusOK, getitem)
}


// POST "/items"
func addItem(c echo.Context) error {
	// Create or open JSON file
	// file, err := os.OpenFile("items.json", os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	c.Logger().Infof("Error message: %s", err)
	// }
	// defer file.Close()


	// Get form data
	Name     := c.FormValue("name")
	Category := c.FormValue("category")
	Image, err   := c.FormFile("image")
	if err != nil {
		log.Fatal(err)
	}
	
	image_path := fmt.Sprintf("%x", Image)

	// Log
	c.Logger().Infof("Receive item: %s", Name)
	c.Logger().Infof("Receive category: %s", Category)
	c.Logger().Infof("Receive image: %s", Image)

	message := fmt.Sprintf("item received: %s", Name)

	// Hash
	hash := sha256.Sum256([]byte(image_path))
	hash_string := fmt.Sprintf("%x", hash)
	Image_name := hash_string + ".jpg"

	// Save image
	imgfile, err := os.OpenFile("images", os.O_RDWR, 0644)
	if err != nil {
		c.Logger().Infof("Error message: %s", err)
	}

	src, err := Image.Open()
	if err != nil {
		c.Logger().Infof("Error message: %s", err)
	}
	defer src.Close()

	f, err := os.Create("images/" + hash_string + ".jpg")
	if err != nil {
		c.Logger().Infof("Error message: %s", err)
	}
	defer f.Close()
	
	io.Copy(f, src)
	io.Copy(f, imgfile)
	defer imgfile.Close()

	// Add Items
	DB_write(Name, Category, Image_name)
	
	// // Encode JSON
	// encoder := json.NewEncoder(file)
	// if err := encoder.Encode(itemindex); err != nil {
	// 	c.Logger().Infof("Error message: %s", err)
	//  }
	res := Response{Message: message}
	return c.JSON(http.StatusOK, res)
}

//GET "/items/:id"
func showItem(c echo.Context) error {
	// Get id & debug
	id, err := strconv.Atoi(c.Param("id")) 
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	if id == 0 { 
		c.Logger().Infof("Error message: Out of range")
	}

	//open JSON file
	file, err := os.Open("items.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer file.Close()

	showitem := ItemIndex{} 

	// Decode
	if err := json.NewDecoder(file).Decode(&showitem); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer file.Close()
	return c.JSON(http.StatusOK, showitem.Items[id-1])
}

//GET "/image/:imageFilename
func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Infof("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}	

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", root)
	e.GET("/items", getItem)
	e.POST("/items", addItem)
	e.GET("/items/:id", showItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))

}