
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

//var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs
var batchIndexStr = "_batchindex"


var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)

var customerOrdersStr = "_customerorders"    // This will  be the key, value will be a list of orders placed by customer - wil be called by Customer

var supplierOrdersStr = "_supplierorders"     // this will be key, value will be a list of orders placed by supplier to logistics


var Count int                          //To keep count of Boxes created
//In our case product is a discrete one like Pendrive, hard disk etc
type Product struct{
       ProductID string               `json:"productid"`
	     Owner string                     `json:"owner"`
       Status string                   `json:"status"`

}



type Batch struct{
       BatchID string               `json:"batchid"`
       Quantity int                 `json:"quantity"`
       Productlist []string      `json:"productlist"`
       Owner string                 `json:"owner"`
       Status string                `json:"status"`

}

type Order struct{
       OrderID string                  `json:"orderid"`
       Owner string                     `json:"owner"`
       Status string                   `json:"status"`
       Quantity int                      `json:"quantity"`
}

// Quantity implies Batch if order placed by Retailer, product if placed by Customer
type SupplierOrder struct {

        OrderID string                `json:"orderid"`
	Towhom string                 `json:"towhom"`
	BatchID string            `json:"BatchID"`

}


type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}


type AllSupplierOrders struct {
        SupplierOrdersList []SupplierOrder  `supplierOrdersList`
}


type Asset struct{
	Owner string        `json:"owner"`
	//BatchIDs []string `json:"BatchIDs"`
	BatchIDs []string `json:"batchIDs"`
	NumberofProducts int `json:"numberofproducts"`
	Supplycoins int `json:"supplycoins"`
}



func main() {
	err := shim.Start(new(SimpleChaincode))

	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	fmt.Printf("every time we enter main function")
}

// Init Newbatchets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {


	var err error

	fmt.Println("Welcome  to  Supply chain management , Deployment has been started...")
	fmt.Printf("Hope for best, Plan for the worst")


       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }

/* Newbatchetting the container list - Making sure the value corNewbatchponding to openOrdersStr is empty */

       var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(batchIndexStr, jsonAsBytes)                     //Newbatchetting - Making milk container list as empty
       if err != nil {
		return nil, err
        }


/*making no of batches to zero*/
	Count = 0

/* Updating the customer and market order list  */
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

/* Newbatchetting the supplier order list  */
	var suporders AllSupplierOrders
	suporderAsBytes,_ := json.Marshal(suporders)
	err = stub.PutState(supplierOrdersStr, suporderAsBytes)                 //So the value for key is null
	if err != nil {
		return nil, err
}
// Newbatchetting the Assets of Supplier,Market, Logistics, Customer

	var emptyasset Asset

	emptyasset.Owner = "Supplier"
	jsonAsBytes, _ = json.Marshal(emptyasset)                // this is the byte format format of empty Asset structure
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // key -Supplier assets and value is empty now --> Supplier has no assets
	emptyasset.Owner = "Retailer"
	jsonAsBytes, _ = json.Marshal(emptyasset)
	err = stub.PutState("RetailerAssets", jsonAsBytes)         // key -Market assets and value is empty now --> Market has no assets
	emptyasset.Owner = "Logistics"
	jsonAsBytes, _ = json.Marshal(emptyasset)
	err = stub.PutState("LogisticsAssets", jsonAsBytes)      // key - Logistics assets and value is empty now --> Logistic has no assets
	emptyasset.Owner = "Customer"
	jsonAsBytes, _ = json.Marshal(emptyasset)
	err = stub.PutState("CustomerAssets", jsonAsBytes)      // key - Customer assets and value is empty now --> Customer has no assets

	if err != nil {
		return nil, err
}
	fmt.Println("Successfully deployed the code and orders and assets are Newbatchet")
	fmt.Printf("Go a head and play around")

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
        }else if function == "Buyproductfrom_Retailer" { //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Buyproductfrom_Retailer(stub, args)
        }else if function == "Vieworderby_Retailer" {  //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Vieworderby_Retailer(stub, args)
        }else if function == "Checkstockby_Retailer" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Checkstockby_Retailer(stub, args)
        }else if function == "Orderto_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Orderto_Supplier(stub, args)
        }else if function == "Vieworderby_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Vieworderby_Supplier(stub, args)
        }else if function == "Checkstockby_Supplier" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Checkstockby_Supplier(stub,args)
        }else if function == "Call_Logistics" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Call_Logistics(stub, args)
        }else if function == "Vieworderby_Logistics"{
                return t.Vieworderby_Logistics(stub,args)
        }else if function == "pickuptheproduct" {
                return t.pickuptheproduct(stub,args)
        }else if function == "Deliverto_Retailer" {
                return t.Deliverto_Retailer(stub,args)
        }
	fmt.Println("invoke did not find func: " + function)

