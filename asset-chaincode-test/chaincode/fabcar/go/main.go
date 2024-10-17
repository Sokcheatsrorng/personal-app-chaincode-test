package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Car represents a car object in the chaincode
type Car struct {
    AssetID string `json:"assetID"`
    Make    string `json:"make"`
    Model   string `json:"model"`
    Color   string `json:"color"`
    Owner   string `json:"owner"`
}

// SmartContract provides functions for managing Cars
type SmartContract struct {
    contractapi.Contract
}

// InitLedger adds some default cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    cars := []Car{
        {AssetID: "CAR0", Make: "Toyota", Model: "Prius", Color: "blue", Owner: "Tomoko"},
        {AssetID: "CAR1", Make: "Ford", Model: "Mustang", Color: "red", Owner: "Brad"},
        {AssetID: "CAR2", Make: "Hyundai", Model: "Tucson", Color: "green", Owner: "Jin Soo"},
        {AssetID: "CAR3", Make: "Volkswagen", Model: "Passat", Color: "yellow", Owner: "Max"},
        {AssetID: "CAR4", Make: "Tesla", Model: "Model S", Color: "black", Owner: "Adriana"},
        {AssetID: "CAR5", Make: "Peugeot", Model: "208", Color: "purple", Owner: "Michel"},
    }

    for _, car := range cars {
        carAsBytes, err := json.Marshal(car)
        if err != nil {
            return fmt.Errorf("failed to marshal car: %v", err)
        }

        err = ctx.GetStub().PutState(car.AssetID, carAsBytes)
        if err != nil {
            return fmt.Errorf("failed to put car into world state: %v", err)
        }
    }

    return nil
}

// CreateCar creates a new car asset with the given details and stores the reference to off-chain storage
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, assetID string, make string, model string, color string, owner string) error {
    car := Car{
        AssetID: assetID, // Reference to off-chain storage
        Make:    make,
        Model:   model,
        Color:   color,
        Owner:   owner,
    }

    // Serialize the car object to JSON
    carAsBytes, err := json.Marshal(car)
    if err != nil {
        return fmt.Errorf("failed to marshal car: %v", err)
    }

    // Put the serialized car into the world state
    err = ctx.GetStub().PutState(assetID, carAsBytes)
    if err != nil {
        return fmt.Errorf("failed to put car into world state: %v", err)
    }

    // Emit an event to notify external systems (off-chain processing)
    err = ctx.GetStub().SetEvent("CarCreated", []byte(assetID))
    if err != nil {
        return fmt.Errorf("failed to emit event: %v", err)
    }

    return nil
}

// QueryCar returns the car stored in the world state with given assetID
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
    log.Println("Querying all cars")

    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        log.Printf("Error getting all cars: %v", err)
        return nil, fmt.Errorf("failed to get all cars: %v", err)
    }
    defer resultsIterator.Close()

    var cars []*Car
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            log.Printf("Error iterating cars: %v", err)
            return nil, fmt.Errorf("failed to iterate cars: %v", err)
        }

        var car Car
        err = json.Unmarshal(queryResponse.Value, &car)
        if err != nil {
            log.Printf("Error unmarshaling car: %v", err)
            return nil, fmt.Errorf("failed to unmarshal car: %v", err)
        }
        log.Printf("Found car: %+v", car)
        cars = append(cars, &car)
    }

    log.Printf("Total cars found: %d", len(cars))
    return cars, nil
}


func main(){
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("Error create fabcar chaincode: %s", err)
        return
    }
    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting fabcar chaincode: %s", err)
    }
}