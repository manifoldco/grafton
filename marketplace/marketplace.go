package marketplace

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-zoo/bone"

	"github.com/manifoldco/grafton"
	"github.com/manifoldco/grafton/connector"
	"github.com/manifoldco/grafton/db"
	"github.com/manifoldco/grafton/marketplace/routes"
)

// FakeMarketplace represents a fake marketplace dashboard and backend server
//  to allow provider to test a full integration experience without integrating
type FakeMarketplace struct {
	Port      uint
	DB        *db.DB
	Connector *connector.FakeConnector
	GC        *grafton.Client
	Server    *http.Server
}

// New creates a new FakeMarketplace based on the passed parameters
func New(connector *connector.FakeConnector, port uint, pAPI *url.URL,
	signer grafton.Signer) *FakeMarketplace {

	fm := &FakeMarketplace{
		Port:      port,
		DB:        connector.DB,
		Connector: connector,
	}

	cURL, err := url.Parse("http://localhost:" + strconv.Itoa(int(connector.Config.Port)))
	if err != nil {
		panic("Failed to parse connector url for Grafton client: " + err.Error())
	}
	fm.GC = grafton.New(pAPI, cURL, signer, nil)

	return fm
}

// StartSync starts the server or returns an error if it couldn't be started
func (m *FakeMarketplace) StartSync() error {
	m.Server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", m.Port),
		Handler: Routes(m),
	}

	return m.Server.ListenAndServe()
}

// Start the server or return an error if it couldn't be started
func (m *FakeMarketplace) Start() {
	m.Server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", m.Port),
		Handler: Routes(m),
	}

	go m.Server.ListenAndServe()
}

// Stop the server or return an error if it couldn't be stopped
func (m *FakeMarketplace) Stop() error {
	if m.Server == nil {
		return errors.New("Cannot not stop a server that has not started")
	}

	return m.Server.Close()
}

// Routes returns all the routes for the HTTP Server as a bone.Mux
func Routes(m *FakeMarketplace) *bone.Mux {
	mux := bone.New()
	mux.GetFunc("/", routes.GetResourcesHandler(m.DB))

	mux.GetFunc("/resources", routes.GetResourcesHandler(m.DB))
	mux.PostFunc("/resources", routes.PostResourcesHandler(m.DB, m.GC, m.Connector))
	mux.PutFunc("/resources/:id", routes.PutResourcesHandler(m.DB))
	mux.DeleteFunc("/resources/:id", routes.DeleteResourcesHandler(m.DB, m.GC, m.Connector))

	// TODO: Core funcs
	// mux.GetFunc("/resources/:id/sso", getResourcesSSOHandler(c))
	// mux.GetFunc("/resources/:id/measures", getResourcesMeasuresHandler(c))

	// TODO: Future funcs
	// mux.GetFunc("/users", getUsersHandler(c))
	// mux.PostFunc("/users", postUsersHandler(c))
	// mux.PutFunc("/users/:id", postUsersHandler(c))
	// mux.DeleteFunc("/users/:id", deleteUsersHandler(c))

	// mux.GetFunc("/teams", getTeamsHandler(c))
	// mux.PostFunc("/teams", postTeamsHandler(c))
	// mux.PutFunc("/teams/:id", postTeamsHandler(c))
	// mux.DeleteFunc("/teams/:id", deleteTeamsHandler(c))

	// mux.GetFunc("/projects", getProjectsHandler(c))
	// mux.PostFunc("/projects", postProjectsHandler(c))
	// mux.PutFunc("/projects/:id", postProjectsHandler(c))
	// mux.DeleteFunc("/projects/:id", deleteProjectsHandler(c))

	return mux
}
