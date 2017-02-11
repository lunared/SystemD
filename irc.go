package main

import (
	"os"
	"os/signal"

	"./bot"
)

//Starts the bot and listens for a ^C
func main() {

	bot := bot.Run()

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		bot.CloseConn()

		os.Exit(0)
	}()
}
