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