return nil, errors.New("Received unknown function invocation: " + function)
}



func  Create_Batch(stub shim.ChaincodeStubInterface, args string ) ( error) {

//args
//No of batches in string format

//var err error
Quantityofbatches,_ := strconv.Atoi(args)             // No of batches to be created
Productsperbatch := 10
owner := "Supplier"
status := "Manufactured"


for j:=0;j<Quantityofbatches;j++{
//Each time this outer loop runs a new batch of products is created
        Count += 1
	batchid := "Batch"+strconv.Itoa(Count)           //This has to be linked to QR Code generated later on
	fmt.Println(batchid)

	batchAsBytes, err := stub.GetState(batchid)
        if err != nil {
		return  errors.New("Failed to get details of given id")
        }


	Newbatch := Batch{}
	json.Unmarshal(batchAsBytes, &Newbatch)

if  Newbatch.BatchID == batchid{

        fmt.Println("batch already exixts")
        fmt.Println("%+v\n",Newbatch)
        return errors.New("This batch already exists")
}

//If not der, creating
	Newbatch.BatchID = batchid
        Newbatch.Owner = owner
        Newbatch.Status = status
        Newbatch.Quantity = Productsperbatch

	for i:=0; i < Newbatch.Quantity ;i++{
	 	 Newproduct := Product{}
	   Newproduct.ProductID = "prod"+strconv.Itoa(Count)+"."+strconv.Itoa(i)
     Newproduct.Owner = owner
     Newproduct.Status = status
     productasbytes ,_ := json.Marshal( Newproduct)
		 stub.PutState( Newproduct.ProductID,productasbytes)

     Newbatch.Productlist = append(Newbatch.Productlist, Newproduct.ProductID)
	}

fmt.Printf("%+v\n", Newbatch)
	batchAsBytes,_ = json.Marshal(Newbatch)
	stub.PutState(Newbatch.BatchID,batchAsBytes)

	//Update batchIndexStr
        batchindexAsBytes, err := stub.GetState(batchIndexStr)
	if err != nil {
		return  errors.New("Failed to get container index")
	}
	var batchIndex []string                                        //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(batchindexAsBytes, &batchIndex)


	batchIndex = append(batchIndex, Newbatch.BatchID)          //append the newly created container to the global container list									//add marble name to index list
	fmt.Println("batch indices in the network: ", batchIndex)
	batchindexAsBytes, _ = json.Marshal(batchIndex)
        err = stub.PutState(batchIndexStr, batchindexAsBytes)


// append the Batch ID to the existing assets of the Supplier

	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)

	supplierasset.BatchIDs = append(supplierasset.BatchIDs, Newbatch.BatchID)
	supplierasset.NumberofProducts += Newbatch.Quantity
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	fmt.Println("Balance of Supplier")
        fmt.Printf("%+v\n", supplierasset)
    //double checking
	supplierassetAsBytes,_ = stub.GetState("SupplierAssets")        // The same key which we used in Init function

	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	fmt.Printf("%+v\n", supplierasset)




}






return nil

}


