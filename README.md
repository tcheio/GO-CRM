# 📇 Mini-CRM (CLI) en Go - CHEIO THOMAS

Un petit CRM en ligne de commande développé en **Go**.  
Il permet de gérer une liste de contacts en mémoire : **ajout, suppression, mise à jour et affichage**.

---

## 🚀 Lancer le projet

Placez-vous dans le dossier où se trouvent tous vos fichiers `.go` puis exécutez :

```bash
go run .
```

---

## 📖 Fonctionnalités

- Mode interactif (menu en boucle)
- Afficher le menu principal
- Ajouter un contact (**ID, Nom, Email**)
- Lister tous les contacts
- Supprimer un contact par ID
- Mettre à jour un contact
- Quitter l’application

---

## ⚡ Utilisation avec flags (ajout direct)

Il est possible d’ajouter un contact directement sans passer par le menu grâce aux **flags** :

```bash
go run . -add -id 1 -name "Alice" -email "alice@example.com"
```

### Options disponibles :
- `-add` → active le mode ajout par flags
- `-id` → identifiant numérique du contact
- `-name` → nom du contact
- `-email` → email du contact

⚠️ **Tous ces champs sont obligatoires** lorsque vous utilisez `-add`.

### Exemples :

```bash
# Ajout d’un contact Alice
go run . -add -id 1 -name "Alice" -email "alice@example.com"

# Ajout d’un contact Bob
go run . -add -id 2 -name "Bob" -email "bob@test.org"
```

---

## 🛠️ Concepts utilisés

Le projet met en pratique plusieurs idiomes Go :

- `for {}` → boucle infinie du menu
- `switch` → gestion des choix du menu
- `map[int]Contact` → stockage des contacts
- *comma ok idiom* → vérification de la présence d’une clé dans la map
- `if err != nil` → gestion des erreurs
- `strconv` → conversion des entrées utilisateur
- `os.Stdin` + `bufio.NewReader` → lecture de la saisie utilisateur
- `flag` → gestion des options en ligne de commande

---

## 🧩 Améliorations possibles

- Sauvegarder/charger les contacts depuis un fichier JSON
- Ajouter des sous-commandes (`mini-crm add | list | update | delete`)
- Écrire des tests unitaires (`go test`) pour valider la logique

---
