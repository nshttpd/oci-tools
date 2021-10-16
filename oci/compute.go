package oci

import (
	"context"

	"runtime"

	"sync"

	"fmt"

	"github.com/oracle/oci-go-sdk/v49/core"
	"github.com/oracle/oci-go-sdk/v49/identity"
)

// a conglomeration of things that make up a OCI Instance, the instance itself, VNICs
// and other things later as needed (attached storage etc. etc.)
type Compute struct {
	Instance    core.Instance
	Vnics       []core.Vnic
	Compartment identity.Compartment
}

func (cc *ClientConfig) GetComputeInstances(compartment identity.Compartment) ([]*Compute, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(cc.Config())
	client.SetRegion(cc.Region)
	if err != nil {
		// when real logging is here DEBUG stuff can go here.
		return nil, err
	}

	var more = true
	var next *string

	instances := make([]core.Instance, 0)

	for more {
		req := core.ListInstancesRequest{CompartmentId: compartment.Id}
		if next != nil {
			req.Page = next
		}
		res, err := client.ListInstances(context.Background(), req)
		if err == nil {
			instances = append(instances, res.Items...)
			if res.OpcNextPage != nil {
				next = res.OpcNextPage
			} else {
				more = false
			}
		} else {
			more = false
			return nil, err
		}
	}

	throttle := make(chan int, runtime.NumCPU())
	var wg sync.WaitGroup
	mux := &sync.Mutex{}

	computes := make([]*Compute, len(instances))
	for x, i := range instances {

		throttle <- 1
		wg.Add(1)
		go func(idx int, inst core.Instance) {
			defer wg.Done()
			var v []core.Vnic
			if inst.LifecycleState == core.InstanceLifecycleStateRunning {
				v, err = cc.GetVnics(inst)
				if err != nil {
					fmt.Printf("error fetching vnic : %v\n", err)
					<-throttle
					return
				}
			}
			mux.Lock()
			computes[idx] = &Compute{Instance: inst, Vnics: v, Compartment: compartment}
			mux.Unlock()
			<-throttle
		}(x, i)
	}
	wg.Wait()
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