func (t *SimpleChaincode) Create_coins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"Retailer/Logistics/Customer",                  "100"
//args[0]                                     args[1]
//targeted owner                         No of supplycoins
var err error
	owner:= args[0]
	OwnerAssets := owner +"Assets"
        assetAsBytes,_ := stub.GetState(OwnerAssets)        // The same key which we used in Init function
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)

	asset.Supplycoins,err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New(" No of coins must be a numeric string")
	}
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(OwnerAssets,assetAsBytes)
	fmt.Println("Balance of " , owner)
        fmt.Printf("%+v\n", asset)


return nil,nil
}



func (t *SimpleChaincode) Buyproductfrom_Retailer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//args[0]      args[1]
//"cus123"       "10"
	var err error
	fmt.Println("Hello customer, welcome ")


	Openorder := Order{}
        Openorder.Owner = "customer"
        Openorder.Status = "Order received by Retailer"
        Openorder.OrderID = args[0]
        Openorder.Quantity, err = strconv.Atoi(args[1])  ///No of units of product the customer wants
	if err != nil {
		return nil, errors.New(" No of products must be a numeric string")
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

func(t *SimpleChaincode)  Vieworderby_Retailer(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET
	fmt.Printf("Hello Retailer, these are the orders placed to  you by customer")


	ordersAsBytes, _ := stub.GetState(customerOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}



func (t *SimpleChaincode)  Checkstockby_Retailer(stub shim.ChaincodeStubInterface, args[]string) ([]byte, error){
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
	quantity := ShipOrder.Quantity        //Here quantity implies no of unit products ..not batches
//fetching assets of retailer
	retailerassetAsBytes, _ := stub.GetState("RetailerAssets")
	Retailerasset := Asset{}
	json.Unmarshal(retailerassetAsBytes, &Retailerasset )

//checking if market has the stock
	if (Retailerasset.NumberofProducts >= quantity && ShipOrder.Status != "Delivered to customer"  ){
		fmt.Println("Enough Number of Products are available, Go ahead and deliver to customer")

//Call Deliver to customer function here
		b,_:= Deliverto_Customer(stub,ShipOrder.OrderID)
		fmt.Println(string(b))
		str := "Delivered to customer"
		return []byte(str), nil

	}else{
	        fmt.Println("Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer")
		      str :=  "Right now there isn't sufficient quantity , Give order to Supplier/Manufacturer"
	        ShipOrder.Status = "In transit" // No matter, where the order placed by market is , for customer we will show it is "in transit"
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
			orders.OpenOrders[i].Status = "In transit "
		         customerordersAsBytes , _ = json.Marshal(orders)
                        stub.PutState(customerOrdersStr,  customerordersAsBytes)
			}
	       }
	  return []byte(str), nil

		//Now we should send details of updated order status to customer, should be done in UI
  // A pop up should come to give order to Supplier

        }


//Should keep a flag kind of thing to see which one happend and catch that flag in UI and display accordingly
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
	quantity := ShipOrder.Quantity
	//market and customer assets
  retailerassetAsBytes, _ := stub.GetState("RetailerAssets")
	Retailerasset := Asset{}
	json.Unmarshal(retailerassetAsBytes, &Retailerasset)
	fmt.Printf("%+v\n", Retailerasset)
	customerassetAsBytes, _ := stub.GetState("CustomerAssets")
	Customerasset := Asset{}
	json.Unmarshal(customerassetAsBytes, &Customerasset)
	fmt.Printf("%+v\n", Customerasset)
if (Retailerasset.NumberofProducts >= quantity ){
	fmt.Println("Inside deliver to customer, market has quantity")

	id := Retailerasset.BatchIDs[0]


	batchAsBytes, err := stub.GetState(id)
        if err != nil {
		return nil, errors.New("Failed to get details of given id")
        }

        Newbatch := Batch{}
        json.Unmarshal(batchAsBytes, &Newbatch)

	fmt.Printf("%+v\n", Newbatch)





   // here we are assuming only one container is der and it has enough stock to provide
	if ( Newbatch.Quantity - quantity >0) {
		fmt.Println("yo yo..its about to complete")

        //updating the container details, bcz it is shared now
		Newbatch.Quantity -= quantity // bringing down the Retailer share of it
    Retailerasset.NumberofProducts -= quantity
    Customerasset.NumberofProducts += quantity

	 Customerasset.BatchIDs = append(Customerasset.BatchIDs,Newbatch.Productlist[0:quantity]...)  //adding prod id to customer i.e now the product is with customer
	 Newbatch.Productlist = Newbatch.Productlist[quantity:]     //Remving the products from the box
        for  k:=0;k<len(Customerasset.BatchIDs);k++{


					Newproduct := Product{}
productasbytes ,_ := stub.GetState(Customerasset.BatchIDs[k])
		json.Unmarshal(productasbytes,&Newproduct)
					Newproduct.Status = "Delivered to customer"
					Newproduct.Owner = "Customer"
					productasbytes ,_ = json.Marshal( Newproduct)
					stub.PutState( Newproduct.ProductID,productasbytes)



				}


//updating assets
  customerassetAsBytes,_ = json.Marshal(Customerasset)
	stub.PutState("CustomerAssets",customerassetAsBytes)

	 retailerassetAsBytes,_ = json.Marshal(Retailerasset)
	 stub.PutState("RetailerAssets",retailerassetAsBytes)

	 batchAsBytes,_ = json.Marshal(Newbatch)
	 stub.PutState(Newbatch.BatchID,batchAsBytes)
//updating orders
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

		b := [3]string{"30", "Customer", "Retailer"}
	           transfer(stub,b)        //Transfer should be automated. So it can't be invoked from UI..Loop hole
	               fmt.Println("FINALLLLLYYYY, END OF THE STORY")

                      return nil,nil
	}else{
	       return nil, errors.New("On a whole market has quantity, but it is divided into container, right now we are not going to that level")
	}
}else{
         return nil, errors.New(" No stock, give order to supplier")
 }

}


