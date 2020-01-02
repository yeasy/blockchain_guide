package main

import (
	"fmt"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	pb "github.com/hyperledger/fabric/protos/peer"


)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
    chaincodePkg := "./test_cc.V1"
    chaincodeTar := "./test_cc.tar"
	ccbytes, err := ioutil.ReadFile(chaincodePkg)
	check(err)
	buf := ccbytes

	depSpec := &pb.ChaincodeDeploymentSpec{}
	err = proto.Unmarshal(buf, depSpec)
	check(err)
	fmt.Printf("chaincodeSpec=%+v\n", depSpec.ChaincodeSpec)

    payload := depSpec.CodePackage
    err = ioutil.WriteFile(chaincodeTar, payload, 0644)
    check(err)
}
