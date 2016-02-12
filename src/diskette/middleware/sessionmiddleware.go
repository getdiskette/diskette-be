package middleware

import (
	"diskette/collections"
	"diskette/tokens"
	"diskette/util"

	"errors"
	"net/http"

	"log"

	"github.com/labstack/echo"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func CreateSessionMiddleware(userCollection *mgo.Collection, jwtKey []byte) echo.HandlerFunc {
	return func(c *echo.Context) error {
		st := c.Request().Header.Get("X-Diskette-Session-Token")

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

		log.Print("---")
		log.Printf("sessionToken %+v", sessionToken)
		log.Printf("userDoc %+v", userDoc)

		c.Set("sessionToken", sessionToken)
		c.Set("userDoc", userDoc)

		return nil
	}
}
