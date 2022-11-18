package models

import (
	"net/http"
	"time"
)

type Balance struct {
	ID          int       `json:"id"`
	UserId      int       `json:"User_id"`
	Total       int       `json:"Total"`
	Reserve     int       `json:"Reserve"`
	OrderNumber string    `json:"Order_number"`
	ServiceId   string    `json:"Service_id"`
	CreatedAt   time.Time `json:"Created_at"`
	Status      string    `json:"Status"`
}

type BalanceList struct {
	Balances []Balance `json:"balances"`
}

func (b *Balance) Bind(r *http.Request) error {
	return nil
}

func (*BalanceList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
