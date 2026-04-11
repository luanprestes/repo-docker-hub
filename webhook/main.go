package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	addr       = ":9000"
	deployScript = "/app/deploy.sh"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Responde imediatamente pro Docker Hub (evita timeout)
	w.WriteHeader(http.StatusOK)

	// Roda o deploy em background
	go func() {
		log.Println("webhook recebido — iniciando deploy...")

		cmd := exec.Command("/bin/sh", deployScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("erro no deploy: %v", err)
			return
		}

		log.Println("deploy concluído com sucesso")
	}()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/dockerhub", webhookHandler)
	mux.HandleFunc("/health", healthHandler)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("webhook listener iniciado em %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
