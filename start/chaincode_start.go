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

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"  //bytes包提供了对字节切片进行读写操作的一系列函数 
	"encoding/json"
	"fmt" //打印等函数的包
	"strconv"//字符转换包

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Coupon structure, with 4 properties.  Structure tags are used by encoding/json library
type Coupon struct {               //优惠劵结构类型
	Number   string `json:"number"` 	//数字字符串
	Amount string `json:"amount"`
	Flag string `json:"flag"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCoupon" {
		return s.queryCoupon(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCoupon" {
		return s.createCoupon(APIstub, args)
	} else if function == "consumeCoupon" { //消耗代金券
		return s.queryAllCoupon(APIstub)
	} else if function == "changeCouponOwner" {
		return s.changeCouponOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCoupon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	couponAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(couponAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	coupons := []Coupon{
		Coupon{Number: "001", Amount: "100", Flag: "0", Owner: "hanshuang"},
	}

	i := 0
	for i < len(coupons) {
		fmt.Println("i is ", i)
		couponAsBytes, _ := json.Marshal(coupons[i])
		APIstub.PutState("COUPON"+strconv.Itoa(i), couponAsBytes)
		fmt.Println("Added", coupons[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCoupon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var coupon = Coupon{Number: args[1], Amount: args[2], Flag: args[3], Owner: args[4]}

	couponAsBytes, _ := json.Marshal(coupon)
	APIstub.PutState(args[0], couponAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) consumeCoupon(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	couponAsBytes, _ := APIstub.GetState(args[0])
	coupon := Coupon{}

	json.Unmarshal(couponAsBytes, &coupon)
	coupon.Flag = args[1]

	carAsBytes, _ = json.Marshal(coupon)
	APIstub.PutState(args[0], couponAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeCouponOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	couponAsBytes, _ := APIstub.GetState(args[0])
	coupon := Coupon{}

	json.Unmarshal(couponAsBytes, &coupon)
	coupon.Owner = args[1]

	carAsBytes, _ = json.Marshal(coupon)
	APIstub.PutState(args[0], couponAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
