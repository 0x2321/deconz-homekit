// Package main is the entry point for the deCONZ HomeKit Bridge application.
// It establishes connections to the deCONZ gateway, retrieves device information,
// creates HomeKit accessories, and starts the HomeKit server to enable control
// of Zigbee devices through Apple HomeKit.
package main

import (
	"context"
	"deconz-homekit/internal/accessoryManager"
	"deconz-homekit/internal/client"
	"deconz-homekit/internal/deconz"
	"deconz-homekit/internal/kvStorage"
	"fmt"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/charmbracelet/log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// go:generate directive to run the device configuration generator
// This creates the necessary device configuration files from the JSON templates
//go:generate go run generateDeviceConfiguration.go

// main is the entry point of the application.
// It initializes the bridge, connects to the deCONZ gateway,
// retrieves device information, and starts the HomeKit server.
func main() {
	// Create a context that can be cancelled on system signals
	ctx := DefaultContext()

	// Initialize the logger with timestamp formatting
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})

	l.Info("Starting bridge...")

	// Initialize the key-value storage for persistent data
	STORAGE_PATH := os.Getenv("STORAGE_PATH")
	if len(STORAGE_PATH) == 0 {
		STORAGE_PATH = "./"
	}
	storage, err := kvStorage.New(STORAGE_PATH + "db.sqlite")
	if err != nil {
		l.Fatalf("Error connecting to the database: %v", err)
	}

	// Get deCONZ gateway IP address from environment variables
	var PHOSCON_IP = os.Getenv("DECONZ_IP")
	if len(PHOSCON_IP) == 0 {
		l.Fatalf("Please provide the ip address of the deCONZ gateway (DECONZ_IP not set)")
	}

	// Get deCONZ gateway port from environment variables, default to 80 if not set
	var PHOSCON_PORT = os.Getenv("DECONZ_PORT")
	if len(PHOSCON_PORT) == 0 {
		PHOSCON_PORT = "80"
	}

	// Retrieve or generate the deCONZ API key for authentication
	apiKeyRaw, err := storage.Get("deconz_api_key")
	if err != nil {
		l.Fatalf("Error querying to the database: %v", err)
	}

	// If no API key exists, request a new one from the deCONZ gateway
	if apiKeyRaw == nil {
		l.Infof("No API key found. Requesting a new one...")

		// Request a new API key from the deCONZ gateway
		apiKeyRaw, err = getApiKey(l, fmt.Sprintf("http://%s:%s", PHOSCON_IP, PHOSCON_PORT))
		if err != nil {
			l.Fatalf("Could not obtain API key: %v", err)
		}

		// Save the new API key to the storage for future use
		if err = storage.Set("deconz_api_key", apiKeyRaw); err != nil {
			l.Fatalf("Could not store API key: %v", err)
		}
	}

	// Connect to the deCONZ API and retrieve gateway configuration
	l.Info("Connecting to deCONZ gateway...")
	api := deconz.NewApiClient(fmt.Sprintf("http://%s:%s", PHOSCON_IP, PHOSCON_PORT), string(apiKeyRaw))
	config, err := api.GetConfiguration()
	if err != nil {
		l.Fatalf("Error getting configuration: %v", err)
	}

	// Retrieve all devices from the deCONZ gateway
	l.Info("Retrieving devices from deCONZ gateway...")
	devices, err := api.GetAllDevices()
	if err != nil {
		l.Fatalf("Failed to get all devices: %+v", err)
	}

	// Create HomeKit accessories for each supported device
	l.Info("Creating HomeKit accessories...")
	am := accessoryManager.NewAccessoryManager(api, devices)

	// Connect to the deCONZ WebSocket event stream for real-time updates
	l.Info("Connecting to deCONZ event stream...")
	_, err = deconz.NewEventClient(ctx, fmt.Sprintf("ws://%s:%d", PHOSCON_IP, config.WebsocketPort), am.ProcessUpdate)
	if err != nil {
		l.Fatalf("WebSocket connection error: %+v", err)
	}

	// Initialize and start the HomeKit server
	l.Info("Starting HomeKit server...")

	// Create a bridge accessory to represent the deCONZ gateway in HomeKit
	b := accessory.NewBridge(accessory.Info{
		Manufacturer: "deCONZ Bridge",
		Name:         fmt.Sprintf("%s %s", config.Name, strings.ReplaceAll(config.BridgeId[:4], ":", "")),
		SerialNumber: config.BridgeId,
		Model:        config.DeviceName,
		Firmware:     config.SwVersion,
	})

	// Create a new HomeKit server with the bridge and all device accessories
	server, err := hap.NewServer(storage, b.A, am.GetAccessories()...)
	if err != nil {
		l.Fatalf("HomeKit server initialization error: %+v", err)
	}

	// set port
	server.Addr = "0.0.0.0:51826"

	// Generate a random 8-digit pairing code for HomeKit setup
	code := uint32(rand.Intn(90000000) + 10000000)
	server.Pin = fmt.Sprintf("%d", code)
	l.Infof("HomeKit pairing code: %s-%s", server.Pin[0:4], server.Pin[4:8])

	// Start the HomeKit server and listen for connections
	if err := server.ListenAndServe(ctx); err != nil {
		l.Fatalf("HomeKit server error: %+v", err)
	}
}

// getApiKey requests and retrieves an API key from the deCONZ gateway.
// It repeatedly attempts to obtain the key until successful, prompting the user
// to press the link button on the gateway when necessary.
//
// Parameters:
//   - log: Logger for output messages
//   - addr: The base URL of the deCONZ gateway
//
// Returns:
//   - []byte: The API key as a byte slice
//   - error: Any error encountered during the process
func getApiKey(log *log.Logger, addr string) ([]byte, error) {
	// Define request and response types for the API key request
	type Request struct {
		DeviceName string `json:"devicetype"`
	}
	type Response []map[string]map[string]interface{}

	// Loop until an API key is successfully obtained
	for {
		// Send a POST request to the deCONZ API to request an API key
		data, err := client.Post[Response](addr+"/api", Request{DeviceName: "HomeKit Bridge"})
		if err != nil {
			// Return any HTTP or network errors
			return nil, err
		}

		// Parse the response to extract the API key
		if result, ok := (*data)[0]["success"]; ok {
			if username, ok := result["username"]; ok {
				log.Info("Successfully obtained an API key")
				return []byte(username.(string)), nil
			}
		}

		// If no key was provided (likely because the link button wasn't pressed),
		// wait and try again
		log.Warn("Please press the link button on your deCONZ gateway to obtain an API key. Retrying in 15s...")
		time.Sleep(15 * time.Second)
	}
}

// DefaultContext creates a context that can be cancelled when the application
// receives an interrupt or termination signal (SIGINT or SIGTERM).
//
// Returns:
//   - context.Context: A cancellable context tied to system signals
func DefaultContext() context.Context {
	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Register for interrupt (Ctrl+C) and termination signals
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine that will cancel the context when a signal is received
	go func() {
		<-c
		// Stop delivering signals to prevent multiple cancellations
		signal.Stop(c)
		// Cancel the context to initiate graceful shutdown
		cancel()
	}()

	return ctx
}
