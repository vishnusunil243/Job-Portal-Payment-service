package initializer

import (
	"github.com/vishnusunil243/Job-Portal-Payment-service/concurrency"
	cronjobs "github.com/vishnusunil243/Job-Portal-Payment-service/cron-jobs"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/usecase"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.PaymentEngine {
	adapter := adapters.NewPaymentAdapter(db)
	usecases := usecase.NewPaymentUsecase(adapter)
	services := service.NewPaymentService(usecases)
	c := concurrency.NewConcurrency(db, services)
	c.Concurrency()
	cronjob := cronjobs.NewCronJob(services, db)
	cronjob.Start()
	return service.NewPaymentEngine(services)
}
