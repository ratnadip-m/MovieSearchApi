package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Year   int     `json:"year"`
	Rating float32 `json:"rating"`
	Genres string  `json:"genres"`
}
type MovieModel struct {
	ID     int      `gorm:"primaryKey"`
	Title  string   `gorm:"not null"`
	Year   string   `gorm:"not null"`
	Rating string   `gorm:"not null"`
	Genres []string `gorm:"-"`
}

var movies []MovieModel
var db *gorm.DB

func main() {
	var err error
	router := gin.Default()
	dsn := "root:cdac@tcp(localhost:3306)/moviedb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// db, err = sql.Open("mysql", "root:cdac@tcp(localhost:3306)/moviedb")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	db.AutoMigrate(&Movie{})

	router.POST("/movies", addMovie)
	router.GET("/searchmovies", GetMovie)
	router.GET("/searchmoviesbyidyeargenres", GetMovieByIDYearRatingGenres)

	router.Run(":8087")
}

func addMovie(c *gin.Context) {
	var movie Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign a new ID to the movie
	// movie.ID = len(movies) + 1

	// Add the movie to the list of movies
	// movies = append(movies, movie)
	// result, err := db.Exce("INSERT INTO orders (id, title, year, rating, generes) VALUES (?, ?, ?, ?, ?)",
	// 	movie.ID, movie.Title, movie.Year, movie.Rating, movie.Genres)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	db.Save(&movie)

	c.JSON(200, "result")
}

func GetMovie(c *gin.Context) {
	var movie Movie
	title := c.Query("title")

	db.Where("title = ?", title).First(&movie)

	if movie.Title != "" {
		c.JSON(http.StatusOK, movie)
	} else {
		// request := gorequest.New()
		// _, body, _ := request.Get("http://www.omdbapi.com/").
		// 	Query("apikey=YOUR_API_KEY").
		// 	Query("t=" + title).
		// 	End()
		// resp, err := http.Get("http://www.omdbapi.com/?t=" + title)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
		// 	return
		// }
		// defer resp.Body.Close()
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, err := http.NewRequest("GET", "http://www.omdbapi.com/?t="+title, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("apikey", "64cbe519")

		// Send the request using the HTTP client
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Parse the response body to a Movie struct
		var omdbMovie struct {
			// ID     int     `json:"id"`
			Title  string  `json:"title"`
			Year   int     `json:"year"`
			Rating float32 `json:"rating"`
			Genres string  `json:"genres"`
		}
		// if err := json.Unmarshal([]byte(resp.Body), &omdbMovie); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		// 	return
		// }
		if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
			return
		}

		if omdbMovie.Title != "" {
			// Save the movie in the database
			movie := Movie{Title: omdbMovie.Title, Year: omdbMovie.Year, Rating: omdbMovie.Rating, Genres: omdbMovie.Genres}
			db.Save(&movie)

			c.JSON(http.StatusOK, movie)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		}
	}
}

func GetMovieByIDYearRatingGenres(c *gin.Context) {
	var movie Movie
	id := c.Query("id")
	year := c.Query("year")
	generes := c.Query("generes")
	rating := c.Query("rating")
	if id != "" {
		db.Where("id = ?", id).First(&movie)

		if movie.ID != 0 {
			c.JSON(http.StatusOK, movie)
			return
		}
	} else if year != "" {
		db.Where("year = ?", year).First(&movie)

		if movie.Year != 0 {
			c.JSON(http.StatusOK, movie)
			return
		}
	} else if generes != "" {
		db.Where("generes = ?", generes).First(&movie)

		if movie.Genres != "" {
			c.JSON(http.StatusOK, movie)
			return
		}
	} else if rating != "" {
		db.Where("rating = ?", rating).First(&movie)

		if movie.Rating != 0 {
			c.JSON(http.StatusOK, movie)
			return
		}
	}
}
