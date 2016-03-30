package main


import (
	"os"
	"fmt"

	"gopkg.in/mgo.v2"
	h "txtpoll/sm/api/handlers"
	"txtpoll/sm/api/config"
	"github.com/gin-gonic/gin"
)

func main() {
	//openFTP()
	db := *InitDB()
	router := gin.Default()
	LoadAPIRoutes(router, &db)
}

func LoadAPIRoutes(r *gin.Engine, db *mgo.Session) {
	public := r.Group("/api/v1")

	//manage users
	userHandler := h.NewUserHandler(db)
	public.GET("/users", userHandler.Index)
	public.POST("/users", userHandler.Create)
	public.POST("/auth", userHandler.Auth)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	fmt.Println("PORT ---> ",port)
	r.Run(fmt.Sprintf(":%s", port))
}

func InitDB() *mgo.Session {
	fmt.Prntln("DB URL ---> ", config.GetString("DB_URL"))
	sess, err := mgo.Dial(config.GetString("DB_URL"))
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database:  %s", err))
	}
	sess.SetSafe(&mgo.Safe{})
	return sess
}