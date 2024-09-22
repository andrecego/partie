package main

import (
	"fmt"
	"net/http"
	"partie-bot/bot"
	"partie-bot/config"

	"github.com/rollbar/rollbar-go"
)

func main() {
	err := config.ReadConfig()
	// configRollbar()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// rollbar.WrapAndWait(bot.Start)
	// defer rollbar.Close()
	bot.Start()

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	}()

	<-make(chan struct{})
	return
}

func configRollbar() {
	rollbar.SetToken(config.RollbarToken)
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot("github.com/andrecego/partie")
}
