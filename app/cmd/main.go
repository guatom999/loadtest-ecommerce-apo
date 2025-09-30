package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/guatom999/ecommerce-product-api/app/databases"
	"github.com/guatom999/ecommerce-product-api/app/handlers"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
	"github.com/guatom999/ecommerce-product-api/app/services"
	"github.com/guatom999/ecommerce-product-api/app/utils"
)

func main() {
	// Load env (use your preferred loader if needed)
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	db := databases.MustOpenPostgres()
	defer db.Close()
	rdb := databases.MustOpenRedis()
	defer rdb.Close()

	userRepo := repositories.NewUserRepo(db)
	productRepo := repositories.NewProductRepo(db)
	jwtMaker := utils.NewJWTMaker()

	authSvc := services.NewAuthService(userRepo, jwtMaker)
	productSvc := services.NewProductService(productRepo)
	cartSvc := services.NewCartService(rdb, productRepo)

	authH := handlers.NewAuthHandler(authSvc)
	prodH := handlers.NewProductHandler(productSvc)
	cartH := handlers.NewCartHandler(cartSvc)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Public
	e.POST("/auth/register", authH.Register)
	e.POST("/auth/login", authH.Login)

	// Protected
	authMW := services.AuthMiddleware(jwtMaker)
	g := e.Group("", authMW)

	// Product
	g.POST("/product/create", prodH.Create)
	g.GET("/product", prodH.List)
	g.GET("/product/:id", prodH.Get)

	// Cart
	g.POST("/cart/add-to-cart", cartH.AddToCart)
	g.GET("/cart", cartH.GetCart)

	addr := ":" + os.Getenv("PORT")
	log.Printf("listening on %s", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
