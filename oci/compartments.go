package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/identity"
)

type Compartments struct {
	compartments []identity.Compartment
}

func (c *ClientConfig) GetCompartments() (Compartments, error) {
	cs := Compartments{}
	var err error

	client, err := identity.NewIdentityClientWithConfigurationProvider(c.Config())
	if err == nil {
		cid, err := c.TenancyOCID()
		if err == nil {
			more := true
			var next *string

			for more {
				req := identity.ListCompartmentsRequest{CompartmentId: &cid}
				if next != nil {
					req.Page = next
				}
				res, err := client.ListCompartments(context.Background(), req)
				if err == nil {
					cs.compartments = append(cs.compartments, res.Items...)
					if res.OpcNextPage != nil {
						next = res.OpcNextPage
					} else {
						more = false
					}
				} else {
					more = false
				}

			}
		}
	}

	return cs, err
}

func (c *Compartments) CompartmentId(cname string) *string {
	for _, comp := range c.compartments {
		if strings.EqualFold(*comp.Name, cname) {
			return comp.Id
		}
	}
	return nil
}

func (c *Compartments) CompartmentIds() []*string {
	r := make([]*string, len(c.compartments))
	for i, c := range c.compartments {
		r[i] = c.Id
	}
	return r
}

func (c *Compartments) Compartment(cname string) *string {
	for _, comp := range c.compartments {
		if strings.EqualFold(*comp.Name, cname) {
			return comp.Name
		}
	}
	return nil
}

func (c *Compartments) Compartments() []identity.Compartment {
	return c.compartments
}
