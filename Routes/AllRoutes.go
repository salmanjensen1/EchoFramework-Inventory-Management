package Routes

import (
	"RiseOfProduceManagement/Auth"
	"RiseOfProduceManagement/Controller"
	"RiseOfProduceManagement/Middleware"
	"RiseOfProduceManagement/configs"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func NormalRoutes(e *echo.Echo) {
	e.GET("/get-product/:productID", Controller.GetProduct)
	e.GET("/search-product/:productName", Controller.SearchProduct)
}

func AuthenticationRoutes(e *echo.Echo) {
	e.POST("/login", Auth.Login)
	e.POST("/register", Auth.Register)
}

func AdminRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.Use(echojwt.WithConfig(configs.Config))
	//admin routes
	a := auth.Group("/forAdmin", Middleware.ValidateToken, Middleware.IsAdmin)
	a.GET("/", configs.Restricted)
	a.POST("/make-admin", Controller.MakeAdmin)
}

func UserRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.Use(echojwt.WithConfig(configs.Config))
	//user routes
	u := auth.Group("/forUser", Middleware.ValidateToken)
	u.GET("/", configs.Restricted)
	u.PUT("/update-user/:userID", Controller.UpdateUser)
	u.DELETE("/delete-user/:userID", Controller.DeleteUser)
	u.GET("/buy-product", Controller.BuyProduct)
	u.GET("/add-money/:money", Controller.AddMoneyToAccount)
}

func ProductRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.Use(echojwt.WithConfig(configs.Config))
	//user routes
	r := auth.Group("/forUser", Middleware.ValidateToken)
	r.GET("/", configs.Restricted)
	r.POST("/create-product", Controller.CreateProduct)
	r.GET("/get-all-product/:sellerID", Controller.GetAllProductsOfASeller)
	r.PUT("/update-product/:productID", Controller.UpdateProduct)
	r.DELETE("/delete-product/:productID", Controller.DeleteProduct)
}
