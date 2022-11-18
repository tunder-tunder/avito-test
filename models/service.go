package models

import (
	"fmt"
	"net/http"
)

type Service struct {
	ID           int    `json:"Service_id"`
	ServiceName  string `json:"Service_Name"`
	Price        int    `json:"Price"`
	Availability bool   `json:"availability"`
}
type ServiceList struct {
	Services []Service `json:"services"`
}

func (s *Service) Bind(r *http.Request) error {
	if s.ServiceName == "" {
		return fmt.Errorf("first name is a required field")
	}
	return nil
}

func (*ServiceList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Service) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
