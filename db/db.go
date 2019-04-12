package db

import (
	"time"

	manifold "github.com/manifoldco/go-manifold"
)

// DB functions as an in-memory database of marketplace entities
type DB struct {
	ResourcesByID         map[manifold.ID]Resource
	CredentialsByResource map[manifold.ID][]Credential
	CredentialsByID       map[manifold.ID]Credential
	MeasuresByResource    map[manifold.ID][]Measure
}

// New creates a new DB instance in-memory
func New() *DB {
	return &DB{
		ResourcesByID:         make(map[manifold.ID]Resource),
		CredentialsByResource: make(map[manifold.ID][]Credential),
		CredentialsByID:       make(map[manifold.ID]Credential),
		MeasuresByResource:    make(map[manifold.ID][]Measure),
	}
}

// PutResource stores the provided resource in the database
func (db *DB) PutResource(r Resource) {
	db.ResourcesByID[r.ID] = r
}

// GetResource returns a resource based on it's id or nil, if it can't be found
func (db *DB) GetResource(id manifold.ID) *Resource {
	r, ok := db.ResourcesByID[id]
	if ok {
		return &r
	}
	return nil
}

// DeleteResource removes a resource and returns true, false if there was no resource
func (db *DB) DeleteResource(id manifold.ID) bool {
	_, ok := db.ResourcesByID[id]
	if ok {
		cs := db.CredentialsByResource[id]
		for _, c := range cs {
			db.DeleteCredential(c.ID)
		}
		delete(db.ResourcesByID, id)
		return true
	}
	return false
}

// PutCredential stores the provided credential, it must have a ResourceID set!
func (db *DB) PutCredential(c Credential) {
	if c.ResourceID.IsEmpty() {
		panic("Supplied credential did not have a resource ID specified")
	}
	db.CredentialsByID[c.ID] = c
	db.CredentialsByResource[c.ResourceID] = append(
		db.CredentialsByResource[c.ResourceID], c)
}

// GetCredential returns a credential based on it's id or nil, if it can't be found
func (db *DB) GetCredential(id manifold.ID) *Credential {
	c, ok := db.CredentialsByID[id]
	if ok {
		return &c
	}
	return nil
}

// GetCredentialsByResource returns a list of credentials or nil, for a ResourceID
func (db *DB) GetCredentialsByResource(id manifold.ID) []Credential {
	c, ok := db.CredentialsByResource[id]
	if ok {
		return c
	}
	return nil
}

// DeleteCredential removes a credential and returns true, false if there was no credential
func (db *DB) DeleteCredential(id manifold.ID) bool {
	c, ok := db.CredentialsByID[id]
	if !ok {
		return false
	}
	cs, ok := db.CredentialsByResource[c.ResourceID]
	if !ok {
		panic("Credential existed, but resource credential list did not")
	}
	// Remove from the resource's list
	for i := len(cs) - 1; i >= 0; i-- {
		if cs[i].ID == c.ID {
			cs = append(cs[:i], cs[i+1:]...)
		}
	}
	db.CredentialsByResource[c.ResourceID] = cs
	// Finally remove the key
	delete(db.CredentialsByID, id)
	return true
}

// PutMeasure stores the provided measure, it must have a ResourceID set!
func (db *DB) PutMeasure(m Measure) {
	if m.ResourceID.IsEmpty() {
		panic("Supplied measure did not have a resource ID specified")
	}
	m.UpdatedAt = time.Now()
	db.MeasuresByResource[m.ResourceID] = append(
		db.MeasuresByResource[m.ResourceID], m)
}

// GetMeasuresByResource returns a list of measures or nil, for a ResourceID
func (db *DB) GetMeasuresByResource(id manifold.ID) []Measure {
	m, ok := db.MeasuresByResource[id]
	if ok {
		return m
	}
	return nil
}
