package mpesa

import (
	"fmt"
	"strconv"
	"time"
)

func ParseCallback(callback *MpesaCallback) (ParsedCallback, error) {
	parsed := ParsedCallback{
		CheckoutRequestID: callback.Body.StkCallback.CheckoutRequestID,
		ResultCode: 	   callback.Body.StkCallback.ResultCode,
		ResultDesc: 	   callback.Body.StkCallback.ResultDesc,
	}

	if callback.Body.StkCallback.ResultCode != 0 {
		return parsed, nil
	}

	for _, item := range callback.Body.StkCallback.CallbackMetadata.Item {
		switch item.Name {
		case "Amount":
			if val, ok := item.Value.(float64); ok {
				parsed.Amount = val
			} else if val, ok := item.Value.(string); ok {
				if amount, err := strconv.ParseFloat(val, 64); err == nil {
					parsed.Amount = amount
				}
			}
		case "MpesaReceiptNumber":
			if val, ok := item.Value.(string); ok {
				parsed.MpesaReceiptNumber = val
			}	
		case "TransactionDate":
			if val, ok := item.Value.(float64); ok {
				dateStr := fmt.Sprintf("%.0f", val)
				if t, err := time.Parse("20060102150405", dateStr); err == nil {
					parsed.TransactionDate = t
				}
			} else if val, ok := item.Value.(string); ok {
				if t, err := time.Parse("20060102150405", val); err == nil {
					parsed.TransactionDate = t
				}
			}
		case "PhoneNumber":
			if val, ok := item.Value.(float64); ok {
				parsed.PhoneNumber = fmt.Sprintf("%.0f", val)
			} else if val, ok := item.Value.(string); ok {
				parsed.PhoneNumber = val
			}	
		}
	}

	return parsed, nil
}

func ValidateCallback(callback *MpesaCallback) error {
	if callback.Body.StkCallback.CheckoutRequestID == "" {
		return fmt.Errorf("missing checkoutRequestID")
	}

	if callback.Body.StkCallback.ResultCode == 0 {
		if len(callback.Body.StkCallback.CallbackMetadata.Item) == 0 {
			return fmt.Errorf("missing callback metadata for successful payment")
		}
	}
	return nil
}
