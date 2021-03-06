package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mediocregopher/radix.v2/redis"
)



//eventually get rid of these global variables
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

func main() {
	//redis
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	//db
	var db *gorm.DB
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&UserAuth{}, &UserProfile{}, &Qotd{}, &QotdAnswerOption{}, &QotdAnswer{}, &FeedbackQuestion{}, &FeedbackAnswer{}, &Kinship{}, &Chat{})

	defer db.Close()

	//websockets
	go wsMessages()

	//routes
	http.Handle("/", http.FileServer(http.Dir("../public/")))
	http.Handle("/public/assets/", http.StripPrefix("/public/assets/", http.FileServer(http.Dir("../public/assets/"))))
	http.Handle("/bundles/", http.StripPrefix("/bundles/", http.FileServer(http.Dir("../bundles/"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/api/profile", profileHandler(db, conn))
	http.Handle("/api/signup", signupHandler(db, conn))
	http.Handle("/api/login", loginHandler(db, conn))
	http.Handle("/api/tokenCheck", tokenHandler(conn))
	http.Handle("/api/feedback", feedbackHandler(db))
	http.Handle("/api/visitCheck", visitHandler(conn));
	http.Handle("/api/ws", wsHandler(conn))
	http.Handle("/api/qotd", qotdHandler(db))
	
	// http.Handle("/api/kinships", kinships)

	//Initialize
	http.ListenAndServe(":8080", nil)
}
