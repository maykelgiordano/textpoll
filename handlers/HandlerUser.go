package handlers

import (
	"net/http"
	"strconv"
	"fmt"
	"time"
	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	m "txtpoll/sm/api/models"
	"gopkg.in/mgo.v2"	
	"gopkg.in/mgo.v2/bson"	
	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	sess *mgo.Session
}

// NewAppoointment factory for AppointmentsController
func NewUserHandler(sess *mgo.Session) *UserHandler {
	return &UserHandler{sess}
}

//fetch list of users
func (handler UserHandler) Index(c *gin.Context) {
	start := -1
	max := 10

	//check if start exists in url parameters
	if c.Query("start") != ""  {
		i,_ := strconv.Atoi(c.Query("start"))
		start = i;
	} else {
		fmt.Println("cant read start query param")
	}

	if c.Query("max") != ""  {
		i,_ := strconv.Atoi(c.Query("max"))
		max = i;
	} 

	fmt.Printf("offset ---> %d max ---> %d\n", start, max)
	users := []m.User{}
	collection := handler.sess.DB("textpolldb").C("users") 
	collection.Find(nil).Select(bson.M{"password": 0}).All(&users)
	c.JSON(http.StatusOK, users)
}

// Create an account
func (handler UserHandler) Create(c *gin.Context) {
	user := m.User{}
	c.Bind(&user)
	collection := handler.sess.DB("textpolldb").C("users") 
	result := m.User{}
	err := collection.Find(bson.M{"email": user.Email}).One(&result)
	//check if email is not existing
	if fmt.Sprintf("%s", err) == "not found" {
		// generate hashed password
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = time.Now().UTC()
		user.Status = "active"
		user.IsDefaultPassword = true
		hasher := md5.New()
    	hasher.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(hasher.Sum(nil))
		collection.Insert(&user)
	    // Create the token
	    token := jwt.New(jwt.SigningMethodHS256)
	    token.Claims["id"] = user.Email
	    token.Claims["iat"] = time.Now().Unix()
	    token.Claims["exp"] = time.Now().Add(time.Second * 3600 * 24).Unix()
	    tokenString, err := token.SignedString([]byte("secret"))
	    if err == nil {
	    	response := &m.CurrentUser{FirstName: user.FirstName, LastName: user.LastName,
		    		ContactNo: user.ContactNo, Email: user.Email, UserRole: user.UserRole, 
		    		IsDefaultPassword: user.IsDefaultPassword, Token: tokenString}
	    	c.JSON(http.StatusCreated,response)
    	} else {
			respond(http.StatusBadRequest,"Failed to create account",c,true)
    	}
	} else {
		respond(http.StatusBadRequest,"Email was already taken",c,true)
	}
}

// Authenticate user
func (handler UserHandler) Auth(c *gin.Context) {
	auth := m.Auth{}
	c.Bind(&auth)

	collection := handler.sess.DB("textpolldb").C("users") 
	result := m.User{}
	error := collection.Find(bson.M{"$and": []bson.M{bson.M{"email": auth.Email}, 
							bson.M{"status": "active"}}}).One(&result)
	if fmt.Sprintf("%s", error) == "not found" {
		respond(http.StatusBadRequest,"Account not found",c,true)
	} else {
		hasher := md5.New()
    	hasher.Write([]byte(auth.Password))
		if result.Password == hex.EncodeToString(hasher.Sum(nil)) {
		    // Create the token
		    token := jwt.New(jwt.SigningMethodHS256)
		    token.Claims["id"] = result.Email
		    token.Claims["iat"] = time.Now().Unix()
		    token.Claims["exp"] = time.Now().Add(time.Second * 3600 * 24).Unix()
		    tokenString, err := token.SignedString([]byte("secret"))
		    if err == nil {
		    	response := &m.CurrentUser{FirstName: result.FirstName, LastName: result.LastName,
		    		ContactNo: result.ContactNo, Email: result.Email, UserRole: result.UserRole, 
		    		IsDefaultPassword: result.IsDefaultPassword, Token: tokenString}
		    	c.JSON(http.StatusOK,response)
	    	} else {
				respond(http.StatusBadRequest,"Failed to create token",c,true)
			}
		} else {
			respond(http.StatusBadRequest,"Invalid password",c,true)
		}
	}
}




