// LandGut-App Backend – JSON-basiert
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Provider struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Hours    string  `json:"hours"`
	Payment  string  `json:"payment"`
	Homepage string  `json:"homepage"`
}

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price,omitempty"`
}

var providers []Provider
var products map[string][]Product

func loadJSONData() {
	// Anbieter laden
	pf, err := os.Open("providers.json")
	if err != nil {
		log.Fatal("Fehler beim Öffnen von providers.json:", err)
	}
	defer pf.Close()
	if err := json.NewDecoder(pf).Decode(&providers); err != nil {
		log.Fatal("Fehler beim Einlesen von providers.json:", err)
	}

	// Produkte laden
	prf, err := os.Open("products.json")
	if err != nil {
		log.Fatal("Fehler beim Öffnen von products.json:", err)
	}
	defer prf.Close()
	if err := json.NewDecoder(prf).Decode(&products); err != nil {
		log.Fatal("Fehler beim Einlesen von products.json:", err)
	}
}

func providersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func providerDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	for _, p := range providers {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.NotFound(w, r)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if p, ok := products[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	} else {
		http.NotFound(w, r)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	loadJSONData()
	
	switch r.URL.Path {
	case "/api/providers":
		providersHandler(w, r)
	case "/api/provider":
		providerDetailHandler(w, r)
	case "/api/products":
		productsHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	loadJSONData()
	http.HandleFunc("/api/providers", providersHandler)
	http.HandleFunc("/api/provider", providerDetailHandler)
	http.HandleFunc("/api/products", productsHandler)
	
	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index_test1_final.html")
		} else {
			http.FileServer(http.Dir(".")).ServeHTTP(w, r)
		}
	})
	
	log.Println("Server läuft auf http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
