package main

import (
	"log"
	"time"

	"github.com/Sonal000/golang-simple-microservice-project/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err := account.NewAccountRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		log.Print(r)
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080")
	s := account.NewAccountService()
	log.Fatal(account.ListenGRPC(s, 8080))
}
