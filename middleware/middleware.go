package middleware

import (
	"github.com/getdiskette/diskette/collections"
	"github.com/getdiskette/diskette/tokens"
	"github.com/getdiskette/diskette/util"

	"errors"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateSessionMiddleware(userCollection *mgo.Collection, jwtKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			st := c.Request().Header().Get("X-Diskette-Session-Token")

			sessionToken, err := tokens.ParseSessionToken(jwtKey, st)
			if err != nil {
				c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
				return err
			}

			var userDoc collections.UserDocument
			err = userCollection.FindId(bson.ObjectIdHex(sessionToken.UserId)).One(&userDoc)
			if err != nil {
				err = errors.New("The session is not valid.")
				c.JSON(http.StatusNotFound, util.CreateErrResponse(err))
				return err
			}

			if sessionToken.CreatedAt.Before(userDoc.SignedOutAt) {
				err = errors.New("The session has expired.")
				c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
				return err
			}

			c.Set("sessionToken", sessionToken)
			c.Set("userDoc", userDoc)
			return next(c)
		}
	}
}

func CreateAdminSessionMiddleware(userCollection *mgo.Collection, jwtKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			st := c.Request().Header().Get("X-Diskette-Session-Token")

			sessionToken, err := tokens.ParseSessionToken(jwtKey, st)
			if err != nil {
				c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
				return err
			}

			var userDoc collections.UserDocument
			err = userCollection.FindId(bson.ObjectIdHex(sessionToken.UserId)).One(&userDoc)
			if err != nil {
				err = errors.New("The session is not valid.")
				c.JSON(http.StatusNotFound, util.CreateErrResponse(err))
				return err
			}

			if sessionToken.CreatedAt.Before(userDoc.SignedOutAt) {
				err = errors.New("The session has expired.")
				c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
				return err
			}

			isAdmin := false
			for _, role := range userDoc.Roles {
				if role == "admin" {
					isAdmin = true
				}
			}

			if !isAdmin {
				err = errors.New("Not admin.")
				c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
				return err
			}

			c.Set("sessionToken", sessionToken)
			c.Set("userDoc", userDoc)

			return next(c)
		}
	}
}
