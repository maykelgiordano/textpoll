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

type PollingPlaceHandler struct {
	sess *mgo.Session
}

// NewPollingPlace factory for PollingPlace Controller
func NewPollingPlaceHandler(sess *mgo.Session) *PollingPlaceHandler {
	return &PollingPlaceHandler{sess}
}

//fetch list of polling place
func (handler PollingPlaceHandler) Index(c *gin.Context) {
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
	places := []m.PollingPlace{}
	collection := handler.sess.DB("textpolldb").C("pollingplace") 
	collection.Find(bson.M{"status": "active"}).Sort("-createdat").All(&places)
	c.JSON(http.StatusOK, places)
}

// Create new polling place
func (handler PollingPlaceHandler) Create(c *gin.Context) {
	pollingPlace := m.PollingPlace{}
	c.Bind(&pollingPlace)
	collection := handler.sess.DB("textpolldb").C("pollingplace") 
	result := m.PollingPlace{}
	err := collection.Find(bson.M{"place": pollingPlace.Place}).One(&result)
	//check if polling place is not existing
	if fmt.Sprintf("%s", err) == "not found" {
		pollingPlace.Id = bson.NewObjectId()
		pollingPlace.CreatedAt = time.Now().UTC()
		pollingPlace.UpdatedAt = time.Now().UTC()
		pollingPlace.Status = "active"
		collection.Insert(&pollingPlace)
		c.JSON(http.StatusCreated,pollingPlace)
	} else {
		respond(http.StatusBadRequest,"Polling place name was already taken",c,true)
	}
}

// Update polling place
func (handler PollingPlaceHandler) Update(c *gin.Context) {
	id := c.Param("id")
	place := m.PollingPlace{}
	c.Bind(&place)
	collection := handler.sess.DB("textpolldb").C("pollingplace") 
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


