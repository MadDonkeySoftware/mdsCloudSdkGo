package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MadDonkeySoftware/mdsCloudSdkGo/sdk"
)

type testCredentials struct {
	AccountID string
	UserName  string
	Password  string
}

func createTestAccount(client *sdk.IdentityClient) *testCredentials {
	uniqueTestName := fmt.Sprintf("Test-%d", time.Now().Unix())
	testCreds := &testCredentials{
		UserName: uniqueTestName,
		Password: "Password",
	}

	registerArgs := sdk.RegisterAccountArgs{
		AccountName:  uniqueTestName,
		Email:        fmt.Sprintf("%s@no.com", uniqueTestName),
		FriendlyName: uniqueTestName,
		Password:     testCreds.Password,
		UserID:       uniqueTestName,
	}

	fmt.Println("== Register ==")
	registerResult, err := client.Register(&registerArgs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AccountID: %s\n", registerResult.AccountID)
	fmt.Printf("Status: %s\n", registerResult.Status)
	fmt.Println()

	fmt.Println("== Authenticate ==")
	authResult, err := client.Authenticate(&sdk.AuthenticateArgs{
		AccountID: registerResult.AccountID,
		Password:  "Password",
		UserID:    uniqueTestName,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s\n", authResult.Token)
	fmt.Println()

	testCreds.AccountID = registerResult.AccountID
	return testCreds
}

func testServerlessFunctions(client *sdk.ServerlessFunctionsClient) {
	fmt.Println("== List ==")
	funcList, err := client.ListFunctions()
	if err != nil {
		panic(err)
	}

	for _, summary := range *funcList {
		fmt.Printf("Name: %s\n", summary.Name)
		fmt.Printf("Orid: %s\n", summary.Orid)
	}
	fmt.Println()

	fmt.Println("== Create ==")
	sourcePath := "/home/malline/git-source/mdsCloudFnProjectMinion/mdsCloudServerlessFunctions-sampleApp.zip"
	funcSummary, err := client.CreateFunction("test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Name: %s\n", funcSummary.Name)
	fmt.Printf("Orid: %s\n", funcSummary.Orid)
	fmt.Println()

	fmt.Println("== Updating Code ==")
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

	fmt.Println("== Details ==")
	funcDetails, err := client.GetFunctionDetails(funcSummary.Orid)
	if err != nil {
		panic(err)
	}
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

func testIdentityClient(client *sdk.IdentityClient, testCreds *testCredentials) {
	fmt.Println("== Update User ==")
	err := client.UpdateUser(&sdk.UpdateUserArgs{
		// AccountID: registerResult.AccountID,
		// Password:  "Password",
		// UserID:    uniqueTestName,
		FriendlyName: fmt.Sprintf("%s-updated", testCreds.UserName),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Update Complete ==")
	fmt.Println()
}

func testQueueServiceClient(client *sdk.QueueServiceClient) {
	// Create
	fmt.Println("== Create Queue ==")
	createResult, err := client.CreateQueue(&sdk.CreateQueueArgs{
		Name:     "TestQueue",
		Resource: "Foo Bar",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Name: %s\n", createResult.Name)
	fmt.Printf("Orid: %s\n", createResult.Orid)
	fmt.Printf("Status: %s\n", createResult.Status)
	fmt.Println()

	// Read
	fmt.Println("== Read Queue ==")
	readResult, err := client.GetQueueDetails(&sdk.GetQueueDetailsArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orid: %s\n", readResult.Orid)
	fmt.Printf("Resource: %s\n", readResult.Resource)
	fmt.Println()

	// Modify
	fmt.Println("== Modify Queue ==")
	err = client.UpdateQueue(&sdk.UpdateQueueArgs{
		Orid:     createResult.Orid,
		Resource: "NULL",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Modify Completed ==")
	fmt.Println()

	fmt.Println("== Read 2 Queue ==")
	readResult2, err := client.GetQueueDetails(&sdk.GetQueueDetailsArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orid: %s\n", readResult2.Orid)
	fmt.Printf("Resource: %s\n", readResult2.Resource)
	fmt.Println()

	// Delete
	fmt.Println("== Delete Queue ==")
	err = client.DeleteQueue(&sdk.DeleteQueueArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Delete Completed ==")
	fmt.Println()
}

func testFileServiceClient(client *sdk.FileServiceClient) {
	// Create
	fmt.Println("== Create Container ==")
	createResult, err := client.CreateContainer(&sdk.CreateContainerArgs{
		Name: "TestContainer",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orid: %s\n", createResult.Orid)
	fmt.Println()

	// Read
	fmt.Println("== Read Container Contents ==")
	readResult, err := client.ListContainerContents(&sdk.ListContainerContentsArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Directories: %s\n", readResult.Directories)
	fmt.Printf("Files: %s\n", readResult.Files)
	fmt.Println()

	// Delete
	fmt.Println("== Delete Container ==")
	err = client.DeleteContainerOrPath(&sdk.DeleteContainerArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Delete Completed ==")
	fmt.Println()
}

const definition string = `{
  "Name": "test",
  "StartsAt": "Success",
  "States": {
    "Success": {
      "Type": "Succeed"
    }
  }
}`

const definition2 string = `{
  "Name": "test",
  "StartsAt": "Wait",
  "States": {
    "Wait": {
      "Type": "Wait",
      "Seconds": "5",
      "Next": "Success"
    },
    "Success": {
      "Type": "Succeed"
    }
  }
}`

func testStateMachineServiceClient(client *sdk.StateMachineServiceClient) {
	// Create
	fmt.Println("== Create State Machine ==")
	createResult, err := client.CreateStateMachine(&sdk.CreateStateMachineArgs{
		Definition: definition,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orid: %s\n", createResult.Orid)
	fmt.Println()

	fmt.Println("== Update State Machine ==")
	updateResult, err := client.UpdateStateMachine(&sdk.UpdateStateMachineArgs{
		Orid:       createResult.Orid,
		Definition: definition2,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Update Completed ==")
	fmt.Printf("Orid: %s\n", updateResult.Orid)
	fmt.Println()

	// Read
	fmt.Println("== Read State Machine ==")
	readResult, err := client.GetStateMachineDetails(&sdk.GetStateMachineDetailsArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orid: %s\n", readResult.Orid)
	fmt.Printf("Name: %s\n", readResult.Name)

	b, _ := json.Marshal(readResult.Definition)
	fmt.Printf("Definition: %s\n", string(b))
	fmt.Println()

	// Delete
	fmt.Println("== Delete State Machine ==")
	deleteResult, err := client.DeleteStateMachine(&sdk.DeleteStateMachineArgs{
		Orid: createResult.Orid,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("== Delete Completed ==")
	fmt.Printf("Orid: %s\n", deleteResult.Orid)
	fmt.Println()
}

func main() {
	enableAuthSemaphore := false
	allowSelfCert := true
	urls := map[string]string{
		"identityUrl": "https://127.0.0.1:8081",
		"nsUrl":       "http://127.0.0.1:8082",
		"qsUrl":       "http://127.0.0.1:8083",
		"fsUrl":       "http://127.0.0.1:8084",
		"sfUrl":       "http://127.0.0.1:8085",
		"smUrl":       "http://127.0.0.1:8086",
	}
	sdkObj := sdk.NewSdk("", "", "", allowSelfCert, enableAuthSemaphore, urls)
	testCreds := createTestAccount(sdkObj.GetIdentityClient())

	sdkObj = sdk.NewSdk(testCreds.AccountID, testCreds.UserName, testCreds.Password, allowSelfCert, enableAuthSemaphore, urls)
	testIdentityClient(sdkObj.GetIdentityClient(), testCreds)
	// testServerlessFunctions(sdkObj.GetServerlessFunctionsClient())
	testQueueServiceClient(sdkObj.GetQueueServiceClient())
	// testFileServiceClient(sdkObj.GetFileServiceClient())
	// testStateMachineServiceClient(sdkObj.GetStateMachineServiceClient())
}