func(t *SimpleChaincode) Orderto_Supplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
// "cus123"           "abcd"
// CustomerOrderID    MarketOrderID

var err error

//fetching the customer order details and ordering 5 times to the Quantity customer asked
CustomerOrderID := args[0]
orderAsBytes, err := stub.GetState(CustomerOrderID)
	if err != nil {
		return  nil, errors.New("Failed to get details of customer order, please make sure your id is correct")
	}
CustomerOrder := Order{}
json.Unmarshal(orderAsBytes, &CustomerOrder)
quantity := CustomerOrder.Quantity

//Generating market order

Openorder := Order{}
Openorder.Owner = "Retailer"
Openorder.Status = "Order placed to Supplier "
Openorder.OrderID = args[1]
Openorder.Quantity = ((quantity - (quantity % 10)) /10) + 1 //Rounding to the nearest number of boxes that can suffice for the quantity asked by customer

orderAsBytes,_ = json.Marshal(Openorder)
stub.PutState(Openorder.OrderID,orderAsBytes)
fmt.Println("your Order has been generated successfully")
fmt.Printf("%+v\n", Openorder)

//Add the new market order to market orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get  existing list of Retailer orders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println(" appended ",Openorder.OrderID,"to existing market orders")
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
        }


return nil,nil
}



func(t *SimpleChaincode)  Vieworderby_Supplier(stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {
// This will be invoked by MARKET- think of UI-View orders- does he pass any parameter there...
// so here also no need of any arguments.

	fmt.Printf("Hello Supplier, these are the orders placed to  you by Retailer")


	ordersAsBytes, _ := stub.GetState(openOrdersStr)
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)
	//This should stop here.. In UI it should display all the orders - beside each order -one button "ship to customer"
	//If we click on any order, it should call query for that OrderID. So it will be enough if we update OrderID and push it to ledger
	fmt.Println(orders)
	 return nil,nil
}




