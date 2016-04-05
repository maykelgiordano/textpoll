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

// Create new precinct
func (handler PrecinctHandler) Create(c *gin.Context) {
	precinct := m.Precinct{}
	c.Bind(&precinct)

	if precinct.PrecinctNo == "" {
		respond(http.StatusBadRequest,"Please specify the precinct no",c,true)
	} else if precinct.RegisteredVoters < 1 || strconv.Itoa(precinct.RegisteredVoters) == "" {
		respond(http.StatusBadRequest,"Please specify the no. of registered voters",c,true)
	} else if precinct.PollingPlaceId == "" {
		respond(http.StatusBadRequest,"Please specify the polling place id of the precinct",c,true)
	} else {
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
}

// Update a precinct
func (handler PrecinctHandler) Update(c *gin.Context) {
	id := c.Param("id")
	precinct := m.Precinct{}
	c.Bind(&precinct)

	if precinct.PrecinctNo == "" {
		respond(http.StatusBadRequest,"Please specify the precinct no",c,true)
	} else if precinct.RegisteredVoters < 1 || strconv.Itoa(precinct.RegisteredVoters) == "" {
		respond(http.StatusBadRequest,"Please specify the no. of registered voters",c,true)
	} else if precinct.PollingPlaceId == "" {
		respond(http.StatusBadRequest,"Please specify the polling place id of the precinct",c,true)
	} else {
		collection := handler.sess.DB("textpolldb").C("precinct") 
		result := m.Precinct{}
		err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)

		if fmt.Sprintf("%s", err) == "not found" {
			respond(http.StatusBadRequest,"Precinct not found",c,true)
		} else {
			otherPrecinct := m.Precinct{}
			err := collection.Find(bson.M{"$and": []bson.M{bson.M{"precinctno": precinct.PrecinctNo}, 
								bson.M{"_id" : bson.M{"$ne" : bson.ObjectIdHex(id)}}}}).One(&otherPrecinct)
			fmt.Println("ERRR ---> ", err)
			if fmt.Sprintf("%s", err) == "not found" {
				change := mgo.Change {
					Update: bson.M{"$set": bson.M{"precinctno": precinct.PrecinctNo,
									"pollingplaceid" : precinct.PollingPlaceId,
									"registeredvoters" : precinct.RegisteredVoters,
									"status" : precinct.Status, "updatedat" : time.Now().UTC()}},
					ReturnNew: true,
				}
				updatePrecinct := m.Precinct{}
				collection.FindId(bson.ObjectIdHex(id)).Apply(change, &updatePrecinct) // Apply
				c.JSON(http.StatusOK,updatePrecinct)
			} else {
				respond(http.StatusBadRequest,"Precint id already used",c,true)
			}
		}
	}
}


