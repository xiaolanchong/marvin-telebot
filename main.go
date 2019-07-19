package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

func runTelebot() {
	fmt.Printf("Starting pulling bot...\n")

	botToken := os.Getenv("BOT_TOKEN")
	if len(botToken) == 0 {
		panic("BOT_TOKEN env variable not set")
	}

	dataRootDir := os.Getenv("DATA_ROOT")

	newBot, updates, err := bot.StartTeleBot(botToken)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			select {
			case <-c:
				os.Exit(1)
			case update := <-updates:
				bot.ProcessTeleBotUpdate(newBot, update, dataRootDir)
			}
		}
	}()
	for {
		time.Sleep(10 * time.Second)
	}
}

func main() {
	// runHeroku()
	runTelebot()
}
