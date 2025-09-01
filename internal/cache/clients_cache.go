package cache

import "go.mau.fi/whatsmeow"

var clients = make(map[string]*whatsmeow.Client)

func GetClient(key string) *whatsmeow.Client {
	return clients[key]
}

func SetClient(key string, client *whatsmeow.Client) {
	clients[key] = client
}

func SetAllClients(newClients map[string]*whatsmeow.Client) {
	clients = newClients
}

func GetAllClients() map[string]*whatsmeow.Client {
	return clients
}

// DeleteClient removes a client from the cache by ID
func DeleteClient(key string) bool {
	if client, exists := clients[key]; exists {
		// Disconnect the client if it's connected
		if client.IsConnected() {
			client.Disconnect()
		}
		delete(clients, key)
		return true
	}
	return false
}

// ClientExists checks if a client exists in the cache
func ClientExists(key string) bool {
	_, exists := clients[key]
	return exists
}

// GetClientCount returns the number of clients in the cache
func GetClientCount() int {
	return len(clients)
}
