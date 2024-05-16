package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	//"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func main() {
	// export AZURE_COSMOS_ENDPOINT=""
	// export AZURE_COSMOS_KEY=""

	//	client, err := createClientAAD()
	//	if err != nil {
	//		log.Fatal("Failed to create client: ", err)
	//	}

	endpoint := os.Getenv("AZURE_COSMOS_ENDPOINT")
	if endpoint == "" {
		log.Fatal("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key := os.Getenv("AZURE_COSMOS_KEY")
	if key == "" {
		log.Fatal("AZURE_COSMOS_KEY could not be found")
	}

	var databaseName = "adventureworks"
	var containerName = "customer"
	var partitionKeyPath = "/customerId"

	item := struct {
		ID           string `json:"id"`
		CustomerId   string `json:"customerId"`
		Title        string
		FirstName    string
		LastName     string
		EmailAddress string
		PhoneNumber  string
		CreationDate string
	}{
		ID:           "1",
		CustomerId:   "1",
		Title:        "Mr",
		FirstName:    "Luke",
		LastName:     "Hayes",
		EmailAddress: "luke12@adventure-works.com",
		PhoneNumber:  "879-555-0197",
		CreationDate: time.Now().UTC().Format("2006-01-02T15:04:05.1234567Z"), //cosmos db recommended date format
	}

	// Map
	/* item := map[string]interface{}{
		"id":           "1",
		"customerId":   "1",
		"title":        "Mr",
		"firstName":    "Luke",
		"lastName":     "Hayes",
		"emailAddress": "luke12@adventure-works.com",
		"phoneNumber":  "879-555-0197",
	}  */

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatal("Failed to create a credential: ", err)
	}

	// Create a CosmosDB client
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatal("Failed to create cosmos db client: ", err)
	}

	// CRUD functions
	err = createDatabase(client, databaseName)
	if err != nil {
		log.Printf("createDatabase failed: %s\n", err)
	}

	//
	err = createContainer(client, databaseName, containerName, partitionKeyPath)
	if err != nil {
		log.Printf("createContainer failed: %s\n", err)
	}

	err = createItem(client, databaseName, containerName, item.CustomerId, item)
	if err != nil {
		log.Printf("createItem failed: %s\n", err)
	}

	err = readItem(client, databaseName, containerName, item.CustomerId, item.ID)
	if err != nil {
		log.Printf("readItem failed: %s\n", err)
	}

	err = replaceItem(client, databaseName, containerName, item.CustomerId, item.ID, item)
	if err != nil {
		log.Printf("replaceItem failed: %s\n", err)
	}

	err = upsertItem(client, databaseName, containerName, item.CustomerId, item.ID, item)
	if err != nil {
		log.Printf("upsertItem failed: %s\n", err)
	}

	err = deleteItem(client, databaseName, containerName, item.CustomerId, item.ID)
	if err != nil {
		log.Printf("deleteItem failed: %s\n", err)
	}
	err = readMany(client, databaseName, containerName, item.CustomerId)
	if err != nil {
		log.Printf("readMany failed: %s\n", err)
	}
}

func createDatabase(client *azcosmos.Client, databaseName string) error {
	//	databaseName := "adventureworks"

	// This is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}
	databaseProperties := azcosmos.DatabaseProperties{ID: databaseName}
	options := &azcosmos.CreateDatabaseOptions{}
	ctx := context.TODO()
	databaseResp, err := client.CreateDatabase(ctx, databaseProperties, options)

	switch {
	case errorIs409(err):
		log.Printf("Database [%s] already exists\n", databaseName)
	case err != nil:
		return err
	default:
		log.Printf("Database [%v] created. ActivityId %s\n", databaseName, databaseResp.ActivityID)
	}
	return nil
}

func createContainer(client *azcosmos.Client, databaseName, containerName, partitionKeyPath string) error {
	//	databaseName = adventureworks
	//	containerName = customer
	//	partitionKeyPath = "/customerId"

	databaseClient, err := client.NewDatabase(databaseName)
	if err != nil {
		return err
	}

	// creating a container
	containerProperties := azcosmos.ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{partitionKeyPath},
		},
	}

	// this is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}

	// setting options upon container creation
	throughputProperties := azcosmos.NewManualThroughputProperties(400) //defaults to 400 if not set
	options := &azcosmos.CreateContainerOptions{
		ThroughputProperties: &throughputProperties,
	}
	ctx := context.TODO()
	containerResponse, err := databaseClient.CreateContainer(ctx, containerProperties, options)

	switch {
	case errorIs409(err):
		log.Printf("Container [%s] already exists\n", containerName)
	case err != nil:
		return err
	default:
		log.Printf("Container [%s] created. ActivityId %s\n", containerName, containerResponse.ActivityID)
	}
	return nil
}

func createItem(client *azcosmos.Client, databaseName, containerName, partitionKey string, item any) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "1"

	/*	item = struct {
			ID           string `json:"id"`
			CustomerId   string `json:"customerId"`
			Title        string
			FirstName    string
			LastName     string
			EmailAddress string
			PhoneNumber  string
			CreationDate string
		}{
			ID:           "1",
			CustomerId:   "1",
			Title:        "Mr",
			FirstName:    "Luke",
			LastName:     "Hayes",
			EmailAddress: "luke12@adventure-works.com",
			PhoneNumber:  "879-555-0197",
			CreationDate: "2014-02-25T00:00:00",
		}
	*/
	// create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client: %s", err)
	}

	// specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	// setting the item options upon creating ie. consistency level
	itemOptions := azcosmos.ItemOptions{
		ConsistencyLevel: azcosmos.ConsistencyLevelSession.ToPtr(),
	}

	// this is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}

	ctx := context.TODO()
	itemResponse, err := containerClient.CreateItem(ctx, pk, b, &itemOptions)

	switch {
	case errorIs409(err):
		log.Printf("Item with partitionkey value %s already exists\n", pk)
	case err != nil:
		return err
	default:
		log.Printf("Status %d. Item %v created. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)
	}

	return nil
}

