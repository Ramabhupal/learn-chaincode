/*
Copyright IBM Corp 2016 All Rights Reserved.

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
	"errors"
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	


var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)

type userandlitres struct{
	User string        `json:"user"`
	Litres int       `json:"litres"`
}

type MilkContainer struct{

        ContainerID string `json:"containerid"`
	Userlist  [2]userandlitres    `json:"userlist"`

}

type Asset struct{
	User string        `json:"user"`
	ContainerIDs []string `json:"containerIDs"`
	LitresofMilk int `json:"litresofmilk"`
	Supplycoins int `json:"supplycoins"`
}

type Order struct{
       OrderID string                  `json:"orderid"`
       User string                     `json:"user"`
       Status string                   `json:"status"`
       Litres int                      `json:"litres"`
}

type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	
	if err != nil {
		return nil, err
	}
	
	var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(containerIndexStr, jsonAsBytes)                     //Resetting - Making milk container list as empty 
       if err != nil {
		return nil, err
        }  
	
	
/* Resetting the customer and market order list  */
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	
	var emptyasset Asset
	
	emptyasset.User = "Supplier"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	stub.PutState("SupplierAssets",jsonAsBytes) 
	emptyasset.User = "Market"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("MarketAssets", jsonAsBytes) 
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "checkstockbysupplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.checkstockbysupplier(stub, args)	
        }
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}


func  Create_milkcontainer(stub shim.ChaincodeStubInterface, args [3]string) ([]byte, error) {
var err error

// "1x223" "supplier" "20" 
// args[0] args[1] args[2] 
	
	if len(args) != 3{
		return nil, errors.New("Please enter all the details")
        }
	fmt.Println("Hold on, we are Creating milkcontainer asset for you")
	
id := args[0]
user := args[1]
litres,err:=strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Litres argument must be a numeric string")
	}
	
// Checking if the container already exists in the network
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details of given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

//If not present, create it and Update ledger, containerIndexStr, Assets of Supplier
//Creation
        res.ContainerID = id
	res.Userlist[0].User=user
	res.Userlist[0].Litres = litres
	milkAsBytes, _ =json.Marshal(res)
        stub.PutState(res.ContainerID,milkAsBytes)
	fmt.Printf("Container created successfully, details are %+v\n", res)
	
	
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	
	supplierasset.ContainerIDs = append(supplierasset.ContainerIDs, res.ContainerID)
	supplierasset.LitresofMilk += res.Userlist[0].Litres
	fmt.Println("Balance of Supplier")
        fmt.Printf("%+v\n", supplierasset)
	
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	
       supplierassetAsBytes,_ = stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	fmt.Printf("%+v\n", supplierasset)
	cid := supplierasset.ContainerIDs[0]
	userAssets := "MarketAssets"
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
	 
	fmt.Println("Updating ",userAssets)
	asset.LitresofMilk += 50
	fmt.Println("appending",cid,"to Market container id list")
        asset.ContainerIDs = append(asset.ContainerIDs,cid)
	fmt.Printf("%+v\n", asset)
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	assetAsBytes,_ = stub.GetState(userAssets)        // The same key which we used in Init function 
	json.Unmarshal( assetAsBytes, &asset)
	 fmt.Printf("%+v\n", asset)
	
	
/*
	for i := 0 ;i < len(supplierasset.containerIDs);i++{
	
            if(supplierasset.containerIDs[i] == "1x223"){

            supplierasset.containerIDs =     append(supplierasset.containerIDs[:i],supplierasset.containerIDs[i+1:]...)
           break
       }	
}
fmt.Printf("%+v\n", supplierasset)
*/	return nil,nil

}


func(t *SimpleChaincode)  checkstockbysupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//Market OrderID
//args[0]
	
Shiporder := Order{}
Shiporder.User = "Market"
Shiporder.Status = "Order placed to Supplier "
Shiporder.OrderID = args[1]
Shiporder.Litres = 5 * 10
//fetching assets of supplier	
	supplierassetAsBytes, _ := stub.GetState("SupplierAssets")
	supplierasset := Asset{}             
	json.Unmarshal(supplierassetAsBytes, &supplierasset )
	fmt.Printf("%+v\n", supplierasset)
//checking if Supplier has the stock	
if (supplierasset.LitresofMilk >= 50 ){
		fmt.Println("Enough stock is available, finding a suitable container.....")
		
	 fmt.Printf("%+v\n", supplierasset)
	cid := supplierasset.ContainerIDs[0]
	containerassetAsBytes, _ := stub.GetState(cid)
	res := MilkContainer{} 
	json.Unmarshal(containerassetAsBytes,&res)
      
	fmt.Println("Found a suitable container, below is the ID of the container, use it while placing order to Logistics")
	fmt.Printf("%+v\n", res)
	   // return nil, errors.New("Supplier has the quantity but not all in one container, this will be covered in next phase")
}else{
	        fmt.Println("Right now there isn't sufficient quantity , Create a new container")
		var b [3]string
		b[0] = "1x223"
		b[1] = "Supplier"
		b[2] = strconv.Itoa(Shiporder.Litres)
		Create_milkcontainer(stub,b)
		
	        fmt.Println("Successfully created container, check stock again to know your container details ") 
	        // can't call function again..loop hole
		return nil,nil
}
	return nil,nil
}


























// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
