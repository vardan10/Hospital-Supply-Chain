package main

import (
	"encoding/json"
	"fmt"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type marble struct {
	Name       			string `json:"name"`			//the fieldtags are needed to keep case from bouncing around
	HospitalName		string `json:"hospitalName"`
	AssetRequested      string `json:"assetRequested"`

	NgoName  			string `json:"owner"`
	NgoSuccess 			bool `json:"ngoSuccess"`
	
	VolunteerName 		string `json:"volunteerName"`
	VolunteerSuccess 	bool `json:"volunteerSuccess"`

	HospitalSuccess 	bool `json:"hospitalSuccess"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "hospitalInvoke" {
		return t.hospitalInvoke(stub, args)
	} else if function == "query" {
		return t.query(stub,args)
	} else if function == "ngoInvoke" {
		return t.ngoInvoke(stub,args)
	} else if function == "volunteerInvoke" {
		return t.volunteerInvoke(stub,args)
	} else if function == "hospitalSuccess" {
		return t.hospitalSuccess(stub,args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// hospitalInvoke - Hospital Creates a request
// ============================================================
func (t *SimpleChaincode) hospitalInvoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       	1       		2     				3
	// "Key", "HospitalName", "AssetRequested", "VolunteerName"
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	request := marble{HospitalName: args[1], AssetRequested: args[2], NgoName: args[3]}

	requestJSONasBytes, err := json.Marshal(request)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(args[0], requestJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ================================================
// query - Get all records
// ================================================
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) ngoInvoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       	1
	// "key", "VoluenteerName"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	requestAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	requestToTransfer := marble{}
	err = json.Unmarshal(requestAsBytes, &requestToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	requestToTransfer.VolunteerName = args[1] //add Volunteer Name
	requestToTransfer.NgoSuccess = true

	marbleJSONasBytes, _ := json.Marshal(requestToTransfer)
	err = stub.PutState(args[0], marbleJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) volunteerInvoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	//   0
	// "key"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	requestAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	requestToTransfer := marble{}
	err = json.Unmarshal(requestAsBytes, &requestToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	requestToTransfer.VolunteerSuccess = true

	marbleJSONasBytes, _ := json.Marshal(requestToTransfer)
	err = stub.PutState(args[0], marbleJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) hospitalSuccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "key"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	requestAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if requestAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	requestToTransfer := marble{}
	err = json.Unmarshal(requestAsBytes, &requestToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	requestToTransfer.HospitalSuccess = true

	marbleJSONasBytes, _ := json.Marshal(requestToTransfer)
	err = stub.PutState(args[0], marbleJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}