func(t *SimpleChaincode)  Checkstockby_Supplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
/***FUNCTIONALITY EXPLAINED*******/
// In UI, beside each order one button to call logistics, one button to check stock
// we will extract details of orderId
//we will exract asset balance of Market
// if enough balance is der --> find a container and show it, if not create a new container (automated) and show it
//At the end of this function we will end up with a container
/*******/


//OrderID should be passed in UI
//fetching order details
//Market OrderID
//args[0]

	OrderID := args[0]
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	ShipOrder := Order{}
	json.Unmarshal(orderAsBytes, &ShipOrder)
	quantity := ShipOrder.Quantity               //Here quantity is no of boxes
//fetching assets of supplier
	supplierassetAsBytes, _ := stub.GetState("SupplierAssets")
	supplierasset := Asset{}
	json.Unmarshal(supplierassetAsBytes, &supplierasset )
	fmt.Printf("%+v\n", supplierasset)
//checking if Supplier has the stock
if (len(supplierasset.BatchIDs) >= quantity  ){
		fmt.Println("Enough number of batches are available, below are the details")
		fmt.Printf("%+v\n", supplierasset)
	  cid := supplierasset.BatchIDs[0]
	  batchassetAsBytes, _ := stub.GetState(cid)
	  Newbatch := Batch{}
	  json.Unmarshal(batchassetAsBytes,&Newbatch)

	  fmt.Printf("%+v\n", Newbatch)

  }else{
	        fmt.Println("Right now there isn't sufficient quantity , Create a new container")
		b := strconv.Itoa(quantity)
	 Create_Batch(stub,b)


	       // fmt.Println("Successfully created container, check stock again to know your container details ")
	        // can't call function again..loop hole
		//return nil,nil
}
	return nil,nil
}

func (t *SimpleChaincode)  Call_Logistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

//args[0]   //ToWhom  //Container ID
//OrderID   //Retailer   //"1x223"

// I think its fair only, in practical case, we will tell adrNewbatchs for a postman to deliver, same thing here also
//Here Postman is Logistics guy, Receiver is market, letter is Container

	ShipOrder := SupplierOrder{}
	ShipOrder.OrderID = args[0]
	ShipOrder.Towhom = args[1]
	ShipOrder.BatchID = args[2]

	orderAsBytes, _ :=json.Marshal(ShipOrder)
	stub.PutState( ShipOrder.OrderID, orderAsBytes)

	fmt.Println("Successfully placed order to Logistics")
	fmt.Println("%+v\n", ShipOrder)


	//Add the new Supplier order to market orders list
	ordersAsBytes, err := stub.GetState(supplierOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get  existing list of  orders placed by Supplier to logistics")
	}
	var suporders AllSupplierOrders
	json.Unmarshal(ordersAsBytes, &suporders)
	suporders.SupplierOrdersList  = append(suporders.SupplierOrdersList, ShipOrder);		//append the new order - Openorder
	fmt.Println(" appended ",ShipOrder.OrderID,"to existing orders placed by Supplier to logistics")
	jsonAsBytes, _ := json.Marshal(suporders)
	err = stub.PutState(supplierOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
        }


	return nil,nil

}



func(t *SimpleChaincode) Vieworderby_Logistics(stub shim.ChaincodeStubInterface, args []string) ( []byte , error) {

	// This will be invoked by Supplier in UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty-but just for syntax-pass something as parameter in angular js


//fetching the Orders
	fmt.Printf("Hello Logistics, here are the orders placed to you by Supplier")
	fmt.Printf("Go a head and do your business")


	ordersAsBytes, _ := stub.GetState(supplierOrdersStr)

	var orders AllSupplierOrders
	json.Unmarshal(ordersAsBytes, &orders)

	fmt.Println(orders)
	return nil,nil
}



