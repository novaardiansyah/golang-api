package payment_service

import (
	"encoding/json"
	"errors"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/gorm"
)

type StoreService interface {
	Store(c *fiber.Ctx) error
}

type storeService struct {
	payment        *repositories.PaymentRepository
	paymentAccount *repositories.PaymentAccountRepository
	generate       *repositories.GenerateRepository
	activityLog    *repositories.ActivityLogRepository
	db             *gorm.DB
}

func NewStoreService(db *gorm.DB) StoreService {
	return &storeService{
		payment:        repositories.NewPaymentRepository(db),
		paymentAccount: repositories.NewPaymentAccountRepository(db),
		generate:       repositories.NewGenerateRepository(db),
		activityLog:    repositories.NewActivityLogRepository(db),
		db:             db,
	}
}

func (s *storeService) Store(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)
	var payload dto.StorePaymentRequest

	validateErrors := s.validate(c, &payload)
	if validateErrors != nil {
		return utils.ValidationError(c, validateErrors)
	}

	s.preparePayload(&payload)
	var result *models.Payment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.createPayment(tx, userId, &payload)
		if err != nil {
			return err
		}

		draft := false
		if payload.IsDraft || payload.IsScheduled {
			draft = true
		}

		if draft == false {
			return s.updateBalances(tx, userId, &payload)
		}

		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	s.saveLog(userId, result)

	return utils.SuccessResponse(c, "Payment created successfully", result)
}

func (s *storeService) preparePayload(payload *dto.StorePaymentRequest) {
	if payload.HasItems {
		payload.Amount = nil
		payload.Name = nil
		payload.TypeID = 1
	}
}

func (s *storeService) createPayment(tx *gorm.DB, userId uint, payload *dto.StorePaymentRequest) (*models.Payment, error) {
	code := s.generate.GetCode("payment", true)
	date, _ := time.Parse("2006-01-02", payload.Date)

	payment, err := s.payment.Create(tx, &models.Payment{
		UserID:             userId,
		Code:               code,
		Name:               payload.Name,
		Date:               date,
		Amount:             payload.Amount,
		TypeID:             payload.TypeID,
		PaymentAccountID:   payload.PaymentAccountID,
		PaymentAccountToID: payload.PaymentAccountToID,
		HasItems:           payload.HasItems,
		IsScheduled:        payload.IsScheduled,
		IsDraft:            payload.IsDraft,
	})

	if err != nil {
		return nil, errors.New("Failed to create payment, please try again")
	}

	return payment, nil
}

func (s *storeService) updateBalances(tx *gorm.DB, userId uint, payload *dto.StorePaymentRequest) error {
	switch payload.TypeID {
	case models.PaymentTypeExpense, models.PaymentTypeIncome:
		return s.handleIncomeOrExpense(tx, userId, payload)
	case models.PaymentTypeTransfer, models.PaymentTypeWithdrawal:
		return s.handleTransferOrWithdrawal(tx, userId, payload)
	}
	return nil
}

func (s *storeService) handleIncomeOrExpense(tx *gorm.DB, userId uint, payload *dto.StorePaymentRequest) error {
	paymentAccount, err := s.paymentAccount.SelectByID(tx, payload.PaymentAccountID, []string{"id", "user_id", "name", "deposit"})

	if err != nil {
		return errors.New("Payment account not found")
	}

	depositChange := paymentAccount.Deposit

	if payload.TypeID == models.PaymentTypeExpense {
		if *payload.Amount > depositChange {
			return errors.New("Insufficient balance for this payment account (e01)")
		}
		depositChange -= *payload.Amount
	} else {
		depositChange += *payload.Amount
	}

	_, err = s.paymentAccount.Update(tx, userId, &models.PaymentAccount{
		ID:      payload.PaymentAccountID,
		Deposit: depositChange,
	}, paymentAccount)

	if err != nil {
		return errors.New("Failed to update payment account, please try again")
	}

	return nil
}

