package types

import "time"

type Product struct {
    Id      string
    Name    string
    Qty     int
    Price   int
}

type Customer struct{
    Id,
    Name, 
    Email,
    Address string
    Number  *int
}

type Order struct{
    Id string
    Name, 
    Email,
    Product     string
    Number      *int
    Status      StatusType
    Quantity,
    Price       int
    Origin, 
    Address     string
    CreatedAt,   
    UpdatedAt   time.Time
}

// enum creation for status type
type StatusType int

const (
    Active StatusType=iota
    Paid
    Shipped
    Delivered
)

func (s StatusType) String() string{
    if(s == Active){
        return "active"
    }else if(s == Paid){
        return "paid"
    }else if(s == Shipped){
        return "shipped"
    }else if(s == Delivered){
        return "delivered"
    }

    return "not a valid type"
}
