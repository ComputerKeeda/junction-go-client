package core

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
)

func GetVRF() {

	stationId, err := chain.GetStationId()
	if err != nil {
		components.Logger.Error(err.Error())
		return
	}
	podNumber, err := chain.GetPodNumber()
	if err != nil {
		components.Logger.Error(err.Error())
		return
	}

	qClient := components.GetQueryClient()
	ctx := context.Background()
	queryResp, err := qClient.FetchVrn(ctx, &types.QueryFetchVrnRequest{
		PodNumber: podNumber,
		StationId: stationId,
	})
	if err != nil {
		fmt.Println("Error fetching VRF: ", err)
		return
	}

	components.Logger.Info(fmt.Sprintf("VRF: %s", queryResp.Details))
}
