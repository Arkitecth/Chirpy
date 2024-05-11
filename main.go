package main

import (
	"log"
	"net/http"
	"strconv"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(200)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	message := "Hits: " + strconv.Itoa(cfg.fileSeverHits)
	w.Write([]byte(message))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileSeverHits = 0
}

type apiConfig struct {
	fileSeverHits int
}

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileSeverHits += 1 
		next.ServeHTTP(w, r)
	})

}

func main() {
	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{fileSeverHits: 0}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	mux.Handle("/app/*", apiCfg.middleWareMetricsInc(middlewareCors(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))

	srv := http.Server {
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}