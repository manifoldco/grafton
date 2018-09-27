package db

import (
	"time"

	"github.com/manifoldco/go-manifold"
)

// ResourceState defines the current state of the resource
type ResourceState string

const (
	// ResourceStateProvisioning defines the state for the resource when its "provisioning"
	ResourceStateProvisioning ResourceState = "provisioning"
	// ResourceStateProvisioned defines the state for the resource when its "provisioned"
	ResourceStateProvisioned ResourceState = "provisioned"
	// ResourceStateProvisionFailed defines the state for the resource when it has failed provisioning
	ResourceStateProvisionFailed ResourceState = "provision-failed"
	// ResourceStateDerovisioning defines the state for the resource when its "deprovisioning"
	ResourceStateDerovisioning ResourceState = "deprovisioning"
	// ResourceStateDeprovisioned defines the state for the resource when its "deprovisioned"
	ResourceStateDeprovisioned ResourceState = "deprovisioned"
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
	// Internal Fields
	State ResourceState `json:"-"`
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
