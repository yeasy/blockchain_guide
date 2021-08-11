// Usage: ./prog inputPkg outputPkg

package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/hyperledger/fabric/protos/peer"
	"io/ioutil"
	"os"
)

func check(msg string, e error) {
	if e != nil {
		fmt.Println(msg)
		panic(e)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s inputPkg outputPkg\n", os.Args[0])
		return
	}

	inputPkg := os.Args[1]
	outputPkg := os.Args[2]
	fmt.Printf("Input pkg=%s, output=%s\n", inputPkg, outputPkg)

	ccbytes, err := ioutil.ReadFile(inputPkg)
	check("Read file failed", err)

	buf := ccbytes
	depSpec := &pb.ChaincodeDeploymentSpec{}
	err = proto.Unmarshal(buf, depSpec)
	check("Unmarshal depSpec failed", err)

	fmt.Printf("chaincodeSpec=%+v\n", depSpec.ChaincodeSpec)

	payload := depSpec.CodePackage
	err = ioutil.WriteFile(outputPkg, payload, 0644)
	check("Write to file failed", err)
}

