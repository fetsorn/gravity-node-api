package client

import (
	"encoding/json"
	"fmt"
	coreconfig "github.com/Gravity-Tech/gravity-core/config"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"net/http"
)

type LedgerClient struct {
	EndpointURL string
}

func NewLedgerClient(endpoint string) *LedgerClient {
	return &LedgerClient{ EndpointURL:endpoint }
}

type responseResult struct {
	Response abcitypes.ResponseQuery `json:"response"`
}
type fetchResponse struct {
	Result *responseResult `json:"result"`
}

func (ledger *LedgerClient) extractData(response fetchResponse) string {
	return string(response.Result.Response.Value)
}

func (ledger *LedgerClient) extractValidatorInfo(response *ctypes.ResultStatus) *ctypes.ValidatorInfo {
	return &response.ValidatorInfo
}

func (ledger *LedgerClient) FetchValidatorStatus() (*ctypes.ResultStatus, error) {
	path := fmt.Sprintf("%v/status", ledger.EndpointURL)

	response, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	byteValue, err := ioutil.ReadAll(response.Body)

	var resultStatus ctypes.ResultStatus

	err = json.Unmarshal(byteValue, &resultStatus)

	if err != nil {
		return nil, err
	}

	return &resultStatus, nil
}

func (ledger *LedgerClient) FetchValidatorDetails() (*coreconfig.ValidatorDetails, error) {
	path := fmt.Sprintf("%v/abci_query?path=\"validatorDetails\"", ledger.EndpointURL)

	response, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	byteValue, err := ioutil.ReadAll(response.Body)

	var parsedResponse *fetchResponse
	err = json.Unmarshal(byteValue, &parsedResponse)

	if err != nil {
		return nil, err
	}

	ledgerDataExtracted := ledger.extractData(*parsedResponse)

	var validatorDetails coreconfig.ValidatorDetails
	err = json.Unmarshal([]byte(ledgerDataExtracted), &validatorDetails)

	if err != nil {
		return nil, err
	}

	return &validatorDetails, nil
}