func(t *SimpleChaincode) pickuptheproduct(stub shim.ChaincodeStubInterface, args []string) ( []byte , error) {

// So in view order, he will see his orders, clicking on the order will show to whom and which container
//There will be a button "pickuptheproduct" which is equivalent to real life pick up --status will be in transit
//There will be one more button there only "Delivertheproduct"
//As of march 3, lets pass market order Id only as argument
//How can we update the order placed by market...without a notification
//here we are passing market order id only
	//args[0] args[1]
	// MarketOrderID, BatchID
	RetailerOrderID := args[0]


	// fetch the order details and update status as "in transit"
	orderAsBytes, err := stub.GetState(RetailerOrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	RetailerOrder := Order{}
	json.Unmarshal(orderAsBytes, &RetailerOrder)

	RetailerOrder.Status = "In transit to Retailer"

	orderAsBytes,err = json.Marshal(RetailerOrder)

	stub.PutState(RetailerOrderID,orderAsBytes)

	fmt.Printf("%+v\n", RetailerOrder)

	fmt.Println("Batch  is in transit")

return nil,nil
}



func(t *SimpleChaincode)  Deliverto_Retailer(stub shim.ChaincodeStubInterface, args []string) ([]byte , error) {

// SupplierOrderID      //MarketOrderID
//args[0]               //args[1]

//So here we will set the Owner name in container ID to the one in Order ID and Status to Delivered - Asset Transfer
// Why should logistics guy check if the supplier actually holds the container?????????
	fmt.Println("Delivering the container to Retailer")
	OrderID := args[0]

//fetch order details
       orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  nil,errors.New("Failed to get openorders")
	}
	ShipOrder := SupplierOrder{}
	json.Unmarshal(orderAsBytes, &ShipOrder)

	BatchID := ShipOrder.BatchID
//fetch container details
	assetAsBytes,err := stub.GetState(BatchID)
	Newbatch := Batch{}
	json.Unmarshal(assetAsBytes, &Newbatch)

	if (Newbatch.Owner == "Supplier"){

	Newbatch.Owner = "Retailer"         //ASSET TRANSFER
	fmt.Println("%+v\n", Newbatch)
	fmt.Println("pushing the updated Batch back to ledger")
	assetAsBytes,err = json.Marshal(Newbatch)
	stub.PutState(BatchID, assetAsBytes)    //Pushing the updated container  back to the ledger

//fetch supplier assets
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
//fetch market assets
	OwnerAssets := "RetailerAssets"
	assetAsBytes,_ := stub.GetState(OwnerAssets)        // The same key which we used in Init function
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
//update market assets
	fmt.Println("Updating ",OwnerAssets)
	asset.NumberofProducts += Newbatch.Quantity
	fmt.Println("appending", BatchID,"to Retailer batch id list")
        asset.BatchIDs = append(asset.BatchIDs,BatchID)
       fmt.Printf("%+v\n", asset)
//update supplierassets

	fmt.Println("Updating Supplier assets..")
	supplierasset.NumberofProducts -= Newbatch.Quantity

	//WRITE A CODE  to remove that container id from supplier id list

		for i := 0 ;i < len(supplierasset.BatchIDs);i++{

            if(supplierasset.BatchIDs[i] == BatchID){

            supplierasset.BatchIDs =  append(supplierasset.BatchIDs[:i],supplierasset.BatchIDs[i+1:]...)
           break
       }
}
	fmt.Printf("%+v\n", supplierasset)

//pushing updated ledger back to ledger
        supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)

	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(OwnerAssets,assetAsBytes)

//double check
	supplierassetAsBytes,_ = stub.GetState("SupplierAssets")        // The same key which we used in Init function
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
        fmt.Printf("%+v\n", supplierasset)

	assetAsBytes,_ = stub.GetState(OwnerAssets)        // The same key which we used in Init function
	json.Unmarshal( assetAsBytes, &asset)
	 fmt.Printf("%+v\n", asset)
//update the RetailerOrder and push back to ledger

	RetailerOrderID := args[1]
        orderAsBytes, err = stub.GetState(RetailerOrderID)
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	RetailerOrder := Order{}
	json.Unmarshal(orderAsBytes, &RetailerOrder)
	RetailerOrder.Status = "Delivered to market"
	orderAsBytes,err = json.Marshal(RetailerOrder)
	stub.PutState(RetailerOrderID,orderAsBytes)
	fmt.Printf("%+v\n", ShipOrder)



	var b [2]string
	b[0] = args[1]
	b[1] = BatchID   //Later on this will be an array of batch IDs

	checktheproduct(stub,b)

	}else
        {
                stub.PutState("delivertomarket",[]byte("failure in this function"))
                //t.read(stub,"setOwner")
								fmt.Printf("Failure in Deliverto_Retailer")
                return nil,nil
        }


