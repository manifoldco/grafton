package connector

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
)

var (
	port         uint
	clientID     = "21jtaatqj8y5t0kctb2ejr6jev5w8"
	clientSecret = "3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw"
	product      = "tester"
)

func makeResource(t *testing.T, plan, region string) *Resource {
	ID, err := manifold.NewID(idtype.Resource)
	if err != nil {
		gm.Expect(err).ToNot(gm.HaveOccurred())
		return nil
	}

	return &Resource{
		ID:        ID,
		Product:   product,
		Plan:      plan,
		Region:    region,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func TestConnector(t *testing.T) {
	gm.RegisterTestingT(t)

	c, err := New(port, clientID, clientSecret, product)
	if err != nil {
		gm.Expect(err).ToNot(gm.HaveOccurred())
	}

	t.Run("a resource is available if added and not if removed", func(t *testing.T) {
		gm.RegisterTestingT(t)

		r := makeResource(t, "high", "aws::us-east-1")
		c.AddResource(r)

		found := c.GetResource(r.ID)
		gm.Expect(found).ToNot(gm.BeNil(), "resource should have been found")

		err := c.RemoveResource(r.ID)
		gm.Expect(err).ToNot(gm.HaveOccurred())

		found = c.GetResource(r.ID)
		gm.Expect(found).To(gm.BeNil())
	})
}
