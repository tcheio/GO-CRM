package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// jsonFileFormat est le format stocké dans le fichier
type jsonFileFormat struct {
	NextID   int        `json:"next_id"`
	Contacts []*Contact `json:"contacts"`
}

type JSONFileStore struct {
	path     string
	contacts map[int]*Contact
	nextID   int
}

// NewJSONFileStore crée un store persistant ; charge si le fichier existe
func NewJSONFileStore(path string) (*JSONFileStore, error) {
	js := &JSONFileStore{
		path:     path,
		contacts: make(map[int]*Contact),
		nextID:   1,
	}
	if err := js.load(); err != nil {
		// On ne bloque pas si fichier absent ; on bloque seulement si autre erreur d'E/S/JSON
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	return js, nil
}

// --- Implémentation Storer ---

func (s *JSONFileStore) Add(contact *Contact) error {
	contact.ID = s.nextID
	s.contacts[contact.ID] = contact
	s.nextID++
	return s.save()
}

func (s *JSONFileStore) GetAll() ([]*Contact, error) {
	out := make([]*Contact, 0, len(s.contacts))
	for _, c := range s.contacts {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *JSONFileStore) GetByID(id int) (*Contact, error) {
	c, ok := s.contacts[id]
	if !ok {
		return nil, errors.New("contact introuvable")
	}
	return c, nil
}

func (s *JSONFileStore) Update(id int, newName, newEmail string) error {
	c, ok := s.contacts[id]
	if !ok {
		return errors.New("contact introuvable")
	}
	if newName != "" {
		c.Name = newName
	}
	if newEmail != "" {
		c.Email = newEmail
	}
	return s.save()
}

func (s *JSONFileStore) Delete(id int) error {
	if _, ok := s.contacts[id]; !ok {
		return errors.New("contact introuvable")
	}
	delete(s.contacts, id)
	return s.save()
}

// --- Persistance ---

func (s *JSONFileStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var payload jsonFileFormat
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	s.contacts = make(map[int]*Contact, len(payload.Contacts))
	for _, c := range payload.Contacts {
		// Clones simples (pointeurs déjà OK)
		s.contacts[c.ID] = c
	}
	if payload.NextID > 0 {
		s.nextID = payload.NextID
	} else {
		// Si pas de nextID, le recalculer
		maxID := 0
		for id := range s.contacts {
			if id > maxID {
				maxID = id
			}
		}
		s.nextID = maxID + 1
	}
	return nil
}

func (s *JSONFileStore) save() error {
	list := make([]*Contact, 0, len(s.contacts))
	for _, c := range s.contacts {
		list = append(list, c)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })

	payload := jsonFileFormat{
		NextID:   s.nextID,
		Contacts: list,
	}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	// Écriture atomique : on écrit d’abord dans un fichier temporaire
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, encoded, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
