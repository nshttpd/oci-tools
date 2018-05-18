package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/core"
)

// a conglomeration of things that make up a OCI Instance, the instance itself, VNICs
// and other things later as needed (attached storage etc. etc.)
type Compute struct {
	instance core.Instance
	vnics    []core.Vnic
}

func (cc *ClientConfig) GetComputeInstances(cid *string) ([]Compute, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(cc.Config())
	if err != nil {
		// when real logging is here DEBUG stuff can go here.
		return nil, err
	}
	req := core.ListInstancesRequest{CompartmentId: cid}
	res, err := client.ListInstances(context.Background(), req)
	if err != nil {
		return nil, err
	}
	computes := make([]Compute, len(res.Items))
	for x, i := range res.Items {
		v, err := cc.GetVnics(i)
		if err != nil {
			return nil, err
		}
		computes[x] = Compute{instance: i, vnics: v}
	}
	return computes, nil
}

func (cc *ClientConfig) GetVnics(i core.Instance) ([]core.Vnic, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(cc.Config())
	if err != nil {
		return nil, err
	}
	req := core.ListVnicAttachmentsRequest{CompartmentId: i.CompartmentId, InstanceId: i.Id}
	v, err := client.ListVnicAttachments(context.Background(), req)
	if err != nil {
		return nil, err
	}
	vnics := make([]core.Vnic, len(v.Items))

	vcnclient, err := core.NewVirtualNetworkClientWithConfigurationProvider(cc.Config())
	if err != nil {
		return nil, err
	}

	for i, va := range v.Items {
		req := core.GetVnicRequest{VnicId: va.VnicId}
		res, err := vcnclient.GetVnic(context.Background(), req)
		if err != nil {
			return nil, err
		}
		vnics[i] = res.Vnic
	}

	return vnics, nil
}

func (c Compute) Instance() core.Instance {
	return c.instance
}

func (c Compute) Vnics() []core.Vnic {
	return c.vnics
}
