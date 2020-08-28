package main

import (
	"fmt"

	"github.com/MadDonkeySoftware/mdsCloudSdkGo/sdk"
)

func testServerlessFunctions(client *sdk.ServerlessFunctionsClient) {
	sourcePath := "/home/malline/git-source/mdsCloudFnProjectMinion/mdsCloudServerlessFunctions-sampleApp.zip"
	funcSummary, err := client.CreateFunction("test")
	if err != nil {
		panic(err)
	}
	fmt.Println("== Create ==")
	fmt.Printf("Name: %s\n", funcSummary.Name)
	fmt.Printf("Orid: %s\n", funcSummary.Orid)
	fmt.Println()

	fmt.Println("== Updateing Code ==")
	err = client.UpdateFunctionCode(funcSummary.Orid, "node", "src/one:main", sourcePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("== Update Complete ==")
	fmt.Println()

	fmt.Println("== Invoking Code ==")
	testBody := map[string]interface{}{"name": "Frito"}
	result, err := client.InvokeFunction(funcSummary.Orid, testBody)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %s\n", result)
	fmt.Println("== Invoke Complete ==")
	fmt.Println()

	funcDetails, err := client.GetFunctionDetails(funcSummary.Orid)
	if err != nil {
		panic(err)
	}
	fmt.Println("== Details ==")
	fmt.Printf("Name: %s\n", funcDetails.Name)
	fmt.Printf("Orid: %s\n", funcDetails.Orid)
	fmt.Printf("Version: %s\n", funcDetails.Version)
	fmt.Printf("Runtime: %s\n", funcDetails.Runtime)
	fmt.Printf("EntryPoint: %s\n", funcDetails.EntryPoint)
	fmt.Printf("Created: %s\n", funcDetails.Created)
	fmt.Printf("LastInvoke: %s\n", funcDetails.LastInvoke)
	fmt.Printf("LastUpdate: %s\n", funcDetails.LastUpdate)
	fmt.Println()

	// Cleanup
	err = client.DeleteFunction(funcSummary.Orid)
	if err != nil {
		panic(err)
	}
	fmt.Println("== Delete ==")
	fmt.Println(fmt.Sprintf("Successfully deleted: %s", funcSummary.Orid))
	fmt.Println()
}

func main() {
	sdkObj := sdk.NewSdk("1", map[string]string{"sfUrl": "http://192.168.5.90:8086"})
	client := sdkObj.GetServerlessFunctionsClient("", "")

	testServerlessFunctions(client)
}
