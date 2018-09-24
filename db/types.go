package db

import (
	"time"

	"github.com/manifoldco/go-manifold"
)

// Resource represents a resource provisioned through Grafton
type Resource struct {
	ID        manifold.ID         `json:"id"`
	Name      manifold.Name       `json:"name"`
	Label     manifold.Label      `json:"label"`
	Plan      manifold.Label      `json:"plan"`
	Product   manifold.Label      `json:"product"`
	Region    string              `json:"region"`
	Features  manifold.FeatureMap `json:"features,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// Credential represents a credential set for a resource
type Credential struct {
	ID          manifold.ID       `json:"id"`
	Keys        map[string]string `json:"keys"`
	CustomNames map[string]string `json:"custom_names"`
	CreatedOn   time.Time         `json:"created_on"`
	// Internal Fields
	ResourceID manifold.ID `json:"-"`
}

// Measure represents a measures set for a resource
type Measure struct {
	ResourceID  manifold.ID      `json:"resource_id"`
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
	Measures    map[string]int64 `json:"measures"`
	// Internal Fields
	UpdatedAt time.Time `json:"-"`
}
