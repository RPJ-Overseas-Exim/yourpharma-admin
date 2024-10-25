package types

type Customer struct{
    Id int
    Name, 
    Email,
    Number,
    Address string
}

type Order struct{
    OrderId int
    Name, 
    Email,
    Number,
    Product string
    Status StatusType
    Quantity,
    Price int
    Origin, 
    CreatedAt,
    Address string
}

// enum creation for status type
type StatusType int

const (
    NewOrder StatusType=iota
    PaymentDone
    Delievered
)

func (s StatusType) String() string{
    if(s == NewOrder){
        return "newOrder"
    }else if(s == PaymentDone){
        return "paymentDone"
    }else if(s == Delievered){
        return "delievered"
    }

    return "not a valid type"
}
