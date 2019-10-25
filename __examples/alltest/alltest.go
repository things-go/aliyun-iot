package main

import (
	"encoding/json"
	"fmt"

	"github.com/thinkgos/aliIOT/model"
)

func main() {
	rsp := model.GwSubDevRegisterResponse{}
	if err := json.Unmarshal([]byte(`{
		  "id": "123",
		  "code": 200,
		  "data": [
		    {
		      "iotId": "12344",
		      "productKey": "1234556554",
		      "deviceName": "deviceName1234",
		      "deviceSecret": "xxxxxx"
		    },
		    {
		      "iotId": "12344",
		      "productKey": "1234556554",
		      "deviceName": "deviceName1234",
		      "deviceSecret": "xxxxxx"
		    }
		  ]
	}`), &rsp); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", rsp)
}

type UserProc struct {
	model.GwNopUserProc
}

func (UserProc) DownstreamGwExtSubDevRegisterReply(m *model.Manager, rsp *model.GwSubDevRegisterResponse) error {

}
