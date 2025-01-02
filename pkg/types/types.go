package types

import "time"

type Product struct {
	Id    string `json:"id"`
	PId   string `json:"product_id"`
	Name  string `json:"name"`
	Qty   int    `json:"quantity"`
	Price int    `json:"price"`
}

type Customer struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Number  *int   `json:"number"`
}

type Order struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Product   string    `json:"product"`
	Number    *int      `json:"number"`
	Status    string    `json:"status"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	Origin    string    `json:"origin"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// enum creation for status type
type StatusType int

const (
	Active StatusType = iota
	Paid
	Shipped
	Delivered
)

func (s StatusType) String() string {
	if s == Active {
		return "active"
	} else if s == Paid {
		return "paid"
	} else if s == Shipped {
		return "shipped"
	} else if s == Delivered {
		return "delivered"
	}

	return "not a valid type"
}
