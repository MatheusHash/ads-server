package main

import (
	"ad-server/config"
	"ad-server/internal/database"
	"ad-server/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Conectar ao banco de dados SQLite
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.DB.Close()
	// Criar um novo roteador
	r := mux.NewRouter()

	// Inicializar as tabelas
	if err := database.InitTables(); err != nil {
		log.Fatalf("Erro ao inicializar tabelas: %v", err)
	}

	// Configurar rotas
	r.HandleFunc("/ad", handlers.GetAd).Methods("GET")
	r.HandleFunc("/ad/create", handlers.CreateAd).Methods("POST") // Novo endpoint para criar anúncio
	r.HandleFunc("/ad/click", handlers.RegisterClick).Methods("POST")
	r.HandleFunc("/ad/find/{slug:[0-9]+}", handlers.GetAdById).Methods("GET")

	// Iniciar servidor
	log.Printf("Servidor de anúncios iniciado na porta %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
