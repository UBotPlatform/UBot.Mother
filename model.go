package main

import "time"

type AccountInfoFromRouter struct {
	ID              string            `json:"id"`
	Token           string            `json:"token"`
	ManagerMetadata map[string]string `json:"manager_metadata"`
}
type AppInfoFromRouter struct {
	ID              string            `json:"id"`
	Token           string            `json:"token"`
	ManagerMetadata map[string]string `json:"manager_metadata"`
}
type ServiceInfoFromMother struct {
	ID       string    `json:"id"`
	Status   int32     `json:"status"`
	LaunchAt time.Time `json:"launch_at"`
}
type AccountInfoFromMother struct {
	ID            string                 `json:"id"`
	BindedService *ServiceInfoFromMother `json:"binded_service,omitempty"`
}
type AppInfoFromMother struct {
	ID            string                 `json:"id"`
	BindedService *ServiceInfoFromMother `json:"binded_service,omitempty"`
}
type AccountProviderInfoFromMother struct {
	ID string `json:"id"`
}
type SuccessResponse struct {
	Success bool `json:"success"`
}
