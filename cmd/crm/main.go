package main

import (
	"bufio"
	"fmt"
	"os"

	"refactor_crm_interface/internal/app"
	"refactor_crm_interface/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Mini CRM v3!")

	app.Run(reader, store)
}
