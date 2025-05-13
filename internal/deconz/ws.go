// Package deconz provides interfaces and types for interacting with the deCONZ REST API.
package deconz

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// RessourceType represents the type of resource in the deCONZ ecosystem.
// This is used to categorize different types of resources in WebSocket events.
type RessourceType string

// Constants defining the different resource types available in deCONZ.
const (
	// SensorsRessource represents sensor resources (motion sensors, temperature sensors, etc.)
	SensorsRessource RessourceType = "sensors"

	// ScenesRessource represents scene resources (predefined settings for groups of devices)
	ScenesRessource RessourceType = "scenes"

	// LightsRessource represents light resources (bulbs, strips, etc.)
	LightsRessource RessourceType = "lights"

	// GroupsRessource represents group resources (collections of lights or other devices)
	GroupsRessource RessourceType = "groups"
)

// EventType represents the type of event that occurred in the deCONZ ecosystem.
// This is used to categorize different types of events in WebSocket messages.
type EventType string

// Constants defining the different event types that can be received from deCONZ.
const (
	// AddedEvent indicates a new resource was added
	AddedEvent EventType = "added"

	// ChangedEvent indicates an existing resource was modified
	ChangedEvent EventType = "changed"

	// DeletedEvent indicates a resource was removed
	DeletedEvent EventType = "deleted"

	// SceneEvent indicates a scene was activated
	SceneEvent EventType = "scene-called"
)

// Messsage represents a WebSocket message from the deCONZ gateway.
// These messages provide real-time updates about changes in the Zigbee network.
// Different fields are populated depending on the event type and resource type.
type Messsage struct {
	// Type is the message type identifier
	Type string `json:"t"`

	// EventType indicates what kind of event occurred (added, changed, deleted, scene-called)
	EventType EventType `json:"e"`

	// RessourceType indicates what kind of resource the event relates to (sensors, lights, etc.)
	RessourceType RessourceType `json:"r"`

	// RessourceID is the identifier of the affected resource (not present for scene-called events)
	RessourceID *string `json:"id,omitempty"`

	// UniqueID is the unique identifier of the affected device (only for light and sensor resources)
	UniqueID *string `json:"uniqueid,omitempty"`

	// GroupID is the identifier of the group (only for scene-called events)
	GroupID *string `json:"gid,omitempty"`

	// SceneID is the identifier of the scene (only for scene-called events)
	SceneID *string `json:"scid,omitempty"`

	// Config contains configuration changes (only for changed events)
	Config *ObjectMap `json:"config,omitempty"`

	// Name contains the name change (only for changed events)
	Name *string `json:"name,omitempty"`

	// State contains state changes (only for changed events)
	State *ObjectMap `json:"state,omitempty"`

	// Group contains group information (only for added events)
	Group *interface{} `json:"group,omitempty"`

	// Light contains light information (only for added events)
	Light *Light `json:"light,omitempty"`

	// Sensor contains sensor information (only for added events)
	Sensor *interface{} `json:"sensor,omitempty"`
}

// EventClient manages a WebSocket connection to the deCONZ gateway.
// It receives real-time events about changes in the Zigbee network.
type EventClient struct {
	// client is the WebSocket connection to the deCONZ gateway
	client *websocket.Conn

	// done is a channel used to signal when the client should stop
	done chan struct{}
}

// NewEventClient creates a new WebSocket connection to the deCONZ gateway.
// It starts a goroutine that listens for events and processes them using the provided function.
//
// Parameters:
//   - ctx: Context for controlling the connection lifecycle
//   - path: The WebSocket URL to connect to
//   - eventFn: A function that will be called for each event received
//
// Returns:
//   - *EventClient: A pointer to the created EventClient
//   - error: Any error encountered during connection setup
func NewEventClient(ctx context.Context, path string, eventFn func(msg *Messsage)) (*EventClient, error) {
	ec := new(EventClient)

	// Establish the WebSocket connection
	c, _, err := websocket.DefaultDialer.DialContext(ctx, path, nil)
	if err != nil {
		log.Printf("[Events] websocket connection error: %+v", err)
		return nil, err
	}
	ec.client = c

	// Create a channel for signaling when to stop
	ec.done = make(chan struct{})

	// Start a goroutine to listen for events
	go func() {
		defer close(ec.done)
		for {
			// Read the next message from the WebSocket
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("[Events] websocket read error: %+v", err)
				continue
			}

			// Parse the message into a Messsage struct
			eventMsg := new(Messsage)
			if err := json.Unmarshal(message, eventMsg); err != nil {
				log.Printf("[Events] message unmarshal error: %+v", err)
				continue
			}

			// Process the event using the provided function
			eventFn(eventMsg)
		}
	}()

	return ec, nil
}

// Stop closes the WebSocket connection and stops the event processing goroutine.
//
// Returns:
//   - error: Any error encountered while closing the connection
func (ec *EventClient) Stop() error {
	close(ec.done)
	return ec.client.Close()
}