func (s *storeService) handleTransferOrWithdrawal(tx *gorm.DB, userId uint, payload *dto.StorePaymentRequest) error {
	paymentAccount, err := s.paymentAccount.SelectByID(tx, payload.PaymentAccountID, []string{"id", "user_id", "name", "deposit"})
	if err != nil {
		return errors.New("Payment account not found")
	}

	paymentAccountTo, err := s.paymentAccount.SelectByID(tx, *payload.PaymentAccountToID, []string{"id", "user_id", "name", "deposit"})
	if err != nil {
		return errors.New("Payment account destination not found")
	}

	balanceOrigin := paymentAccount.Deposit
	balanceTo := paymentAccountTo.Deposit

	if balanceOrigin < *payload.Amount {
		return errors.New("Insufficient balance for this payment account (e02)")
	}

	balanceOrigin -= *payload.Amount
	balanceTo += *payload.Amount

	_, err = s.paymentAccount.Update(tx, userId, &models.PaymentAccount{
		ID:      payload.PaymentAccountID,
		Deposit: balanceOrigin,
	}, paymentAccount)

	if err != nil {
		return errors.New("Failed to update payment account, please try again")
	}

	_, err = s.paymentAccount.Update(tx, userId, &models.PaymentAccount{
		ID:      *payload.PaymentAccountToID,
		Deposit: balanceTo,
	}, paymentAccountTo)

	if err != nil {
		return errors.New("Failed to update payment account destination, please try again")
	}

	return nil
}

func (s *storeService) validate(c *fiber.Ctx, payload *dto.StorePaymentRequest) map[string][]string {
	rules := govalidator.MapData{
		"amount":                []string{"numeric"},
		"date":                  []string{"required", "date:yyyy-mm-dd"},
		"name":                  []string{"max:255"},
		"type_id":               []string{"required", "numeric"},
		"payment_account_id":    []string{"required", "numeric"},
		"payment_account_to_id": []string{"numeric"},
		"has_items":             []string{"bool"},
		"is_scheduled":          []string{"bool"},
		"is_draft":              []string{"bool"},
		"request_view":          []string{"bool"},
	}

	errs := utils.ValidateJSON(c, payload, rules)
	if errs != nil {
		return errs
	}

	validationErrs := make(map[string][]string)

	if payload.HasItems == false {
		if payload.Amount == nil || *payload.Amount < 1 {
			validationErrs["amount"] = []string{"This field is required when the payment has no items", "This field must be greater than 0"}
		}

		if payload.Name == nil || *payload.Name == "" {
			validationErrs["name"] = []string{"This field is required when the payment has no items"}
		}
	}

	if payload.TypeID == 3 || payload.TypeID == 4 {
		if payload.PaymentAccountToID == nil {
			validationErrs["payment_account_to_id"] = []string{"This field is required when the category is transfer or widrawal."}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (s *storeService) saveLog(userId uint, result *models.Payment) {
	logProps := dto.PaymentLogProperties{
		ID:                 result.ID,
		UserID:             result.UserID,
		Code:               result.Code,
		Name:               result.Name,
		Date:               result.Date,
		Amount:             result.Amount,
		HasItems:           result.HasItems,
		IsScheduled:        result.IsScheduled,
		IsDraft:            result.IsDraft,
		Attachments:        result.Attachments,
		TypeID:             result.TypeID,
		PaymentAccountID:   result.PaymentAccountID,
		PaymentAccountToID: result.PaymentAccountToID,
	}

	properties, _ := json.Marshal(logProps)

	err := s.activityLog.Store(&models.ActivityLog{
		Event:       "Created",
		LogName:     "Resource",
		Description: "Payment Created by Nova Ardiansyah (Hardcode)",
		SubjectType: utils.String("App\\Models\\Payment"),
		SubjectID:   &result.ID,
		CauserType:  "App\\Models\\User",
		CauserID:    userId,
		Properties:  properties,
	})

	if err != nil {
		log.Println("Transaction successfully saved, but failed to save activity log", err)
	}
}
