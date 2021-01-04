package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dhnt/oci/pkg/registry"
)

func main() {
	var port int
	var store string

	flag.IntVar(&port, "port", 5000, "port")
	flag.StringVar(&store, "store", "memory:", "image store URI")

	flag.Parse()

	addr := fmt.Sprintf(":%d", port)

	logger := log.New(os.Stdout, "registry: ", log.LstdFlags)
	logger.Printf("Server listening on %s...", addr)

	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	cfg := registry.Config{
		Store: store,
	}
	router.Handle("/", registry.New(&cfg, registry.Logger(logger)))

	s := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
