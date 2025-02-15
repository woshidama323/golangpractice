package filecoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	lotusbig "github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/chain/types"
	builtin2 "github.com/filecoin-project/specs-actors/v7/actors/builtin"
	"github.com/urfave/cli"
)

type GetAddressByMinerID struct {
	ChainProvider struct {
		Api    v0api.FullNodeStruct
		Closer jsonrpc.ClientCloser
	}
}

func NewMiner(providerUrl string) (*GetAddressByMinerID, error) {

	var res v0api.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), providerUrl, "Filecoin",
		[]interface{}{
			&res.CommonStruct.Internal,
			&res.Internal,
		}, nil)

	if err != nil {
		return nil, err
	}
	return &GetAddressByMinerID{
		ChainProvider: struct {
			Api    v0api.FullNodeStruct
			Closer jsonrpc.ClientCloser
		}{
			Api:    res,
			Closer: closer,
		},
	}, nil
}

func GetMinerInfo(urlstr, minerid string) ([]string, error) {
	getMinerInfo := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.StateMinerInfo",
		"params": []interface{}{
			minerid,
			nil,
		},
		"id": 1,
	}

	// res := &HeadInfo{}
	value, _ := json.Marshal(getMinerInfo)

	resp, err := http.Post(urlstr, "application/json", bytes.NewBuffer(value))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var FResp ResultDataType
	err = json.Unmarshal(body, &FResp)
	if err != nil {
		return nil, err
	}
	// if FResp.Error.IsError() {
	// 	return nil, FResp.Error.Error()
	// }
	// if FResp.Result == nil {
	// 	return nil, errors.New("Result is nil ")
	// }
	// _d, err := json.Marshal(FResp.Result)
	// if err != nil {
	// 	return nil, err
	// }

	// err = json.Unmarshal(_d, res)

	fmt.Printf("current result is :%v", FResp)
	fmt.Printf("current result body is :%v", string(body))

	listaddresses := []string{}
	listaddresses = append(listaddresses, FResp.Result.ControlAddresses...)
	listaddresses = append(listaddresses, FResp.Result.Multiaddrs...)
	listaddresses = append(listaddresses, FResp.Result.Owner, FResp.Result.Worker)
	return listaddresses, nil
}

func GetAddressInfo(urlstr string, address string) (string, error) {

	getAddresses := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.StateAccountKey",
		"params": []interface{}{
			address,
			nil,
		},
		"id": 1,
	}

	// res := &HeadInfo{}
	value, _ := json.Marshal(getAddresses)

	resp, err := http.Post(urlstr, "application/json", bytes.NewBuffer(value))
	if err != nil {

		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

type ResType struct {
}
type MinerInfo struct {
	Owner                      string
	Worker                     string
	NewWorker                  string
	ControlAddresses           []string
	WorkerChangeEpoch          uint64
	PeerId                     *string
	Multiaddrs                 []string
	WindowPoStProofType        int
	SectorSize                 int
	WindowPoStPartitionSectors uint64
	ConsensusFaultElapsed      int
}

type ResultDataType struct {
	JSONRpc string `json:"jsonrpc"`
	Result  struct {
		Owner                      string
		Worker                     string
		NewWorker                  string
		ControlAddresses           []string
		WorkerChangeEpoch          int64
		PeerId                     string
		Multiaddrs                 []string
		WindowPoStProofType        int64
		SectorSize                 int64
		WindowPoStPartitionSectors int64
		ConsensusFaultElapsed      int64
	}
	ID int64 `json:"id"`
}

//test wallet
var walletCmd = &cli.Command{
	Name:  "wallet",
	Usage: "wallet operating",
	Action: func(c *cli.Context) error {

		// filecoinwallet.CreateWallet()
		// actor.TestActor()

		type Version uint32
		major := 1
		minor := 3
		patch := 0
		x := Version(uint32(major)<<16 | uint32(minor)<<8 | uint32(patch))
		fmt.Println("....", x)

		v := uint32(x)

		// major := (v & majorOnlyMask) >> 16)
		fmt.Println("major....", (v&majorOnlyMask)>>16)

		//测算cid的创建过程

		m := &types.Message{
			To:    builtin2.StoragePowerActorAddr,
			From:  builtin2.SystemActorAddr,
			Nonce: 34,
			Value: lotusbig.Zero(),

			GasLimit:   123,
			GasFeeCap:  lotusbig.NewInt(234),
			GasPremium: lotusbig.NewInt(234),

			Method: 6,
			Params: []byte("hai"),
		}

		fmt.Printf("current cid: %+v", m.Cid().String())
		return nil
	},
}
