package main

import (
	"dockertest/connect"
	discord_service "dockertest/discord"
	"dockertest/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDiscordBot()
	router := InitHttpsServer()
	err := connect.DG.Open()
	if err != nil {
		log.Fatal()
	}
	defer connect.DG.Close()

	fmt.Println("bot runing ...")

	go func() {
		if err := router.Run(":8000"); err != nil {
			log.Fatal(err)
		}
	}()

	connect.SetCommands(connect.DG)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}

func InitDiscordBot() {
	connect.Postgres()
	connect.CreateDiscordSession()
	connect.DG.AddHandler(discord_service.HandleComands)
	connect.SetCommands(connect.DG)
}

func InitHttpsServer() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)
	router.POST("/nofity_long", service.Cronnotified30minute)
	return router
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
