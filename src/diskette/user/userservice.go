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
	ConfirmSignup(c *echo.Context) error
	Signin(c *echo.Context) error
	ForgotPassword(c *echo.Context) error
	ResetPassword(c *echo.Context) error
}

type impl struct {
	db     *mgo.Database
	jwtKey []byte
}

func NewUserService(db *mgo.Database, jwtKey []byte) UserService {
	return &impl{db, jwtKey}
}

// Example:
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

	return self.createUser(c, request.Name, request.Email, request.Password, request.Language, false)
}

func (self impl) createUser(c *echo.Context, name, email, password, language string, isConfirmed bool) error {
	count, err := self.db.C(USER_COLLECTION).Find(bson.M{"email": email}).Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, util.CreateErrResponse(errors.New("Error: This email address is already being used.")))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	userDoc := UserDoc{
		Name:        name,
		Email:       email,
		HashedPass:  hashedPass,
		Language:    language,
		CreatedAt:   time.Now(),
		IsSuspended: false,
	}

	var tokenStr string

	if isConfirmed {
		userDoc.ConfirmedAt = time.Now()

	} else {
		userDoc.ConfirmationKey = uuid.NewV4().String()

		token := confirmationToken{key: userDoc.ConfirmationKey}

		tokenStr, err = token.toString(self.jwtKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
		}
	}

	err = self.db.C(USER_COLLECTION).Insert(userDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(bson.M{"confirmationToken": tokenStr}))
}

// http POST localhost:5025/user/confirm token=<confirmation_token>
func (self impl) ConfirmSignup(c *echo.Context) error {
	var request struct {
		Token string `json:"token"`
	}
	c.Bind(&request)

	if request.Token == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'token'")))
	}

	token, err := parseConfirmationToken(self.jwtKey, request.Token)
	if err != nil || token.key == "" {
		return c.JSON(http.StatusForbidden, util.CreateErrResponse(err))
	}

	return self.db.C(USER_COLLECTION).Update(
		bson.M{"confirmationKey": token.key},
		bson.M{
			"$set": bson.M{
				"confirmedAt": time.Now(),
			},
		},
	)
}

// http POST localhost:5025/user/signin email=joe.doe@gmail.com password=abc
func (self impl) Signin(c *echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.Bind(&request)

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'email'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'password'")))
	}

	var userDoc UserDoc
	err := self.db.C(USER_COLLECTION).Find(bson.M{"email": request.Email}).One(&userDoc)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	if userDoc.ConfirmedAt.Before(userDoc.CreatedAt) {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The user has not confirmed the account.")))
	}

	if userDoc.IsSuspended {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The user is suspended.")))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDoc.HashedPass), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The password didn't match.")))
	}

	token := SessionToken{
		Id:        uuid.NewV4().String(),
		UserId:    userDoc.Id.Hex(),
		CreatedAt: time.Now(),
	}

	tokenStr, err := token.toString(self.jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(bson.M{"sessionToken": tokenStr}))
}

// http POST localhost:5025/user/forgot-password email=joe.doe@gmail.com
func (self impl) ForgotPassword(c *echo.Context) error {
	var request struct {
		Email string `json:"email"`
	}
	c.Bind(&request)

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'email'")))
	}

	resetKey := uuid.NewV4().String()

	err := self.db.C(USER_COLLECTION).Update(
		bson.M{"email": request.Email},
		bson.M{
			"$set": bson.M{
				"resetKey":         resetKey,
				"requestedResetAt": time.Now(),
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	token := resetToken{key: resetKey}
	tokenStr, err := token.toString(self.jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(bson.M{"resetToken": tokenStr}))
}

// http POST localhost:5025/user/reset-passwort token=<reset_token> password=123
func (self impl) ResetPassword(c *echo.Context) error {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	c.Bind(&request)

	if request.Token == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'token'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Error: Missing parameter 'password'")))
	}

	token, err := parseResetToken(self.jwtKey, request.Token)
	if err != nil || token.key == "" {
		return c.JSON(http.StatusForbidden, util.CreateErrResponse(err))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	err = self.db.C(USER_COLLECTION).Update(
		bson.M{"resetKey": token.key},
		bson.M{
			"$set": bson.M{
				"resetKey":         "",
				"requestedResetAt": time.Unix(0, 0),
				"hashedPass":       hashedPass,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The token doesn't exist.")))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
