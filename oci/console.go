package oci

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/core"
)

func (cc *ClientConfig) CaptureConsole(ocid string) (*core.ConsoleHistory, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(cc.Config())
	if err != nil {
		return nil, err
	}
	req := core.CaptureConsoleHistoryRequest{
		CaptureConsoleHistoryDetails: core.CaptureConsoleHistoryDetails {
			InstanceId: &ocid,
		},
	}

	res, err := client.CaptureConsoleHistory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("console capture error : %v", err)
	}

	return &res.ConsoleHistory, nil
}

func (cc *ClientConfig) CollectConsole(ch *core.ConsoleHistory, len int) (*string, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(cc.Config())

	if err != nil {
		return nil, err
	}

	res, err := client.GetConsoleHistoryContent(context.Background(),
		core.GetConsoleHistoryContentRequest{
			InstanceConsoleHistoryId: ch.Id,
			Length: &len,
		})

	if err != nil {
		return nil, err
	}

	return res.Value, nil

}