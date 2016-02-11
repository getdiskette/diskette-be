package rest

import (
	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

// Service that exposes the common REST operations:
// GET, POST, PUT, DELETE
type Service interface {
	Get(c *echo.Context) error
	Post(c *echo.Context) error
	Put(c *echo.Context) error
	Delete(c *echo.Context) error
}

type serviceImpl struct {
	db *mgo.Database
}

// NewService creates an instance of rest.Service that allows REST operations
// on the specified mongodb database collections.
func NewService(db *mgo.Database) Service {
	return &serviceImpl{db}
}
