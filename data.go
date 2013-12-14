package main 

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("Food already exists")
)

// The DB interface contains all the db methods available
// to the Food type. 
type DB interface {
	Get(id int) *Food 
	GetAll() []*Food
	Find(name, description, venue string) []*Food
	Create(f *Food) (int, error)
	Update(f *Food) error
	Delete(id int)
}

// Thread-safe in-memory map of foods.
type foodsDB struct {
	sync.RWMutex
	m 	map[int]*Food 
	seq int 
}

// The database instance.
var db DB 

func init() {
	db = &foodsDB {
		m: make(map[int]*Food),
	}

	// Fill the database
	db.Create(&Food{Id: 1, Name: "Pizza", Description: "Yummm!", Venue: "Domino's Pizza"})
	db.Create(&Food{Id: 2, Name: "Burger", Description: "Wow!", Venue: "Relish"})
	db.Create(&Food{Id: 3, Name: "Pita", Description: "Dang!", Venue: "Pita Pit"})
}

// GetAll returns all foods from the database.
func (db *foodsDB) GetAll() []*Food {
	db.RLock()
	defer db.RUnlock()
	if len(db.m) == 0 {
		return nil
	}
	ar := make([]*Food, len(db.m))
	i := 0
	for _, v := range db.m {
		ar[i] = v
		i++
	}
	return ar 
}

// Find returns all albums that match the search criteria.
func (db *foodsDB) Find(name, description, venue string) []*Food {
	db.RLock()
	defer db.RUnlock()
	var res []*Food 
	for _, v := range db.m {
		if v.Name == name || name == "" {
			if v.Description == description || description == "" {
				if v.Venue == venue || venue == "" {
					res = append(res, v)
				}
			}
		}
	}
	return res 
}

// Get returns the food identified by the id, or nil.
func (db *foodsDB) Get(id int) *Food {
	db.RLock()
	defer db.RUnlock()
	return db.m[id]
}

// Create adds a new food and returns its id, or an error.
func (db *foodsDB) Create(f *Food) (int, error) {
	db.Lock()
	defer db.Unlock()
	// Return an error if Name-Venue already exists
	if !db.isUnique(f) {
		return 0, ErrAlreadyExists
	}
	// Get the unique ID
	db.seq++
	f.Id = db.seq
	// Store
	db.m[f.Id] = f 
	return f.Id, nil
}

// Update changes the food identified by the id. It returns an error if the 
// updated food is a duplicate.
func (db *foodsDB) Update(f *Food) error {
	db.Lock()
	defer db.Unlock()
	if !db.isUnique(f) {
		return ErrAlreadyExists
	}
	db.m[f.Id] = f
	return nil
}

// Delete removes the food identified by the id from the database. It is a no-op
// if the id doesn't exist.
func (db *foodDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()
	delete(db.m, id)
}

// Checks if the food already exists in the database, based on the Name and Venue
// fields.
func (db *foodsDB) isUnique(f *Food) bool {
	for _, v := range db.m {
		if v.Name == f.Name && v.Venue == f.Venue && v.Id != f.Id {
			return false
		}
	}
	return true
}

// the Food data structure, serializable in JSON, XML, and text using the Stringer interface
type Food struct {
	XMLName 		xml.Name 	`json:"-" xml:"food"`
	Id 				int 		`json:"id" xml:"id,attr"`
	Name 			string 		`json:"name" xml:"name"`
	Description 	string 		`json:"description" xml:"description"`
	Venue 			string 		`json:"venue" xml:"venue"`
}

func (f *Food) String() string {
	return fmt.Sprintf("%s: %s -- served by: %s", f.Name, f.Description, f.Venue)
}
