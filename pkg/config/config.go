package config

import (
	"github.com/MIKE9708/s4t-sdk-go/pkg"
)


type ProviderConfig struct {
	S4tClient *s4t.Client
}


func SetUpProvider() (*ProviderConfig, error){
	c := s4t.Client{}
	client, err := c.GetClientConnection() 
	
	if err != nil {
		return nil, err
	}

	return &ProviderConfig{
		S4tClient: client,
	}, nil

} 
