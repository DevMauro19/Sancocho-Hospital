package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/state", app.stateHandler)
	mux.HandleFunc("/api/attack", app.attackHandler)
	mux.HandleFunc("/api/wake", app.wakeHandler)
	mux.HandleFunc("/api/assign", app.assignHandler)
	mux.HandleFunc("/api/patients", app.addPatientHandler)
	mux.HandleFunc("/api/patients/", app.patientDetailHandler)
	mux.HandleFunc("/api/reset", app.resetHandler)
	mux.Handle("/", http.FileServer(http.Dir("cmd/web/static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	addr := ":" + port
	log.Printf("Costeñitos web disponible en http://127.0.0.1:%s", port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
