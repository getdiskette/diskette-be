package admin

import (
	"diskette/collections"
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo/bson"

	"github.com/labstack/echo"
)

// http POST localhost:5025/admin/create-user email="joe.doe@gmail.com" password="123" roles:='["admin"]' profile:='{"name": "Joe Doe", "lang": "en"}' X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) CreateUser(c *echo.Context) error {
	var request struct {
		Email    string                 `json:"email"`
		Password string                 `json:"password"`
		Roles    []string               `json:"roles"`
		Profile  map[string]interface{} `json:"profile"`
	}
	c.Bind(&request)

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'email'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'password'")))
	}

	if request.Roles == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'roles'")))
	}

	if request.Profile == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'profile'")))
	}

	count, err := service.userCollection.Find(bson.M{"email": request.Email}).Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, util.CreateErrResponse(errors.New("This email address is already being used.")))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	userDoc := collections.UserDocument{
		Id:          bson.NewObjectId(),
		Email:       request.Email,
		HashedPass:  hashedPass,
		Roles:       request.Roles,
		Profile:     request.Profile,
		CreatedAt:   time.Now(),
		IsSuspended: false,
		ConfirmedAt: time.Now(),
	}

	err = service.userCollection.Insert(userDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(userDoc))
}
