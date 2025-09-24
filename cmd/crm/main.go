package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"refactor_crm_interface/internal/app"
	"refactor_crm_interface/internal/storage"
)

func main() {
	// Fichier JSON de persistance
	dataDir := "data"
	_ = os.MkdirAll(dataDir, 0o755)
	jsonPath := filepath.Join(dataDir, "contacts.json")

	// Choix du store : JSON persistant
	store, err := storage.NewJSONFileStore(jsonPath)
	if err != nil {
		fmt.Println("Erreur d'initialisation du stockage JSON :", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Bienvenue dans le Mini CRM (persistance JSON) ✅")
	app.Run(reader, store)
}
