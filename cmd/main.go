package main

import (
	"ad-server/config"
	"ad-server/internal/database"
	"ad-server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	// Conectar ao banco de dados SQLite
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.DB.Close()

	// Inicializar as tabelas
	if err := database.InitTables(); err != nil {
		log.Fatalf("Erro ao inicializar tabelas: %v", err)
	}

	// Configurar rotas
	http.HandleFunc("/ad", handlers.GetAd)
	http.HandleFunc("/ad/create", handlers.CreateAd) // Novo endpoint para criar anúncio
	// http.HandleFunc("/click", handlers.RegisterClick)

	// Iniciar servidor
	log.Printf("Servidor de anúncios iniciado na porta %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
