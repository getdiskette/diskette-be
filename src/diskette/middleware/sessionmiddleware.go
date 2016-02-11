package middleware

import (
	"errors"
	"net/http"

	"diskette/collections"
	"diskette/tokens"
	"diskette/util"

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
		userCollection.FindId(bson.ObjectIdHex(sessionToken.UserID)).One(&userDoc)

		if sessionToken.CreatedAt.Before(userDoc.RejectSessionsBefore) {
			return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The session has expired.")))
		}

		return nil
	}
}
