package middleware

import (
	"diskette/collections"
	"diskette/tokens"
	"diskette/util"

	"errors"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func CreateSessionMiddleware(userCollection *mgo.Collection, jwtKey []byte) echo.HandlerFunc {
	return func(c *echo.Context) error {
		st := c.Query("st")

		sessionToken, err := tokens.ParseSessionToken(jwtKey, st)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
		}

		var userDoc collections.UserDocument
		err = userCollection.FindId(bson.ObjectIdHex(sessionToken.UserId)).One(&userDoc)
		if err != nil {
			return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The session is not valid.")))
		}

		if sessionToken.CreatedAt.Before(userDoc.RejectSessionsBefore) {
			err = errors.New("The session has expired.")
			c.JSON(http.StatusUnauthorized, util.CreateErrResponse(err))
			return err
		}

		c.Set("sessionToken", sessionToken)
		c.Set("userDoc", userDoc)

		return nil
	}
}
