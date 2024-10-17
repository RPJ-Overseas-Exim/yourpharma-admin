package adminHandlers

import (
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/templ/adminViews"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
    comp := adminViews.HomeIndex("Home", adminViews.Home())
    return comp.Render(c.Request().Context(), c.Response().Writer)
}

func Customers(c echo.Context) error {
    data := []types.Customer{
        {
            OrderId: "ordahfs",
            Name: "Muzammil",
            Email: "email@gmail.com",
            Product: "product1",
            Quantity: 180,
            Price: 480,
        },
        {
            OrderId: "orsddahfs",
            Name: "Muzammil2",
            Email: "email@gmail.com",
            Product: "product1",
            Quantity: 90,
            Price: 280,
        },
        {
            OrderId: "ordsdaahfs",
            Name: "Muzammil3",
            Email: "email@gmail.com",
            Product: "product1",
            Quantity: 180,
            Price: 480,
        },
        {
            OrderId: "osdfsaordahfs",
            Name: "Muzammil4",
            Email: "email@gmail.com",
            Product: "product1",
            Quantity: 90,
            Price: 280,
        },
    }

    comp := adminViews.CustomersIndex("Customers", adminViews.CustomersData(data))
    return comp.Render(c.Request().Context(), c.Response().Writer)
}
