package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"fmt"
	
	"marvin-telebot/bot"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func runHeroku() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

func runPuller() {
	bot_token := os.Getenv("BOT_TOKEN")
	if len(bot_token) == 0 {
		panic("BOT_TOKEN env variable not set")
	}

	puller, _ := bot.New(500, bot_token)
	defer puller.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for {
			select
			{
			case <- c:
				os.Exit(1)
			case update := <- puller.UpdatesChannel:
				fmt.Printf("New msg: %+v\n", update)
			}
		}
	}()
	for { time.Sleep(10 * time.Second) }
}

func main() {
	// runHeroku()
	runPuller()

}
