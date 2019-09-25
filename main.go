package main

import (
	"github.com/Sean-Der/fail2go"
	"github.com/Strum355/log"
	"github.com/go-chi/chi"

	"fail2rest/api"
	"fail2rest/services"

	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Initialise logger
	log.InitSimpleLogger(&log.Config{})

	// Gather environment variables
	port := "8080"
	fail2ban := "/var/run/fail2ban/fail2ban.sock"

	// Initialise Fail2Ban connection
	log.WithFields(log.Fields{
		"fail2ban socket": fail2ban,
	}).Info("Initialising fail2ban connection")

	conn := fail2go.Newfail2goConn(fail2ban)

	// Start HTTP Server
	log.WithFields(log.Fields{
		"port": port,
	}).Info("Initialising HTTP Server")

	r := chi.NewRouter()

	// Register service with Consul
	consul := services.ConsulService{
		ConsulHost:  "127.0.0.1:8500",
		ConsulToken: "",
		ServiceAddr: "127.0.0.1",
		Port:        8080,
		TTL:         time.Second * 5,
		Secret:      "abcd",
	}
	err := consul.Setup()
	if err != nil {
		log.WithError(err).Error("Could not setup Consul service")
		os.Exit(1)
	}
	err = consul.Register()
	if err != nil {
		log.WithError(err).Error("Could not register with Consul service")
		os.Exit(1)
	}

	// Initialise API
	api := api.API{Fail2Conn: conn, Secret: consul.Secret}
	api.Register(r)

	err = http.ListenAndServe(":"+fmt.Sprint(port), r)
	if err != nil {
		log.WithError(err).Error("Error serving HTTP")
	}
}