package main

import (
	"fmt"
	"gin-simple/pkg/settings"
	"gin-simple/routers"
	"log"
	"net/http"
)

func main() {

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf("localhost:%d", settings.HTTPPort),
		Handler:        router,
		ReadTimeout:    settings.ReadTimeout,
		WriteTimeout:   settings.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Error on listen:%[1]v", err)
	}
}
