package restservice

import (
	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type RestService interface {
	Get(c *echo.Context) error
	Post(c *echo.Context) error
	Put(c *echo.Context) error
	Delete(c *echo.Context) error
}

type impl struct {
	db *mgo.Database
}

func NewRestService(db *mgo.Database) RestService {
	return &impl{db}
}
