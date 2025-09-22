package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Contact struct {
	ID    int
	Name  string
	Email string
}

var contacts = map[int]Contact{}

func main() {
	// Mode FLAGS
	addFlag := flag.Bool("add", false, "Ajouter un contact via flags puis quitter")
	idFlag := flag.Int("id", 0, "ID du contact (obligatoire avec -add)")
	nameFlag := flag.String("name", "", "Nom du contact (obligatoire avec -add)")
	emailFlag := flag.String("email", "", "Email du contact (obligatoire avec -add)")
	flag.Parse()

	if *addFlag {
		if *idFlag == 0 || strings.TrimSpace(*nameFlag) == "" || strings.TrimSpace(*emailFlag) == "" {
			fmt.Println("Erreur: -id, -name et -email sont obligatoires avec -add")
			os.Exit(1)
		}
		if err := validateEmail(*emailFlag); err != nil {
			fmt.Println("Email invalide:", err)
			os.Exit(1)
		}
		if _, exists := contacts[*idFlag]; exists {
			fmt.Printf("Erreur: un contact avec l'ID %d existe déjà.\n", *idFlag)
			os.Exit(1)
		}
		contacts[*idFlag] = Contact{ID: *idFlag, Name: *nameFlag, Email: *emailFlag}
		fmt.Println("Contact ajouté via flags ✅")
		printOne(contacts[*idFlag])
		return
	}

	// MENU
	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
		fmt.Print("> Choix: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			continue
		case "2":
			if err := handleAdd(reader); err != nil {
				fmt.Println("❌", err)
			} else {
				fmt.Println("✅ Contact ajouté.")
			}
		case "3":
			listAll()
		case "4":
			if err := handleDelete(reader); err != nil {
				fmt.Println("❌", err)
			} else {
				fmt.Println("✅ Contact supprimé.")
			}
		case "5":
			if err := handleUpdate(reader); err != nil {
				fmt.Println("❌", err)
			} else {
				fmt.Println("✅ Contact mis à jour.")
			}
		case "6":
			fmt.Println("👋 Au revoir.")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func printMenu() {
	fmt.Println("\n====== Mini-CRM ======")
	fmt.Println("1) Afficher ce menu (boucle)")
	fmt.Println("2) Ajouter un contact (ID, Nom, Email)")
	fmt.Println("3) Lister tous les contacts")
	fmt.Println("4) Supprimer un contact par ID")
	fmt.Println("5) Mettre à jour un contact")
	fmt.Println("6) Quitter")
}

func handleAdd(r *bufio.Reader) error {
	id, err := askInt(r, "ID")
	if err != nil {
		return err
	}
	if _, exists := contacts[id]; exists {
		return fmt.Errorf("un contact avec l'ID %d existe déjà", id)
	}
	name, err := askNonEmpty(r, "Nom")
	if err != nil {
		return err
	}
	email, err := askNonEmpty(r, "Email")
	if err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	contacts[id] = Contact{ID: id, Name: name, Email: email}
	return nil
}

func handleDelete(r *bufio.Reader) error {
	id, err := askInt(r, "ID à supprimer")
	if err != nil {
		return err
	}
	if _, ok := contacts[id]; !ok {
		return fmt.Errorf("aucun contact avec l'ID %d", id)
	}
	delete(contacts, id)
	return nil
}

func handleUpdate(r *bufio.Reader) error {
	id, err := askInt(r, "ID à mettre à jour")
	if err != nil {
		return err
	}
	c, ok := contacts[id]
	if !ok {
		return fmt.Errorf("aucun contact avec l'ID %d", id)
	}

	fmt.Printf("Nom actuel: %s (laisser vide pour conserver)\n", c.Name)
	name, _ := readLine(r, "Nouveau nom")
	fmt.Printf("Email actuel: %s (laisser vide pour conserver)\n", c.Email)
	email, _ := readLine(r, "Nouvel email")

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name != "" {
		c.Name = name
	}
	if email != "" {
		if err := validateEmail(email); err != nil {
			return err
		}
		c.Email = email
	}
	contacts[id] = c
	return nil
}

func listAll() {
	if len(contacts) == 0 {
		fmt.Println("Aucun contact.")
		return
	}
	fmt.Println("\n--- Contacts ---")
	for _, c := range contacts {
		printOne(c)
	}
}

func printOne(c Contact) {
	fmt.Printf("- ID: %d | Nom: %s | Email: %s\n", c.ID, c.Name, c.Email)
}

// Utils

func askInt(r *bufio.Reader, label string) (int, error) {
	s, err := readLine(r, label)
	if err != nil {
		return 0, err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("valeur vide")
	}
	n, convErr := strconv.Atoi(s)
	if convErr != nil {
		return 0, fmt.Errorf("entrez un nombre valide (%v)", convErr)
	}
	return n, nil
}

func askNonEmpty(r *bufio.Reader, label string) (string, error) {
	s, err := readLine(r, label)
	if err != nil {
		return "", err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("valeur vide")
	}
	return s, nil
}

func readLine(r *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s: ", label)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func validateEmail(s string) error {
	s = strings.TrimSpace(s)
	if s == "" || strings.Count(s, "@") != 1 {
		return errors.New("format email invalide")
	}
	parts := strings.Split(s, "@")
	local, domain := parts[0], parts[1]
	if local == "" || domain == "" || strings.Contains(domain, "..") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return errors.New("format email invalide")
	}
	if !strings.Contains(domain, ".") {
		return errors.New("domaine invalide")
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r < 32 {
			return errors.New("caractères invalides dans l'email")
		}
	}
	return nil
}
