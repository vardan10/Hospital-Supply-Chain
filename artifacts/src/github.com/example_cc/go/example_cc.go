/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main


import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("example_cc0")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	hospitalName   string `json:"hospitalName"`
	assetRequested  string `json:"assetRequested"`
	volunteerName string `json:"volunteerName"`
	volunteerSuccess bool `json:"volunteerSuccess"`
	ngoName  string `json:"owner"`
	ngoSuccess bool `json:"ngoSuccess"`
	hospitalSuccess bool `json:"hospitalSuccess"`


}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "hospitalRequest" {
		return t.hospitalRequest(stub, args)
	} else if function == "getHospitalRequest" {
		return t.getHospitalRequest(stub,args)
	} else if function == "volunteerRequest" {
		return t.volunteerRequest(stub,args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (t *SimpleChaincode) hospitalRequest(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var request = SimpleChaincode{hospitalName: args[1], assetRequested: args[2], volunteerName: args[3]}

	ledgerRequest, _ := json.Marshal(request)
	APIstub.PutState(args[0], ledgerRequest)

	return shim.Success(nil)
}

func (t *SimpleChaincode) getHospitalRequest(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	returnResponce,_ := APIstub.GetState(args[0])

	jsonResp := "{\"Name\":\"Vardan\",\"Amount\":\"Test\"}"
	return shim.Success(jsonResp)
}

func (t *SimpleChaincode) volunteerRequest(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	return shim.Success(nil)
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
