package payment_service

import (
	"errors"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"golang-api/pkg/utils"

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
	db             *gorm.DB
}

func NewStoreService(db *gorm.DB) StoreService {
	return &storeService{
		payment:        repositories.NewPaymentRepository(db),
		paymentAccount: repositories.NewPaymentAccountRepository(db),
		generate:       repositories.NewGenerateRepository(db),
		db:             db,
	}
}

func (s *storeService) Store(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	var payload dto.StorePaymentRequest

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

	errs := utils.ValidateJSON(c, &payload, rules)
	if errs != nil {
		return utils.ValidationError(c, errs)
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
		return utils.ValidationError(c, validationErrs)
	}

	if payload.HasItems {
		payload.Amount = nil
		payload.Name = nil
		payload.TypeID = 1
	}

	code := s.generate.GetCode("payment", true)
	var result *models.Payment

	incomeOrExpense := false
	transferOrWithdrawal := false

	switch payload.TypeID {
	case models.PaymentTypeExpense, models.PaymentTypeIncome:
		incomeOrExpense = true
	case models.PaymentTypeTransfer, models.PaymentTypeWithdrawal:
		transferOrWithdrawal = true
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error

		result, err = s.payment.Create(tx, &models.Payment{
			UserID:             userId.(uint),
			Code:               code,
			Name:               payload.Name,
			Amount:             payload.Amount,
			TypeID:             payload.TypeID,
			PaymentAccountID:   payload.PaymentAccountID,
			PaymentAccountToID: payload.PaymentAccountToID,
			HasItems:           payload.HasItems,
			IsScheduled:        payload.IsScheduled,
			IsDraft:            payload.IsDraft,
		})

		if err != nil {
			return errors.New("Failed to create payment, please try again")
		}

		paymentAccount, err := s.paymentAccount.FindByID(payload.PaymentAccountID)
		if err != nil {
			return errors.New("Payment account not found")
		}

		depositChange := paymentAccount.Deposit

		if incomeOrExpense {
			if payload.TypeID == models.PaymentTypeExpense {
				if *payload.Amount > depositChange {
					return errors.New("Insufficient balance for this payment account (e01)")
				}
				depositChange -= *payload.Amount
			} else {
				depositChange += *payload.Amount
			}

			_, err = s.paymentAccount.Update(tx, payload.PaymentAccountID, &models.PaymentAccount{
				Deposit: depositChange,
			})

			if err != nil {
				return errors.New("Failed to update payment account, please try again")
			}
		} else if transferOrWithdrawal {
			paymentAccountTo, err := s.paymentAccount.FindByID(*payload.PaymentAccountToID)
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

			_, err = s.paymentAccount.Update(tx, payload.PaymentAccountID, &models.PaymentAccount{
				Deposit: balanceOrigin,
			})

			if err != nil {
				return errors.New("Failed to update payment account, please try again")
			}

			_, err = s.paymentAccount.Update(tx, *payload.PaymentAccountToID, &models.PaymentAccount{
				Deposit: balanceTo,
			})

			if err != nil {
				return errors.New("Failed to update payment account destination, please try again")
			}
		}

		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, "Payment created successfully", result)
}