func readItem(client *azcosmos.Client, databaseName, containerName, partitionKey, itemId string) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "1"
	//	itemId = "1"

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client: %s", err)
	}

	// Specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Read an item
	ctx := context.TODO()
	itemResponse, err := containerClient.ReadItem(ctx, pk, itemId, nil)
	if err != nil {
		return err
	}

	itemResponseBody := struct {
		ID           string `json:"id"`
		CustomerId   string `json:"customerId"`
		Title        string
		FirstName    string
		LastName     string
		EmailAddress string
		PhoneNumber  string
		CreationDate string
	}{}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(itemResponseBody, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("Read item with customerId %s\n", itemResponseBody.CustomerId)
	fmt.Printf("%s\n", b)

	log.Printf("Status %d. Item %v read. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)

	return nil
}

func deleteItem(client *azcosmos.Client, databaseName, containerName, partitionKey, itemId string) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "1"
	//	itemId = "1"

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client:: %s", err)
	}
	// Specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Delete an item
	ctx := context.TODO()

	res, err := containerClient.DeleteItem(ctx, pk, itemId, nil)
	if err != nil {
		return err
	}

	log.Printf("Status %d. Item %v deleted. ActivityId %s. Consuming %v Request Units.\n", res.RawResponse.StatusCode, pk, res.ActivityID, res.RequestCharge)

	return nil
}

func replaceItem(client *azcosmos.Client, databaseName, containerName, partitionKey, itemID string, item any) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "1"
	//	itemId = "1"

	/*	item = struct {
			ID           string `json:"id"`
			CustomerId   string `json:"customerId"`
			Title        string
			FirstName    string
			LastName     string
			EmailAddress string
			PhoneNumber  string
			CreationDate string
		}{
			ID:           "1",
			CustomerId:   "2",
			Title:        "Mr",
			FirstName:    "Luke",
			LastName:     "Hayes",
			EmailAddress: "luke12@adventure-works.com",
			PhoneNumber:  "879-555-0197",
			CreationDate: "2014-02-25T00:00:00",
		}
	*/
	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client:: %s", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)

	//b, err := json.Marshal(item)
	//if err != nil {
	//	log.Fatal(err)
	//}

	b, err := json.MarshalIndent(item, "", "    ")
	if err != nil {
		return err
	}

	ctx := context.TODO()
	itemResponse, err := containerClient.ReplaceItem(ctx, pk, itemID, b, nil /*options*/)
	if err != nil {
		var responseErr *azcore.ResponseError
		errors.As(err, &responseErr)
		log.Fatal(responseErr)
	}

	fmt.Printf("Updated item with customerId %s\n", partitionKey)
	fmt.Printf("%s\n", b)

	log.Printf("Status %d. Item %v updated ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)

	return nil
}

func upsertItem(client *azcosmos.Client, databaseName, containerName, partitionKey, itemID string, item any) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "2"
	//	itemId = "1"

	/*	item = struct {
			ID           string `json:"id"`
			CustomerId   string `json:"customerId"`
			Title        string
			FirstName    string
			LastName     string
			EmailAddress string
			PhoneNumber  string
			CreationDate string
		}{
			ID:           "1",
			CustomerId:   "2",
			Title:        "Mr",
			FirstName:    "Luke",
			LastName:     "Hayes",
			EmailAddress: "luke12@adventure-works.com",
			PhoneNumber:  "879-555-0197",
			CreationDate: "2014-02-25T00:00:00",
		}
	*/
	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client:: %s", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)

	b, err := json.MarshalIndent(item, "", "    ")
	if err != nil {
		return err
	}

	ctx := context.TODO()
	itemResponse, err := containerClient.UpsertItem(ctx, pk, b, nil /*options*/)
	if err != nil {
		var responseErr *azcore.ResponseError
		errors.As(err, &responseErr)
		log.Fatal(responseErr)
	}

	fmt.Printf("Upserted item with customerId %s\n", partitionKey)
	fmt.Printf("%s\n", b)

	log.Printf("Status %d. Item %v upserted ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)

	return nil
}

func readMany(client *azcosmos.Client, databaseName, containerName, partitionKey string) error {
	customerId := "2"
	pk := azcosmos.NewPartitionKeyString(customerId)

	container, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return err
	}
	log.Printf("Printing customerId [%v] in [%v\\%v]", pk, databaseName, containerName)

	query := "SELECT * " +
		"FROM c WHERE c.customerId in ('2','3')"

	fmt.Printf("Print out customers\n")

	queryPager := container.NewQueryItemsPager(query, pk, &azcosmos.QueryOptions{PopulateIndexMetrics: true})
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())
		if err != nil {
			return err
		}
		for _, item := range queryResponse.Items {
			map1 := map[string]interface{}{}
			err := json.Unmarshal(item, &map1)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(map1, "", "    ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
		}
		log.Printf("Query page received with %d items. Status %d. ActivityId %s. Consuming %v RU\n", len(queryResponse.Items), queryResponse.RawResponse.StatusCode, queryResponse.ActivityID, queryResponse.RequestCharge)
	}
	return nil
}

/*
func createClientAAD() (*azcosmos.Client, error) {
	var config = newConfig()
	var endpoint = config.endpoint

	cred, err := azidentity.NewClientSecretCredential (
		config.tenantID,
		config.clientID,
		config.clientSecret,
		nil /*options*)
	if err != nil {
		return nil, err
	}

	client, err := azcosmos.NewClient(endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
*/
