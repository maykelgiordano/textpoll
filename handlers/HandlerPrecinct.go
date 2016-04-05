package handlers

import (
	"net/http"
	"strconv"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	m "txtpoll/sm/api/models"
	"gopkg.in/mgo.v2"	
	"gopkg.in/mgo.v2/bson"	
)

type PrecinctHandler struct {
	sess *mgo.Session
}

// NewPrecinct factory for Precinct Controller
func NewPrecinctHandler(sess *mgo.Session) *PrecinctHandler {
	return &PrecinctHandler{sess}
}

//fetch list of precincts place
func (handler PrecinctHandler) Index(c *gin.Context) {
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
	precincts := []m.Precinct{}
	collection := handler.sess.DB("textpolldb").C("precinct") 
	collection.Find(bson.M{"status": "active"}).Sort("-createdat").All(&precincts)
	c.JSON(http.StatusOK, precincts)
}

// Create new polling place
func (handler PrecinctHandler) Create(c *gin.Context) {
	precinct := m.Precinct{}
	c.Bind(&precinct)
	collection := handler.sess.DB("textpolldb").C("precinct") 
	result := m.Precinct{}
	err := collection.Find(bson.M{"precinctno": precinct.PrecinctNo}).One(&result)
	
	if fmt.Sprintf("%s", err) == "not found" {
		precinct.Id = bson.NewObjectId()
		precinct.CreatedAt = time.Now().UTC()
		precinct.UpdatedAt = time.Now().UTC()
		precinct.Status = "active"
		collection.Insert(&precinct)
		c.JSON(http.StatusCreated,precinct)
	} else {
		respond(http.StatusBadRequest,"Precinct id already existing",c,true)
	}
}

// Update polling place
func (handler PrecinctHandler) Update(c *gin.Context) {
	id := c.Param("id")
	place := m.PollingPlace{}
	c.Bind(&place)
	collection := handler.sess.DB("textpolldb").C("precinct") 
	result := m.PollingPlace{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	//check if polling place record exists
	if fmt.Sprintf("%s", err) == "not found" {
		respond(http.StatusBadRequest,"Polling place record not found",c,true)
	} else {
		//check if polling place name exists
		otherPlace := m.PollingPlace{}
		err := collection.Find(bson.M{"$and": []bson.M{bson.M{"place": place.Place}, 
							bson.M{"_id" : bson.M{"$ne" : bson.ObjectIdHex(id)}}}}).One(&otherPlace)
		fmt.Println("ERRR ---> ", err)
		if fmt.Sprintf("%s", err) == "not found" {
			change := mgo.Change {
				Update: bson.M{"$set": bson.M{"place": place.Place,
								"status" : place.Status, "updatedat" : time.Now().UTC()}},
				ReturnNew: true,
			}
			updatePlace := m.PollingPlace{}
			collection.FindId(bson.ObjectIdHex(id)).Apply(change, &updatePlace) // Apply
			c.JSON(http.StatusOK,updatePlace)
		} else {
			respond(http.StatusBadRequest,"Polling place name was already taken",c,true)
		}
	}
}


