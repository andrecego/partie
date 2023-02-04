package main

import (
	"fmt"
	"partie-bot/bot"
	"partie-bot/config"

	"github.com/rollbar/rollbar-go"
)

func main() {
	err := config.ReadConfig()
	configRollbar()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rollbar.WrapAndWait(bot.Start)
	defer rollbar.Close()

	<-make(chan struct{})
	return
}

func configRollbar() {
	rollbar.SetToken(config.RollbarToken)
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot("github.com/andrecego/partie")
}
