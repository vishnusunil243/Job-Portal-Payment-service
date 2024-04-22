package service

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"github.com/vishnusunil243/Job-Portal-Payment-service/entities"
	"github.com/vishnusunil243/Job-Portal-Payment-service/helperstruct"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/usecase"
	"github.com/vishnusunil243/Job-Portal-Payment-service/kafka"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/grpc"
)

type PaymentService struct {
	usecases usecase.PaymentUsecaseInterface
	UserConn pb.UserServiceClient
}

func NewPaymentService(usecases usecase.PaymentUsecaseInterface) *PaymentService {
	userConn, _ := grpc.Dial("user-service:8081", grpc.WithInsecure())
	return &PaymentService{
		usecases: usecases,
		UserConn: pb.NewUserServiceClient(userConn),
	}
}

func (p *PaymentService) subscriptionPayment(c *gin.Context) {
	subId := c.Query("plan_id")
	userId := c.Query("user_id")
	planInfo, err := p.usecases.GetSubscriptionPlanById(subId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error retrieving subscription plans",
			Error:      err.Error(),
		})
		return
	}
	usersubsciption, err := p.usecases.GetUserSubscription(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error getting user info",
			Error:      err.Error(),
		})
	}
	if !usersubsciption.SubscribedTill.IsZero() {
		durationExpiry := time.Until(usersubsciption.SubscribedTill)
		if durationExpiry > 7*24*time.Hour {
			c.HTML(200, "alreadyPaid.html", gin.H{})
			return
		}
	}
	client := razorpay.NewClient(os.Getenv("RAZORPAYID"), os.Getenv("RAZORPAYSECRET"))
	data := map[string]interface{}{
		"amount":   planInfo.Amount * 100,
		"currency": "INR",
		"receipt":  "test_receipt_id",
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error creating razorpay client",
			Error:      err.Error(),
		})
		return
	}
	value := body["id"]
	razorPayId := value.(string)
	c.HTML(200, "payment.html", gin.H{
		"total_amount": planInfo.Amount,
		"total":        planInfo.Amount,
		"orderid":      razorPayId,
		"plan_id":      subId,
		"userId":       userId,
	})
}
func (p *PaymentService) verifyPayment(c *gin.Context) {
	paymentRef := c.Query("payment_ref")
	userId := c.Query("user_id")
	subId := c.Query("plan_id")
	userID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "please chose a valid user to proceed",
			Error:      err.Error(),
		})
	}
	// idStr := c.Query("plan_id")
	// planId := strings.ReplaceAll(idStr, " ", "")
	if paymentRef == "" {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "payment failed",
		})
		return
	}
	if err := p.usecases.AddPayment(entities.Payment{
		UserId:     userID,
		PaymentRef: paymentRef,
	}); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "error while adding payment",
			Error:      err.Error(),
		})
		return
	}
	plan, err := p.usecases.GetSubscriptionPlanById(subId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while retrieving plan info",
			Error:      err.Error(),
		})
	}
	if err := p.usecases.AddUserSubscription(userId, subId, plan.Duration); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while adding subscription to user",
		})
		return
	}
	if _, err := p.UserConn.UpdateSubscription(context.Background(), &pb.UpdateSubscriptionRequest{
		UserId:       userId,
		Subscription: true,
	}); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error updating subscription status",
			Error:      err.Error(),
		})
		return
	}
	userData, err := p.UserConn.GetUser(context.Background(), &pb.GetUserById{
		Id: userId,
	})
	if err != nil {
		log.Print("error retrieving user info")
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while retrieving user info",
			Error:      err.Error(),
		})
		return
	}
	if err := kafka.SubscribedMessage(userData.Email, plan.Duration); err != nil {
		log.Print("error while sending message", err)
	}
	c.JSON(http.StatusOK, helperstruct.Response{
		StatusCode: 200,
		Message:    "payment verified",
		Data:       true,
	})

}
func (p *PaymentService) servePaymentSuccessPage(c *gin.Context) {
	c.HTML(200, "paymentVerified.html", gin.H{})
}
func (p *PaymentService) addSubscriptionPlan(c *gin.Context) {
	var req entities.Subscription
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "please provide valid subscriprion plan",
			Error:      err.Error(),
		})
		return
	}

	if !helper.IsValidDuration(req.Duration) {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "invalid duration",
			Error:      "please enter a valid duration",
		})
		return
	}
	check, err := p.usecases.GetSubscriptionByDuration(req.Duration)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error adding subscription plan",
			Error:      err.Error(),
		})
		return
	}
	if check.Duration != "" {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "plan for this duration is already present please add a new plan",
			Error:      "plan for this duration is already present please add a new plan",
		})
		return
	}
	if helper.IsNegative(req.Amount) {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "invalid amount",
			Error:      "please enter a valid amount",
		})
		return
	}
	if err := p.usecases.AddSubscriptionPlan(req); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while adding subscription plan",
			Error:      err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, helperstruct.Response{
		StatusCode: 200,
		Message:    "plan added successfully",
	})
}
func (p *PaymentService) updateSubscriptionPlan(c *gin.Context) {
	var req entities.Subscription
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &helperstruct.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Error:      err.Error(),
		})
		return
	}
	if !helper.IsValidDuration(req.Duration) {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "invalid duration",
			Error:      "please enter a valid duration",
		})
		return
	}
	if helper.IsNegative(req.Amount) {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "invalid amount",
			Error:      "please enter a valid amount",
		})
		return
	}
	subId := c.Query("sub_id")
	sId, err := uuid.Parse(subId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error parsing subscription id",
			Error:      err.Error(),
		})
	}
	req.Id = sId
	if err := p.usecases.UpdateSubscriptionPlan(req); err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error updating subscription plan",
			Error:      err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, helperstruct.Response{
		StatusCode: 200,
		Message:    "subscription plan updated successfully",
	})
}
func (p *PaymentService) getAllSubscriptionPlans(c *gin.Context) {
	plans, err := p.usecases.GetAllSubscriptionPlans()
	if err != nil {
		c.JSON(http.StatusBadRequest, helperstruct.Response{
			StatusCode: 400,
			Message:    "error while fetching plans",
			Error:      err.Error(),
		})
	}
	c.JSON(http.StatusOK, helperstruct.Response{
		StatusCode: 200,
		Message:    "successfully fetched subscriptions",
		Data:       plans,
	})
}
