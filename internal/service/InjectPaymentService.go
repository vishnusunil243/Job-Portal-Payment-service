package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

type PaymentEngine struct {
	Srv *PaymentService
}

func NewPaymentEngine(srv *PaymentService) *PaymentEngine {
	return &PaymentEngine{
		Srv: srv,
	}
}
func (engine *PaymentEngine) Start(addr string) {
	r := gin.New()
	r.Use(gin.Logger())
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	r.GET("/plans", engine.Srv.getAllSubscriptionPlans)
	r.POST("/subscribe", engine.Srv.addSubscriptionPlan)
	r.GET("/subscribe/payment", engine.Srv.subscriptionPayment)
	r.GET("/payment/verify", engine.Srv.verifyPayment)
	r.GET("/payment/verified", engine.Srv.servePaymentSuccessPage)
	r.LoadHTMLGlob("*.html")
	http.ListenAndServe(addr, handler)
}
