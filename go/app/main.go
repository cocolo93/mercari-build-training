package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"path"
	"strings"
	"crypto/sha256"
	"strconv"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)
const (
	DB_PATH = "db/mercari.sqlite3"
)
const (
	ImgDir = "images"
)
type Response struct {
	Message    string `json:"message"`
}

// Define structure
type ItemIndex struct {
	Items      []Item
}
type Item struct {	
	Name       string 
	Category   string
	Image_name string
}

// GET "/"
func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

// GET "/items"
func getItem(c echo.Context) error {
	// Open DB
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer db.Close()
	

	// レコード読み込み
	var items []Item
	rows, err := db.Query("SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer rows.Close()

	// レコード取り出し
	for rows.Next() {
		var add Item
        err := rows.Scan(&add.Name, &add.Category, &add.Image_name)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
        }
		items = append(items, add)
    }

	return c.JSON(http.StatusOK, items)
}


// POST "/items"
func addItem(c echo.Context) error {
	// Open DB
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer db.Close()

	var itemindex ItemIndex
	var item Item

	// Get form data
	item.Name     = c.FormValue("name")
	item.Category = c.FormValue("category")
	image, err   := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	// Open image
	src, err := image.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer src.Close()
	image_path := fmt.Sprintf("%x", image)

	// Log
	c.Logger().Infof("Receive item: %s", item.Name)
	c.Logger().Infof("Receive category: %s", item.Category)
	c.Logger().Infof("Receive image: %s", image_path)

	message := fmt.Sprintf("item received: %s", item.Name)

	// Hash
	hash := sha256.Sum256([]byte(image_path))
	hash_string := fmt.Sprintf("%x", hash)
	item.Image_name = hash_string + ".jpg"

	// Save image
	f, err := os.Create("images/" + hash_string + ".jpg")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer f.Close()

	io.Copy(f, src)

	// Add Items
	itemindex.Items = append(itemindex.Items, item)

	// categoriesテーブルに追加
	var category_id int
	// 入力したcategoryが存在する時
	err = db.QueryRow("SELECT id FROM categories WHERE name = ?", item.Category).Scan(&category_id)
	// 存在しない時
	if err == sql.ErrNoRows {
		cmd1 := "INSERT INTO categories (name) VALUES (?)"
		_, err := db.Exec(cmd1, item.Category)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		}
		err = db.QueryRow("SELECT id FROM categories WHERE name = ?", item.Category).Scan(&category_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		}
	}

	

	// items テーブルに追加
	cmd2 := "INSERT INTO items (name, category_id, image_name) VALUES(?, ?, ?)"
	_, err = db.Exec(cmd2, item.Name, category_id, item.Image_name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	res := Response{Message: message}
	return c.JSON(http.StatusOK, res)
}

//GET "/items/:id"
func showItem(c echo.Context) error {
	// Get id
	id, err := strconv.Atoi(c.Param("id")) 
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	
	// Open DB
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer db.Close()

	// レコード読み込み
	rows, err := db.Query("SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer rows.Close()

	// レコード取り出し
	showItems := ItemIndex{}
	for rows.Next() {
		var add Item
        err := rows.Scan(&add.Name, &add.Category, &add.Image_name)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
        }
		showItems.Items = append(showItems.Items, add)
    }
	
	// debug (idが0以下または配列の長さを超えるとき)
	length := len(showItems.Items)
	if id <= 0 || id > length{ 
		c.Logger().Infof("out of range")
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, showItems.Items[id-1])
}

// GET "/image/:imageFilename"
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

// GET "/search"
func searchItem (c echo.Context) error {
	keyword := c.FormValue("keyword")
	// Open DB
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer db.Close()

	// レコード読み込み
	var items []Item
	rows, err := db.Query("SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id WHERE items.name LIKE CONCAT('%', ?, '%') ", keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}
	defer rows.Close()

	// レコード取り出し
	for rows.Next() {
		var add Item
        err := rows.Scan(&add.Name, &add.Category, &add.Image_name)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
        }
		items = append(items, add)
    }

	return c.JSON(http.StatusOK, items)
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
	e.GET("/search", searchItem)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