return nil,nil
}



func  checktheproduct(stub shim.ChaincodeStubInterface, args [2]string) ( error) {

// args[0] args[1]
// MarketOrderID, BatchID
	fmt.Println("Let us check the product")
	OrderID := args[0]
	BatchID := args[1]
//fetch order details
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{}
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details
       assetAsBytes,_ := stub.GetState(BatchID)
	Deliveredbatch := Batch{}
	json.Unmarshal(assetAsBytes, &Deliveredbatch)

//check and transfer coins
	if (Deliveredbatch.Owner == "Retailer" && Deliveredbatch.Quantity == ShipOrder.Quantity * 10) {

		fmt.Println("Thanks, I got  the right product, transferring amount to Supplier/Manufacturer")
		var b [3]string
		b[0]= strconv.Itoa(ShipOrder.Quantity *30)   //price per batch 30
		b[1] = "Retailer"
		b[2] = "Supplier"

		err = transfer(stub,b)
		if err!=nil{
			return err
		}
	  b[0]= strconv.Itoa(25)
		b[1] = "Supplier"
		b[2] = "Logistics"
		err = transfer(stub,b)
		if err!=nil{
			return err
		}
		return nil
}else{
    stub.PutState("checktheproduct",[]byte("failure"))
		fmt.Println("I didn't get the right product")
    return nil
 }


return nil


}



func transfer( stub shim.ChaincodeStubInterface, args [3]string) ( error) {

//args[0]             args[1]         args[2]
//No of supplycoin      Sender         Reciever
	//lets keep it simple for now, just fetch the coin from ledger, change Ownername to Supplier and End of Story
	transferamount,_ := strconv.Atoi(args[0])
	sender := args[1]                               // this thing should be given by us in UI background
	receiver := args[2]                            // this will be given by the Owner on web page

	fmt.Println( sender, "transferring", transferamount, "coins to", receiver)

        senderAssets := sender +"Assets"
        senderassetAsBytes,_ := stub.GetState(senderAssets)        // The same key which we used in Init function
	senderasset := Asset{}
	json.Unmarshal( senderassetAsBytes, &senderasset)


	receiverAssets := receiver+"Assets"
        receiverassetAsBytes,_ := stub.GetState(receiverAssets)        // The same key which we used in Init function
	receiverasset := Asset{}
	json.Unmarshal( receiverassetAsBytes, &receiverasset)

	if ( senderasset.Supplycoins >= transferamount){

	senderasset.Supplycoins -= transferamount
	receiverasset.Supplycoins += transferamount

        senderassetAsBytes,_=  json.Marshal(senderasset)
	stub.PutState(senderAssets,  senderassetAsBytes)
	fmt.Println("Balance of " , sender)
       fmt.Printf("%+v\n", senderasset)

	receiverassetAsBytes,_=  json.Marshal(receiverasset)
	stub.PutState( receiverAssets,receiverassetAsBytes)
	fmt.Println("Balance of " , receiver)
        fmt.Printf("%+v\n", receiverasset)
		return  nil
	}else {
		str := "Failed to transfer amount from" + sender + "to" + receiver
		return  errors.New(str)
	}


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
	var key, jsonNewbatchp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonNewbatchp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonNewbatchp)
	}

	return valAsbytes, nil
}
