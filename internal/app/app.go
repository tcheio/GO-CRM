package app

import (
	"bufio"
	"fmt"
	"refactor_crm_interface/internal/storage"
	"strconv"
	"strings"
)

func Run(reader *bufio.Reader, store storage.Storer) {
	for {
		printMenu()
		fmt.Print("Votre choix : ")

		choice := readUserChoice(reader)

		switch choice {
		case 1:
			handleAddContact(reader, store)
		case 2:
			handleListContacts(store)
		case 3:
			handleUpdateContact(reader, store)
		case 4:
			handleDeleteContact(reader, store)
		case 5:
			fmt.Println("👋 Au revoir !")
			return
		default:
			fmt.Println("Option invalide, veuillez réessayer.")
		}
	}
}

func printMenu() {
	fmt.Println("\n--- Menu Principal ---")
	fmt.Println("1. Ajouter un contact")
	fmt.Println("2. Lister les contacts")
	fmt.Println("3. Mettre à jour un contact")
	fmt.Println("4. Supprimer un contact")
	fmt.Println("5. Quitter")
}

// --- Handlers ---

func handleAddContact(reader *bufio.Reader, storer storage.Storer) {
	fmt.Print("Entrez le nom du contact : ")
	name := readLine(reader)

	fmt.Print("Entrez l'email du contact : ")
	email := readLine(reader)

	contact := &storage.Contact{
		Name:  name,
		Email: email,
	}
	err := storer.Add(contact)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	fmt.Printf("Contact '%s' ajouté avec l'ID %d.\n", contact.Name, contact.ID)
}

func handleListContacts(store storage.Storer) {
	contacts, err := store.GetAll()
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}

	if len(contacts) == 0 {
		fmt.Println("Aucun contact à afficher.")
		return
	}

	fmt.Println("\n--- Liste des Contacts ---")
	for _, contact := range contacts {
		fmt.Printf("ID: %d, Nom: %s, Email: %s\n", contact.ID, contact.Name, contact.Email)
	}
}

func handleUpdateContact(reader *bufio.Reader, store storage.Storer) {
	fmt.Print("Entrez l'ID du contact à mettre à jour : ")
	id := readInteger(reader)
	if id == -1 {
		return
	}

	existingContact, err := store.GetByID(id)
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}

	fmt.Printf("Mise à jour de '%s'. Laissez vide pour conserver la valeur actuelle.\n", existingContact.Name)

	fmt.Printf("Nouveau nom (%s) : ", existingContact.Name)
	newName := readLine(reader)

	fmt.Printf("Nouvel email (%s) : ", existingContact.Email)
	newEmail := readLine(reader)

	err = store.Update(id, newName, newEmail)
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}

	fmt.Println("Contact mis à jour avec succès ✅")
}

func handleDeleteContact(reader *bufio.Reader, store storage.Storer) {
	fmt.Print("Entrez l'ID du contact à supprimer : ")
	id := readInteger(reader)
	if id == -1 {
		return
	}

	err := store.Delete(id)
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
		return
	}

	fmt.Printf("Le contact avec l'ID %d a été supprimé.\n", id)
}

// --- Utils ---

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readUserChoice(reader *bufio.Reader) int {
	choice, err := strconv.Atoi(readLine(reader))
	if err != nil {
		return -1
	}
	return choice
}

func readInteger(reader *bufio.Reader) int {
	id, err := strconv.Atoi(readLine(reader))
	if err != nil {
		fmt.Println("Erreur : ID invalide. Veuillez entrer un nombre.")
		return -1
	}
	return id
}
