package db

import (
	"time"

	manifold "github.com/manifoldco/go-manifold"
)

// DB functions as an in-memory database of marketplace entities
type DB struct {
	resourcesByID         map[manifold.ID]Resource
	credentialsByResource map[manifold.ID][]Credential
	credentialsByID       map[manifold.ID]Credential
	measuresByResource    map[manifold.ID][]Measure
}

// New creates a new DB instance in-memory
func New() *DB {
	return &DB{
		resourcesByID:         make(map[manifold.ID]Resource),
		credentialsByResource: make(map[manifold.ID][]Credential),
		credentialsByID:       make(map[manifold.ID]Credential),
		measuresByResource:    make(map[manifold.ID][]Measure),
	}
}

// PutResource stores the provided resource in the database
func (db *DB) PutResource(r Resource) {
	db.resourcesByID[r.ID] = r
}

// GetResource returns a resource based on it's id or nil, if it can't be found
func (db *DB) GetResource(id manifold.ID) *Resource {
	r, ok := db.resourcesByID[id]
	if ok {
		return &r
	}
	return nil
}

// DeleteResource removes a resource and returns true, false if there was no resource
func (db *DB) DeleteResource(id manifold.ID) bool {
	_, ok := db.resourcesByID[id]
	if ok {
		cs, _ := db.credentialsByResource[id]
		for _, c := range cs {
			db.DeleteCredential(c.ID)
		}
		delete(db.resourcesByID, id)
		return true
	}
	return false
}

// PutCredential stores the provided credential, it must have a ResourceID set!
func (db *DB) PutCredential(c Credential) {
	if c.ResourceID.IsEmpty() {
		panic("Supplied credential did not have a resource ID specified")
	}
	db.credentialsByID[c.ID] = c
	db.credentialsByResource[c.ResourceID] = append(
		db.credentialsByResource[c.ResourceID], c)
}

// GetCredential returns a credential based on it's id or nil, if it can't be found
func (db *DB) GetCredential(id manifold.ID) *Credential {
	c, ok := db.credentialsByID[id]
	if ok {
		return &c
	}
	return nil
}

// GetCredentialsByResource returns a list of credentials or nil, for a ResourceID
func (db *DB) GetCredentialsByResource(id manifold.ID) []Credential {
	c, ok := db.credentialsByResource[id]
	if ok {
		return c
	}
	return nil
}

// DeleteCredential removes a credential and returns true, false if there was no credential
func (db *DB) DeleteCredential(id manifold.ID) bool {
	c, ok := db.credentialsByID[id]
	if !ok {
		return false
	}
	cs, ok := db.credentialsByResource[c.ResourceID]
	if !ok {
		panic("Credential existed, but resource credential list did not")
	}
	// Remove from the resource's list
	for i := len(cs) - 1; i >= 0; i-- {
		if cs[i].ID == c.ID {
			cs = append(cs[:i], cs[i+1:]...)
		}
	}
	db.credentialsByResource[c.ResourceID] = cs
	// Finally remove the key
	delete(db.credentialsByID, id)
	return true
}

// PutMeasure stores the provided measure, it must have a ResourceID set!
func (db *DB) PutMeasure(m Measure) {
	if m.ResourceID.IsEmpty() {
		panic("Supplied measure did not have a resource ID specified")
	}
	m.UpdatedAt = time.Now()
	db.measuresByResource[m.ResourceID] = append(
		db.measuresByResource[m.ResourceID], m)
}

// GetMeasuresByResource returns a list of measures or nil, for a ResourceID
func (db *DB) GetMeasuresByResource(id manifold.ID) []Measure {
	m, ok := db.measuresByResource[id]
	if ok {
		return m
	}
	return nil
}
