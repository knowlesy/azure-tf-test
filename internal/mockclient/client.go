package mockclient

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Resource represents a generic mock resource saved in the database
type Resource struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Client simulates the Azure Resource Manager API locally
type Client struct {
	mu       sync.RWMutex
	FilePath string
	State    map[string]Resource // Key is Resource ID
}

// NewClient initializes the Mock Client, loading from file if exists
func NewClient(filePath string) (*Client, error) {
	c := &Client{
		FilePath: filePath,
		State:    make(map[string]Resource),
	}

	// Try to load existing state
	data, err := os.ReadFile(filePath)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.State)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mock db: %w", err)
		}
	} else if !os.IsNotExist(err) && err != nil {
		return nil, fmt.Errorf("failed to read mock db: %w", err)
	}

	return c, nil
}

// Save writes a resource to the mock DB and persists to disk
func (c *Client) Save(resourceType, id string, properties map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Normalizing ID to lowercase for consistent lookup
	normalizedID := strings.ToLower(id)

	c.State[normalizedID] = Resource{
		ID:         id,
		Type:       resourceType,
		Properties: properties,
	}

	return c.persist()
}

// Read retrieves a resource by ID. Returns false if not found.
func (c *Client) Read(id string) (Resource, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res, ok := c.State[strings.ToLower(id)]
	return res, ok
}

// Delete removes a resource from the mock DB
func (c *Client) Delete(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.State, strings.ToLower(id))
	return c.persist()
}

// FindResourcesByTypeAndProperty searches for resources of a certain type that have a matching property value.
// E.g., searching for NetworkInterfaces with a specific subnet ID
func (c *Client) FindResourcesByTypeAndProperty(resourceType, propertyKey, expectedValue string) []Resource {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var results []Resource
	for _, res := range c.State {
		if res.Type == resourceType {
			// Basic property lookup. We might need nested lookups depending on how we structure the properties map
			if val, ok := res.Properties[propertyKey]; ok {
				if strVal, ok := val.(string); ok && strings.EqualFold(strVal, expectedValue) {
					results = append(results, res)
				}
			}
		}
	}
	return results
}

// persist writes the current state back to the JSON file
func (c *Client) persist() error {
	data, err := json.MarshalIndent(c.State, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mock db: %w", err)
	}
	return os.WriteFile(c.FilePath, data, 0644)
}

// GenerateResourceID creates a standard Azure Resource ID
func GenerateResourceID(subscriptionID, resourceGroup, provider, resourceType, resourceName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s",
		subscriptionID, resourceGroup, provider, resourceType, resourceName)
}

// GenerateSubnetID creates an Azure Subnet ID
func GenerateSubnetID(vnetID, subnetName string) string {
	return fmt.Sprintf("%s/subnets/%s", vnetID, subnetName)
}
