package payment_service

import (
	"errors"
	"fmt"
	"golang-api/internal/dto"
	"golang-api/internal/models"
	"golang-api/internal/repositories"
	"strconv"
	"strings"

	"golang-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AttachItemsService interface {
	AttachMultipleItems(c *fiber.Ctx) error
}

type attachItemsService struct {
	payment        *repositories.PaymentRepository
	paymentItem    *repositories.PaymentItemRepository
	paymentAccount *repositories.PaymentAccountRepository
	item           *repositories.ItemRepository
	generate       *repositories.GenerateRepository
	db             *gorm.DB
}

func NewAttachItemsService(db *gorm.DB) AttachItemsService {
	return &attachItemsService{
		payment:        repositories.NewPaymentRepository(db),
		paymentItem:    repositories.NewPaymentItemRepository(db),
		paymentAccount: repositories.NewPaymentAccountRepository(db),
		item:           repositories.NewItemRepository(db),
		generate:       repositories.NewGenerateRepository(db),
		db:             db,
	}
}

func (s *attachItemsService) AttachMultipleItems(c *fiber.Ctx) error {
	userName := c.Locals("user_name").(string)
	paymentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid payment ID")
	}

	payment, err := s.payment.FindByID(paymentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Payment not found")
	}

	if !payment.HasItems {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "This payment does not support items")
	}

	if payment.TypeID != models.PaymentTypeExpense {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Only expense payments can have items attached")
	}

	var payload dto.AttachMultipleItemsRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if len(payload.Items) == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, "Items are required")
	}

	validationErrs := s.validate(payload)
	if validationErrs != nil {
		return utils.ValidationError(c, validationErrs)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		paymentItems, itemNotes, txErr := s.resolveAndCreateItems(tx, uint(paymentID), payload.Items)
		if txErr != nil {
			return txErr
		}

		if txErr = s.paymentItem.CreateBatch(tx, paymentItems); txErr != nil {
			return errors.New("Failed to attach items")
		}

		oldAmount := int64(0)
		if payment.Amount != nil {
			oldAmount = *payment.Amount
		}
		newAmount := oldAmount + payload.TotalAmount

		existingName := ""
		if payment.Name != nil {
			existingName = *payment.Name
		}
		note := strings.Trim(existingName+", "+strings.Join(itemNotes, ", "), ", ")

		if txErr = s.payment.UpdateFields(tx, uint(paymentID), payment.UserID, userName, map[string]interface{}{
			"amount": newAmount,
			"name":   note,
		}); txErr != nil {
			return errors.New("Failed to update payment")
		}

		return s.updateDeposit(tx, payment, oldAmount, newAmount, userName)
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SimpleSuccessResponse(c, "Items attached successfully")
}

func (s *attachItemsService) validate(payload dto.AttachMultipleItemsRequest) map[string][]string {
	validationErrs := make(map[string][]string)

	for i, item := range payload.Items {
		prefix := "items." + strconv.Itoa(i) + "."
		if strings.TrimSpace(item.Name) == "" {
			validationErrs[prefix+"name"] = []string{"The name field is required"}
		}
		if item.Qty < 1 {
			validationErrs[prefix+"qty"] = []string{"The qty field must be at least 1"}
		}
		if item.Amount < 0 {
			validationErrs[prefix+"amount"] = []string{"The amount field must be a positive number"}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (s *attachItemsService) resolveAndCreateItems(tx *gorm.DB, paymentID uint, items []dto.AttachMultipleItemsItem) ([]models.PaymentItem, []string, error) {
	var paymentItems []models.PaymentItem
	var itemNotes []string

	for _, item := range items {
		var itemID uint

		if item.ItemID != nil {
			itemID = *item.ItemID
		} else {
			existingItem, err := s.item.FindByName(item.Name)
			if err != nil {
				newItem := &models.Item{
					Name:   item.Name,
					Amount: item.Amount,
					TypeID: 1,
					Code:   s.generate.GetCode("item", true),
				}
				if createErr := s.item.CreateWithTx(tx, newItem); createErr != nil {
					return nil, nil, errors.New("Failed to create item")
				}
				itemID = newItem.ID
			} else {
				itemID = existingItem.ID
			}
		}

		total := item.Amount * int64(item.Qty)
		itemCode := s.generate.GetCode("payment_item", true)

		paymentItems = append(paymentItems, models.PaymentItem{
			PaymentID: paymentID,
			ItemID:    itemID,
			ItemCode:  itemCode,
			Quantity:  item.Qty,
			Price:     item.Amount,
			Total:     total,
		})

		itemNotes = append(itemNotes, fmt.Sprintf("%s (x%d)", item.Name, item.Qty))
	}

	return paymentItems, itemNotes, nil
}

func (s *attachItemsService) updateDeposit(tx *gorm.DB, payment *models.Payment, oldAmount int64, newAmount int64, userName string) error {
	paymentAccount, err := s.paymentAccount.SelectByID(tx, payment.PaymentAccountID, []string{"id", "user_id", "name", "deposit"})
	if err != nil {
		return errors.New("Payment account not found")
	}

	deposit := paymentAccount.Deposit
	deposit += oldAmount
	deposit -= newAmount

	_, err = s.paymentAccount.Update(tx, payment.UserID, userName, &models.PaymentAccount{
		ID:      payment.PaymentAccountID,
		Deposit: deposit,
	}, paymentAccount)

	if err != nil {
		return errors.New("Failed to update payment account balance")
	}

	return nil
}
