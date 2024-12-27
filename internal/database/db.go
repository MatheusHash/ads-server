package database

import (
	"ad-server/config"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Connect abre a conexão com o banco de dados SQLite
func Connect(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		return err
	}

	// Verificar se a conexão está funcionando
	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("Conectado ao banco de dados SQLite.")
	return nil
}

// InitTables cria as tabelas necessárias para o aplicativo
func InitTables() error {
	createAdsTable := `
    CREATE TABLE IF NOT EXISTS ads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category TEXT NOT NULL,
        content TEXT NOT NULL,
        impressions INTEGER DEFAULT 0
    );`

	_, err := DB.Exec(createAdsTable)
	if err != nil {
		return err
	}
	log.Println("Tabela `ads` verificada/criada com sucesso.")

	createAccountsTable := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(65) NOT NULL,
			type TEXT CHECK(type IN ('company', 'personal'))
		)
	`

	_, err = DB.Exec(createAccountsTable)
	if err != nil {
		return err
	}

	log.Println("Tabela `accounts` verficada/criada com sucesso.")

	return nil
}
