package user

import (
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const USER_COLLECTION = "user"

type UserService interface {
	Signup(c *echo.Context) error
}

type impl struct {
	db     *mgo.Database
	jwtKey []byte
}

func NewUserService(db *mgo.Database, jwtKey []byte) UserService {
	return &impl{db, jwtKey}
}

// POST /user/signup BODY={doc}
// examples:
// http POST localhost:5025/user/signup name="Joe Doe" email=joe.doe@gmail.com password=abc language=en
func (self impl) Signup(c *echo.Context) error {
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Language string `json:"language"`
	}
	c.Bind(&request)

	if request.Name == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'name'")))
	}

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'email'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'password'")))
	}

	if request.Language == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'language'")))
	}

	confirmationTokenStr, err := self.createUser(request.Name, request.Email, request.Password, request.Language, false)
	if err != nil {
		return c.JSON(http.StatusConflict, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(confirmationTokenStr))
}

func (self impl) createUser(name, email, password, language string, isConfirmed bool) (confirmationTokenStr string, err error) {
	count, err := self.db.C(USER_COLLECTION).Find(bson.M{"email": email}).Count()
	if err != nil {
		return
	}

	if count > 0 {
		return "", errors.New("Error: This email address is already being used.")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	userDoc := UserDoc{
		Name:        name,
		Email:       email,
		HashedPass:  hashedPass,
		Language:    language,
		CreatedAt:   time.Now(),
		IsSuspended: false,
	}

	if isConfirmed {
		userDoc.ConfirmedAt = time.Now()

	} else {
		userDoc.ConfirmationKey = uuid.NewV4().String()

		confirmationToken := privateConfirmationToken{
			database: self.db.Name,
			key:      userDoc.ConfirmationKey,
			language: userDoc.Language,
		}

		confirmationTokenStr, err = confirmationToken.toString(self.jwtKey)
		if err != nil {
			return
		}
	}

	err = self.db.C(USER_COLLECTION).Insert(userDoc)
	return
}
