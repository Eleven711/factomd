package wsapi

import (
	"fmt"
	"testing"
	"github.com/FactomProject/factomd/common/primitives"
	//"encoding/json"
)

type AddressRequest struct {
	Address string `json:"address"`
}

func TestHandleFactoidBalance(t *testing.T) {

	jsonStr := []byte(`{"jsonrpc":"2.0","id":0,"params":{"Address":"FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"},"method":"factoid-balance"}`)

	j, err := primitives.ParseJSON2Request(string(jsonStr))
	t.Log(j)
	params := j.Params
	
	if err != nil {
		t.Log("%s",err)
	}

	t.Log("params: %s",params)

	HandleV2FactoidBalance(params,t)	
}

func HandleV2FactoidBalance(params interface{},t *testing.T)  {
	strs := params.(map[string]interface {}) 
	fadr := strs["Address"]

	fmt.Println(fadr)	
}