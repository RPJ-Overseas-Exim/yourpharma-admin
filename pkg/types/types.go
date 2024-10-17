package types

type Customer struct{
    OrderId string
    Name, Email, Product string
    Quantity, Price int
}

type Order struct{
    Id, Name, Email, Product string
    Quantity, Price int
    Status StatusType
}

// enum creation for status type
type StatusType int

const (
    newOrder StatusType=iota
    paymentDone
    delievered
)

func (s StatusType) String() string{
    if(s == newOrder){
        return "newOrder"
    }else if(s == paymentDone){
        return "paymentDone"
    }else if(s == delievered){
        return "delievered"
    }

    return "not a valid type"
}
