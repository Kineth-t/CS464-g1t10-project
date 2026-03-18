package main

import (
	"log"
	"net/http"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
	"github.com/Kineth-t/CS464-g1t10-project/internal/router"
	"github.com/Kineth-t/CS464-g1t10-project/internal/handler"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

func main() {
	// Repos
	phoneRepo := repository.NewPhoneRepository()
	userRepo  := repository.NewUserRepository()
	cartRepo  := repository.NewCartRepository()

	// Services
	phoneSvc := service.NewPhoneService(phoneRepo)
	authSvc  := service.NewAuthService(userRepo)
	cartSvc  := service.NewCartService(cartRepo, phoneRepo)

	// Handlers
	ph := handler.NewPhoneHandler(phoneSvc)
	ah := handler.NewAuthHandler(authSvc)
	ch := handler.NewCartHandler(cartSvc)

	r := router.Setup(ph, ah, ch)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}