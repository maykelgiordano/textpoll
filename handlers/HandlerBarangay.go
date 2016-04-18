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

type BarangayHandler struct {
	sess *mgo.Session
}

// NewBarangay factory for BarangayController
func NewBarangayHandler(sess *mgo.Session) *BarangayHandler {
	return &BarangayHandler{sess}
}

//fetch list of barangay
func (handler BarangayHandler) Index(c *gin.Context) {
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
	brgys := []m.Barangay{}
	collection := handler.sess.DB("textpolldb").C("barangay") 
	collection.Find(nil).Sort("-createdat").All(&brgys)
	c.JSON(http.StatusOK, brgys)
}

// Create new barangay
func (handler BarangayHandler) Create(c *gin.Context) {
	brgy := m.Barangay{}
	c.Bind(&brgy)
	collection := handler.sess.DB("textpolldb").C("barangay") 
	result := m.Barangay{}
	err := collection.Find(bson.M{"barangayname": brgy.BarangayName}).One(&result)
	//check if barangay name is not existing
	if fmt.Sprintf("%s", err) == "not found" {
		// generate hashed password
		brgy.Id = bson.NewObjectId()
		brgy.CreatedAt = time.Now().UTC()
		brgy.UpdatedAt = time.Now().UTC()
		brgy.Status = "active"
		collection.Insert(&brgy)
		c.JSON(http.StatusCreated,brgy)
	} else {
		respond(http.StatusBadRequest,"Barangay name was already taken",c,true)
	}
}

// Update barangay
func (handler BarangayHandler) Update(c *gin.Context) {
	id := c.Param("id")
	brgy := m.Barangay{}
	c.Bind(&brgy)
	collection := handler.sess.DB("textpolldb").C("barangay") 
	result := m.Barangay{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	//check if barangay record exists
	if fmt.Sprintf("%s", err) == "not found" {
		respond(http.StatusBadRequest,"Barangay record not found",c,true)
	} else {
		//check if barangay name exists
		otherBrgy := m.Barangay{}
		err := collection.Find(bson.M{"$and": []bson.M{bson.M{"barangayname": brgy.BarangayName},
							bson.M{"_id" : bson.M{"$ne" : bson.ObjectIdHex(id)}}}}).One(&otherBrgy)
		fmt.Println("ERRR ---> ", err)
		if fmt.Sprintf("%s", err) == "not found" {
			//collection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"barangayname": brgy.BarangayName,
			//						"population" : brgy.Population, "updatedat" : time.Now().UTC()}})

			change := mgo.Change {
				Update: bson.M{"$set": bson.M{"barangayname": brgy.BarangayName,
								"updatedat" : time.Now().UTC()}},
				ReturnNew: true,
			}
			updatedBrgy := m.Barangay{}
			collection.FindId(bson.ObjectIdHex(id)).Apply(change, &updatedBrgy) // Apply
			c.JSON(http.StatusOK,updatedBrgy)
		} else {
			respond(http.StatusBadRequest,"Barangay name was already taken",c,true)
		}
	}
}

//show polling places in a barangay
func (handler BarangayHandler) Show(c *gin.Context) {
	id := c.Param("id")
	places := []m.PollingPlace{}
	collection := handler.sess.DB("textpolldb").C("pollingplace") 
	collection.Find(bson.M{"$and" : []bson.M{bson.M{"status": "active"},bson.M{"barangayid" : id}}}).Sort("-createdat").All(&places)
	c.JSON(http.StatusOK, places)
}


