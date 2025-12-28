package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
)

type FieldRule struct {
	Type     string
	Required bool
	Min      int
	Max      int
}

func ValidateJSONMap(c *fiber.Ctx, fieldRules map[string]FieldRule) (map[string]interface{}, map[string][]string) {
	data := make(map[string]interface{})
	errors := make(map[string][]string)

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		errors["_error"] = []string{"Invalid JSON body"}
		return nil, errors
	}

	for field, rule := range fieldRules {
		value, exists := data[field]
		display := strings.ReplaceAll(field, "_", " ")

		if rule.Required && (!exists || value == nil) {
			errors[field] = append(errors[field], fmt.Sprintf("The %s field is required.", display))
			continue
		}

		if !exists || value == nil {
			continue
		}

		switch rule.Type {
      case "bool", "boolean":
        if _, ok := value.(bool); !ok {
          errors[field] = append(errors[field], fmt.Sprintf("The %s field must be true or false.", display))
        }
      case "string":
        if s, ok := value.(string); ok {
          if rule.Min > 0 && len(s) < rule.Min {
            errors[field] = append(errors[field], fmt.Sprintf("The %s field must be at least %d characters.", display, rule.Min))
          }
          if rule.Max > 0 && len(s) > rule.Max {
            errors[field] = append(errors[field], fmt.Sprintf("The %s field must not be greater than %d characters.", display, rule.Max))
          }
        } else {
          errors[field] = append(errors[field], fmt.Sprintf("The %s field must be a string.", display))
        }
      case "numeric", "number":
        switch value.(type) {
        case float64, float32, int, int64:
        default:
          errors[field] = append(errors[field], fmt.Sprintf("The %s field must be a number.", display))
        }
      case "email":
        if s, ok := value.(string); ok {
          if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
            errors[field] = append(errors[field], fmt.Sprintf("The %s field must be a valid email address.", display))
          }
        } else {
          errors[field] = append(errors[field], fmt.Sprintf("The %s field must be a string.", display))
        }
      }
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return data, nil
}

func ValidateJSONStruct(c *fiber.Ctx, data interface{}, rules govalidator.MapData) map[string][]string {
	body := c.Body()
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")

	opts := govalidator.Options{
		Request: req,
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	errs := v.ValidateJSON()

	if len(errs) > 0 {
		errors := make(map[string][]string)
		for field, msgs := range errs {
			errors[field] = msgs
		}
		return errors
	}

	return nil
}

func ValidateJSONStructWithMessages(c *fiber.Ctx, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	body := c.Body()
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")

	opts := govalidator.Options{
		Request:  req,
		Data:     data,
		Rules:    rules,
		Messages: messages,
	}

	v := govalidator.New(opts)
	errs := v.ValidateJSON()

	if len(errs) > 0 {
		errors := make(map[string][]string)
		for field, msgs := range errs {
			errors[field] = msgs
		}
		return errors
	}

	return nil
}

func GetString(data map[string]interface{}, key string) *string {
	if val, ok := data[key]; ok && val != nil {
		if s, ok := val.(string); ok {
			return &s
		}
	}
	return nil
}

func GetBool(data map[string]interface{}, key string) *bool {
	if val, ok := data[key]; ok && val != nil {
		if b, ok := val.(bool); ok {
			return &b
		}
	}
	return nil
}
