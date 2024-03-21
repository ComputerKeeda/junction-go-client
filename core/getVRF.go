package core

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
)

func GetVRF() *types.VrfRecord {

	stationId, err := chain.GetStationId()
	if err != nil {
		components.Logger.Error(err.Error())
		return nil
	}
	podNumber, err := chain.GetPodNumber()
	if err != nil {
		components.Logger.Error(err.Error())
		return nil
	}

	qClient := components.GetQueryClient()
	ctx := context.Background()
	queryResp, err := qClient.FetchVrn(ctx, &types.QueryFetchVrnRequest{
		PodNumber: podNumber,
		StationId: stationId,
	})
	if err != nil {
		fmt.Println("Error fetching VRF: ", err)
		return nil
	}

	components.Logger.Info(fmt.Sprintf("VRF: %s", queryResp.Details)) // types.VrfRecord

	return queryResp.Details
}
