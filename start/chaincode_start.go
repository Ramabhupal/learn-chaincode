

/******** SUPPLY CHAIN MANAGEMENT*****/



package main

import (
	"errors"
	"fmt"

	"encoding/json"



	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	

var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)
type MilkContainer struct{

        ContainerID string `json:"containerid"`
        User string        `json:"user"`

        Litres string        `json:"litres"`

}

type SupplyCoin struct{

        CoinID string `json:"coinid"`
        User string        `json:"user"`
}

type Order struct{
        OrderID string `json:"orderid"`
       User string `json:"user"`
       Status string `json:"status"`
       Litres string    `json:"litres"`
}

type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}

type Asset struct{
	  User string        `json:"user"`
	conatinerIDs []string `json:"containerids"`
	coinIds []string `json:"coinids"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func(t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
                                                          
}

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network 
       if err != nil {
		return nil, err
}
	
/* Making sure the value corresponding to containerIndexStr  is empty */

      var empty []string
	jsonAsBytes, _ := json.Marshal(empty)                   //create an empty string
	err = stub.PutState(containerIndexStr, jsonAsBytes)                 //Resetting - Making milk container list as empty 
	if err != nil {
		return nil, err
}  
	/*
        err = stub.PutState(coinIndexStr, jsonAsBytes)                 //Making coin list as empty
        if err != nil {
                return nil, err
}
*/  
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	// Resetting the Assets of Supplier for test case- later on we can do for all of them
	var empty Asset
	jsonAsBytes, _ = json.Marshal(empty)
	err = stub.PutState("SupplierAssets",jsonAsbytes)
	
	
        return nil, nil

}



func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
        }else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		return t.Create_milkcontainer(stub, args)      
        }else if function == "Create_coin" {		//creates a coin - invoked by market
		return t.Create_coin(stub, args)	
        } else if function == "Order_milk"{
		return t.Order_milk(stub,args)
	} else if function == "View_order{
	        return t.View_order(stub,args)
        }else if function == "init_logistics"{
	        return t.init_logistics(stub,args)
        }else if function == "set_user"{
	        return t.set_user(stub,args)
        }else if function == "checktheproduct"{
	       return t.checktheproduct(stub,args)
        }else if function == "cointransfer"{
	       return t.cointransfer(stub,args)
        }

       fmt.Println("invoke didn't find the function")
       return nil,nil

}


func (t *SimpleChaincode) Create_milkcontainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error

// "1x22" "supplier" 20 
// args[0] args[1] args[2] 

id := args[0]
user := args[1]
litres :=args[2] 
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details og given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

res.ContainerID = id
res.User = user
res.Litres = litres
milkAsBytes, _ =json.Marshal(res)

stub.PutState(res.ContainerID,milkAsBytes)
	
	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string                           //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(containerAsBytes, &containerIndex)	
	
	//append the newly created container to the global container list
	containerIndex = append(containerIndex, res.ContainerID)									//add marble name to index list
	fmt.Println("! container index: ", containerIndex)
	jsonAsBytes, _ := json.Marshal(containerIndex)
        err = stub.PutState(containerIndexStr, jsonAsBytes)

	 // append the container ID to the existing assets of the Supplier
	
	supplierassetAsBytes := stub.GetState("SupplierAssets")
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	supplierasset.containerIDs = append(supplierasset.containerIDs, res.ContainerID)
	supplierassetAsBytes = json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)

	
	
	
return nil,nil

}


func (t *SimpleChaincode) Create_coin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"1x245" "Market/Logistics"
id := args[0]
user:= args[1]

coinAsBytes , err := stub.GetState(id)
if err != nil{
              return nil, errors.New("Failed to get details of given id")
} 

res :=SupplyCoin{}

json.Unmarshal(coinAsBytes, &res)

if res.CoinID == id{

          fmt.Println("Coin already exists")
          fmt.Println(res)
          return nil,errors.New("This coin already exists")
}

res.CoinID = id
res.User = user

coinAsBytes, _ = json.Marshal(res)
stub.PutState(id,coinAsBytes)
return nil,nil
}
 

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

if function == "read" {						//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the variable to query")
	}


	name = args[0]
	valAsbytes, err := stub.GetState(name)				//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for  \"}"
		return nil, errors.New(jsonResp)
	}

return valAsbytes, nil										       //send it onward
}

func (t *SimpleChaincode) Order_milk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "pending"
Openorder.OrderID = "abcd"
Openorder.Litres = args[0]
orderasbytes,_ := json.Marshal(Openorder)
	var err error
err = stub.PutState(Openorder.OrderID,orderasbytes)
	
if err != nil {
		return nil, err
}

//Add the new order to the orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println("! appended open to orders")
	jsonAsBytes, _ = json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}
	
return nil,nil
}

func (t *SimpleChaincode) View_order(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	// This will be invoked by Supplier- think of UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty
/* fetching the Orders*/
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	
	
/*fetching the containers*/	
	
	containerasBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string             //an array to clone container indices
	json.Unmarshal(containerAsBytes, &containerIndex)
	
	containerAsBytes := stub.GetState(containerIndex[0])
	
	res := MilkContainer{} 
json.Unmarshal(containerAsBytes, &res)
	
	if res.Litres == orders.OpenOrders[0].Litres {
		fmt.Println("Found a suitable container")
		stub.PutState("hi",[]byte("Your product will be shipped soon"))
		orders.OpenOrders[0].Status = "Ready to be Shipped"
		//t.init_logistics(stub,orders.OpenOrders[0].OrderId, containerIndex[0])
	}
	
return nil,nil	
}

func (t *SimpleChaincode) init_logistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	fmt.Println("Inside Init logistics function")
	OrderID = args[0]
	ContainerID = args[1]
	
	
	ordersAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(ordersAsBytes, &ShipOrder)
	
	ShipOrder.Status = "In transit"
	 
	ordersAsBytes = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,ordersAsBytes)
	//t.set_user(stub,OrderID,ContainerID)
	
	
}


func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
// OrderId  ContainerID
	
//So here we will set the user name in conatiner ID to the one in Order ID and Status to Delivered
	
	OrderID = args[0]
	ContainerID = args[1]

ordersAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(ordersAsBytes, &ShipOrder)
	
	assetAsBytes := stub.GetState(ContainerID)
	container := MilkContainer{}
	json.Unmarshal(assetAsBytes, &container)
	
	container.User = ShipOrder.User
	
	//Pushing the updated container  back to the ledger
	assetAsBytes = json.marshal(container.User)
	stub.PutState(ContainerID, assetAsBytes)
	ShipOrder.Status = "Delivered"
	
	//pushing the updated Order back to ledger
	
	return nil,nil
	//t.checktheproduct(stub,OrderID,ContainerID)
}


func (t *SimpleChaincode) checktheproduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	
	OrderID = args[0]
	ContainerID = args[1]

	ordersAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(ordersAsBytes, &ShipOrder)

       assetAsBytes := stub.GetState(ContainerID)
	Deliveredcontainer := MilkContainer{}
	json.Unmarshal(assetAsBytes, &Deliveredcontainer)
	
	if Deliveredcontainer.User == "Market" && DeliveredContainer.Litres == ShipOrder.Litres {
		
		fmt.Println("Thanks, I got the product")
		stub.PutState("Market Response",[]byte("Product received"))
		//t.cointransfer(stub,coinid) coinid -hard code it and send the coin id created by market
		return nil,nil
       }

	return nil,nil


}



func (t *SimpleChaincode) cointransfer( stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	
	//lets keep it simple for now, just fetch the coin from ledger, change username to Supplier and End of Story
	CoinID = args[0]
	
	assetAsBytes := stub.GetState(CoinID)
	Transfercoin := SupplyCoin{}
	json.Unmarshal(assetAsBytes, &Transfercoin)
	
	if (Transfercoin.User == "Market")    // check if the market guy actually holds coin in his name
	{
		Transfercoin.User == "Supplier"
		assetAsBytes = json.Unmarhsal(Transfercoin)
		stub.PutState(CoinID, assetAsBytes)
		return nil,nil
		
	}else
	{
		fmt.Println("There was some issue in transferring")
		return nil,nil
	}

	
return nil,nil
	
	
/* END OF STORY*/

}
/*
func (t *SimpleChaincode) init_supplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
orderasbytes, _ := stub.GetState(args[0])
Performorder := Order{}
json.Unmarshal(orderasbytes, &Performorder)
containerasbytes, _ := stub.GetState("1x23")
Container := MilkContainer{}
json.Unmarshal(containerasbytes, &Container)
if(Container.Litres == Performorder.Litres){
fmt.Println("Hurray, we got want u want")
Performorder.Status="received"
orderasbytes,_=json.Marshal(Performorder)
stub.PutState(args[0],orderasbytes)
var a []string
a[0] = args[0]
a[1] ="1x23" 
t.init_logistics(stub,a)
return nil,nil
} else{
fmt.Println("Sorry")
return nil,nil
}
}
func (t *SimpleChaincode) init_logistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
orderid := args[0]
fmt.Println("Recevied order for shipping is %s ",orderid)
orderasbytes, _ := stub.GetState(orderid)
Shiporder := Order{}
json.Unmarshal(orderasbytes, &Shiporder)
Shiporder.Status = "Shipped and in transit"
orderasbytes,_ = json.Marshal(Shiporder)
stub.PutState(orderid, orderasbytes)
var a []string
a[0] = args[0]
a[1] ="1x23"
t.completedelivery(stub,a)
return nil,nil
}
func (t *SimpleChaincode) completedelivery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
orderid := args[0]
orderasbytes, _ := stub.GetState(orderid)
Shiporder := Order{}
json.Unmarshal(orderasbytes, &Shiporder)
milkasbytes, _ := stub.GetState(args[1])
milkcont := MilkContainer{}
json.Unmarshal(milkasbytes, &milkcont)
milkcont.User="Market"
Shiporder.Status = "Delivered"
orderasbytes,_ = json.Marshal(Shiporder)
stub.PutState(orderid, orderasbytes)
milkasbytes,_ = json.Marshal(milkcont)
stub.PutState(args[1],milkasbytes)
var a []string
a[0] = args[0]
a[1] =args[1]
t.checkproduct(stub,a)
return nil,nil
}
func (t *SimpleChaincode) checkproduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
milkcontainerid := args[1]
milkasbytes, _ := stub.GetState(milkcontainerid)
milkcontainer := MilkContainer{}
json.Unmarshal(milkasbytes,&milkcontainer)
if(milkcontainer.User == "Market"){
var a []string
a[0] = "1x245"
       t.init_cointransfer(stub,a)
       return nil,nil
}else{
      return nil,errors.New("Couldn't transfer,please try again")
}
return nil,nil
}
func (t *SimpleChaincode) init_cointransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
coinid := args[0]
coinasbytes, _ := stub.GetState(coinid)
Finalcoin:= SupplyCoin{}
json.Unmarshal(coinasbytes,&Finalcoin)
Finalcoin.User = "Supplier"
coinasbytes,_ = json.Marshal(Finalcoin)
stub.PutState(coinid, coinasbytes)
return nil,nil
}
*/

