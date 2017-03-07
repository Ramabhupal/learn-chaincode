
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

var customerOrdersStr = "_customerorders"    // This will  be the key, value will be a list of orders placed by customer - wil be called by Customer

var supplierOrdersStr = "_supplierorders"     // this will be key, value will be a list of orders placed by supplier to logistics

type userandlitres struct{
	User string        `json:"user"`
	Litres int       `json:"litres"`
}

type MilkContainer struct{

        ContainerID string `json:"containerid"`
	Userlist  [2]userandlitres    `json:"userlist"`

}

type Order struct{
       OrderID string                  `json:"orderid"`
       User string                     `json:"user"`
       Status string                   `json:"status"`
       Litres int                      `json:"litres"`
}


type SupplierOrder struct {
   
        OrderID string                `json:"orderid"`
	Towhom string                 `json:"towhom"`
	ContainerID string            `json:"containerid"`
	
}


type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}


type AllSupplierOrders struct {
        SupplierOrdersList []SupplierOrder  `supplierOrdersList`
}
	

type Asset struct{
	User string        `json:"user"`
	ContainerIDs []string `json:"containerIDs"`
	LitresofMilk int `json:"litresofmilk"`
	Supplycoins int `json:"supplycoins"`
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	fmt.Printf("every time we enter main function")
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	
	var err error
	
	fmt.Println("Welcome tothe Supply chain management Phase 1, Deployment has been started, do as u want")
	fmt.Printf("Hope this entire flow will go nicely")
 
       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }
	
/* Resetting the container list - Making sure the value corresponding to openOrdersStr is empty */
	
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
	err = stub.PutState(customerOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	
/* Resetting the supplier order list  */
	var suporders AllSupplierOrders
	suporderAsBytes,_ := json.Marshal(suporders)
	err = stub.PutState(supplierOrdersStr, suporderAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
// Resetting the Assets of Supplier,Market, Logistics, Customer
	
	var emptyasset Asset
	
	emptyasset.User = "Supplier"
	jsonAsBytes, _ = json.Marshal(emptyasset)                // this is the byte format format of empty Asset structure
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // key -Supplier assets and value is empty now --> Supplier has no assets
	emptyasset.User = "Market"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("MarketAssets", jsonAsBytes)         // key -Market assets and value is empty now --> Market has no assets
	emptyasset.User = "Logistics"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("LogisticsAssets", jsonAsBytes)      // key - Logistics assets and value is empty now --> Logistic has no assets
	emptyasset.User = "Customer"
	jsonAsBytes, _ = json.Marshal(emptyasset) 
	err = stub.PutState("CustomerAssets", jsonAsBytes)      // key - Customer assets and value is empty now --> Customer has no assets
	
	if err != nil {       
		return nil, err
}
	fmt.Println("Successfully deployed the code and orders and assets are reset")
	
return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	}else if function == "Create_coins" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coins(stub, args)	
        }else if function == "BuyMilkfrom_Retailer" { //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.BuyMilkfrom_Retailer(stub, args)	
        }else if function == "Vieworderby_Market" {  //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Vieworderby_Market(stub, args)	
        }else if function == "Checkstockby_Market" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Checkstockby_Market(stub, args)	
        }
	
	fmt.Println("invoke did not find func: " + function)

return nil, errors.New("Received unknown function invocation: " + function)
}


func (t *SimpleChaincode) Create_coins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"Market/Logistics/Customer",                  "100"
//args[0]                                     args[1]
//targeted owner                         No of supplycoins     
var err error
	user:= args[0]
	userAssets := user +"Assets"
        assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)

	asset.Supplycoins,err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New(" No of coins must be a numeric string")
	}
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	fmt.Println("Balance of " , user)
        fmt.Printf("%+v\n", asset)


return nil,nil
}



func (t *SimpleChaincode) BuyMilkfrom_Retailer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//args[0]      args[1]
//"cus123"       "10"
	var err error
	fmt.Println("Hello customer, welcome ")

	
	Openorder := Order{}
        Openorder.User = "customer"
        Openorder.Status = "Order received by Market"
        Openorder.OrderID = args[0]
        Openorder.Litres, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New(" No of coins must be a numeric string")
	}
	fmt.Println("Hello customer, your order has been generated successfully, you can track it with id in the following details")
	fmt.Println("%+v\n",Openorder)
        orderAsBytes,_ := json.Marshal(Openorder)
	stub.PutState(Openorder.OrderID,orderAsBytes)
	
	customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended",  Openorder.OrderID,"to existing customer orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(customerOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}

	return nil,nil
}

