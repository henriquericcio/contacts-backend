package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/mgo.v2/bson"
)

type contact struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName   string        `json:"firstName" bson:"firstName" binding:"required"`
	LastName    string        `json:"lastName" bson:"lastName"  binding:"required"`
	PhoneNumber string        `json:"phoneNumber" bson:"phoneNumber"`
}

type repoi interface {
	len() int
	getAll() []contact
	getByID(id string) contact
	store(c *contact)
	remove(c contact)
	close()
}

func main() {
	//environment --------
	godotenv.Load()
	environment := os.Getenv("ENVIRONMENT")

	//database --------
	var repo repoi
	repo = newRepoMemory()
	defer repo.close()

	//api --------
	e := echo.New()

	if environment == "development" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.HideBanner = true

	e.POST("/contacts", func(c echo.Context) error {
		payload := new(contact)
		if err := c.Bind(payload); err != nil {
			return err
		}

		//todo: review responsibility
		//payload.ID = uuid.New().String()

		repo.store(payload)

		return c.NoContent(http.StatusCreated)
	})

	e.GET("/contacts", func(c echo.Context) error {
		return c.JSON(http.StatusOK, repo.getAll())
	})

	e.GET("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		co := repo.getByID(id)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, co)
	})

	e.PUT("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		payload := new(contact)
		if err := c.Bind(payload); err != nil {
			return err
		}

		co := repo.getByID(id)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		co.FirstName = payload.FirstName
		co.LastName = payload.LastName
		co.PhoneNumber = payload.PhoneNumber

		repo.store(&co)

		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		id := c.Param("id")

		co := repo.getByID(id)
		if co == (contact{}) {
			return c.NoContent(http.StatusNotFound)
		}

		repo.remove(co)

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8010"))
}
