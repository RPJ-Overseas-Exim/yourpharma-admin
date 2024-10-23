package dataset

type StatusType int 
const (
    Active StatusType = iota
    Paid
    Shipped
    Delivered
)
func (s StatusType) String() string{
    return [...]string{"Active", "Paid", "Shipped", "Delivered"}[s]
}
func (s StatusType) EnumIndex() int{
    return int(s)
}

type Customer struct{
    OrderId int
    Name, 
    Email,
    PhoneNo,
    Product string
    Status StatusType
    Quantity,
    Price int
    Origin, 
    CreatedAt,
    Address string
}
var Customers = []Customer {
    {
        OrderId: 1,
        Name: "Alan",
        Email: "alan@gmail.com",
        PhoneNo: "89756513213",
        Product: "Product 1",
        Status: Active,
        Quantity: 90,
        Price: 280,
        Origin: "Dash",
        CreatedAt: "date and time",
        Address: "Someplace",
    },
    {
        OrderId: 2,
        Name: "Alan2",
        Email: "alan@gmail.com",
        PhoneNo: "8513213",
        Product: "Product 1",
        Status: Active,
        Quantity: 90,
        Price: 280,
        Origin: "Dash",
        CreatedAt: "date and time",
        Address: "Someplace2",
    },
    {
        OrderId: 3,
        Name: "Alan3",
        Email: "alan@gmail.com",
        PhoneNo: "856513213",
        Product: "Product 1",
        Status: Delivered,
        Quantity: 90,
        Price: 280,
        Origin: "Dash",
        CreatedAt: "date and time",
        Address: "Someplace3",
    },
    {
        OrderId: 4,
        Name: "Alan4",
        Email: "alan@gmail.com",
        PhoneNo: "8756513213",
        Product: "Product 1",
        Status: Shipped,
        Quantity: 90,
        Price: 280,
        Origin: "Dash",
        CreatedAt: "date and time",
        Address: "Someplace4",
    },
    {
        OrderId: 5,
        Name: "Alan5",
        Email: "alan@gmail.com",
        PhoneNo: "8956513213",
        Product: "Product 1",
        Status: Active,
        Quantity: 90,
        Price: 280,
        Origin: "Dash",
        CreatedAt: "date and time",
        Address: "Someplace5",
    },
    {
        OrderId: 6,
        Name: "Alan6",
        Email: "alan@gmail.com",
        PhoneNo: "9756513213",
        Product: "Product 1",
        Status: Paid,
        Quantity: 90,
        Price: 280,
        Origin: "Web",
        CreatedAt: "date and time",
        Address: "Someplace6",
    },
}