func(t *SimpleChaincode)  Vieworderby_Market(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET- think of UI-View orders- does he pass any parameter there...
// so here also no need of any arguments.
	
	fmt.Printf("Hello Market, these are the orders placed to  you by customer")
	
	
	ordersAsBytes, _ := stub.GetState(customerOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}


func (t *SimpleChaincode)  Checkstockby_Market(stub shim.ChaincodeStubInterface, args[]string) ([]byte, error){
	// In UI, beside each order one button to ship to customer, one button to check stock
	// we will extract details of orderId
	//we will exract asset balance of Market
	// if enough balance is der to deliver display "yes", if not der "no"
	//no tirggering is needed
	//OrderID should be passed in UI
//fetching order details
	OrderID := args[0]
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	quantity := ShipOrder.Litres
//fetching assets of market	
	marketassetAsBytes, _ := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset )
	
//checking if market has the stock	
	if (Marketasset.LitresofMilk >= quantity ){
		fmt.Println("Enough stock is available, Go ahead and deliver for customer")
		
//Call Deliver to customer function here
		b,_:= Deliverto_Customer(stub,ShipOrder.OrderID)
		fmt.Println(string(b))
		str := "Delivered to customer"
		return []byte(str), nil
		
	}else{
	        fmt.Println("Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer")
		str :=  "Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer"
	        ShipOrder.Status = "In transit to customer" // No matter, where the order placed by market is , for customer we will show it is "in transit"
	        orderAsBytes,err = json.Marshal(ShipOrder)
                stub.PutState(OrderID,orderAsBytes)  
		
		customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)	
	
	
		for i :=0; i<len(orders.OpenOrders);i++{
			if (orders.OpenOrders[i].OrderID == ShipOrder.OrderID){
			orders.OpenOrders[i].Status = "In transit to customer"
		         customerordersAsBytes , _ = json.Marshal(orders)
                        stub.PutState(customerOrdersStr,  customerordersAsBytes)
			}
	       }
	  return []byte(str), nil
		
		//Now we should send details of updated order status to customer, should be done in UI

		
        }
	
	return nil,nil

}

func Deliverto_Customer(stub shim.ChaincodeStubInterface ,args string) ([]byte,error){

	//args[0] 
	//OrderID  
	
	fmt.Println("Inside deliver to customer function")
//customer order
	OrderID := args
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	fmt.Println("%+v\n", ShipOrder)
	quantity := ShipOrder.Litres
	fmt.Println(quantity)
//market and customer assets
        marketassetAsBytes, _ := stub.GetState("MarketAssets")
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset)
	fmt.Printf("%+v\n", Marketasset) 
	customerassetAsBytes, _ := stub.GetState("CustomerAssets")
	Customerasset := Asset{}             
	json.Unmarshal(customerassetAsBytes, &Customerasset)
	fmt.Printf("%+v\n", Customerasset) 
if (Marketasset.LitresofMilk >= quantity ){
	fmt.Println("Inside deliver to customer, market has quantity")
	
	id := Marketasset.ContainerIDs[0]
	
	
	milkAsBytes, err := stub.GetState(id) 
        if err != nil {
		return nil, errors.New("Failed to get details of given id") 
        }

        res := MilkContainer{} 
        json.Unmarshal(milkAsBytes, &res)
		
	fmt.Printf("%+v\n", res)
	
	
	
	
	
   // here we are assuming only one container is der and it has enough stock to provide
	if ( res.Userlist[0].Litres - quantity >0) {
		fmt.Println("yo yo..its about to complete")
                    
   //updating the container details, bcz it is shared now
		res.Userlist[0].Litres -= quantity // bringing down the market share of it
		res.Userlist[1].User = "Customer"
		res.Userlist[1].Litres = quantity
		fmt.Printf("%+v\n", res)
		milkAsBytes, _ =json.Marshal(res)
                stub.PutState(res.ContainerID,milkAsBytes)
		
  //updating customer assets
		
	              Customerasset.LitresofMilk += quantity
		if ( len(Customerasset.ContainerIDs) == 0){
		      fmt.Println("This is the first container of customer")
	   Customerasset.ContainerIDs = append(Customerasset.ContainerIDs ,id)
		}
		fmt.Printf("%+v\n", Customerasset)
			    Marketasset.LitresofMilk -= quantity
	
	              customerassetAsBytes,_ = json.Marshal(Customerasset)
	              stub.PutState("CustomerAssets",customerassetAsBytes)
	
	               marketassetAsBytes,_ = json.Marshal(Marketasset)
	               stub.PutState("MarketAssets",marketassetAsBytes)
	
	               ShipOrder.Status ="Delivered to Customer"
	               fmt.Printf("%+v\n", ShipOrder)
	               orderAsBytes,err = json.Marshal(ShipOrder)
                       stub.PutState(OrderID,orderAsBytes)
	
        customerordersAsBytes, err := stub.GetState(customerOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(customerordersAsBytes, &orders)				
	
		for i :=0; i<len(orders.OpenOrders);i++{
			if (orders.OpenOrders[i].OrderID == ShipOrder.OrderID){
			orders.OpenOrders[i].Status = "Delivered to customer"
		         customerordersAsBytes , _ = json.Marshal(orders)
                        stub.PutState(customerOrdersStr,  customerordersAsBytes)
			}
	       }
	  
		//b := [3]string{"30", "Customer", "Market"}
	           //    transfer(stub,b)        //Transfer should be automated. So it can't be invoked from UI..Loop hole
	               fmt.Println("FINALLLLLYYYY, END OF THE STORY")
         
                      return nil,nil
	}else{
	       return nil, errors.New("On a whole market has quantity, but it is divided into container, right now we are not going to that level")
	}
}else{
         return nil, errors.New(" No stock, give order to supplier")
 }

}







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
