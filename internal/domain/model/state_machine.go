package model

import (
	"errors"
	"strings"
)

var (
	ErrInvalidOrderStateTransition   = errors.New("invalid order state transition")
	ErrInvalidPaymentStateTransition = errors.New("invalid payment state transition")
)

// Order States
const (
	OrderStatePending   = "pending"
	OrderStateProcessing = "processing"
	OrderStateCompleted = "completed"
	OrderStateFailed    = "failed"
	OrderStateCancelled = "cancelled"
	OrderStateShipped   = "shipped"
	OrderStateDelivered = "delivered"
	OrderStateReturned  = "returned"
	OrderStateRefunded  = "refunded"
)

// Payment States
const (
	PaymentStatePending   = "pending"
	PaymentStateSuccess   = "success"
	PaymentStateFailed    = "failed"
	PaymentStateCancelled = "cancelled"
	PaymentStateRefunded  = "refunded"
)

// Order State Transitions
var orderStateTransitions = map[string]map[string]bool{
	OrderStatePending: {
		OrderStateProcessing: true,
		OrderStateCompleted:  true,
		OrderStateFailed:     true,
		OrderStateCancelled:  true,
	},
	OrderStateProcessing: {
		OrderStateCompleted: true,
		OrderStateFailed:    true,
	},
	OrderStateCompleted: {
		OrderStateShipped:  true,
		OrderStateRefunded: true,
	},
	OrderStateFailed: {
		OrderStatePending: true, // Retry
	},
	OrderStateCancelled: {},
	OrderStateShipped: {
		OrderStateDelivered: true,
		OrderStateReturned:  true,
	},
	OrderStateDelivered: {},
	OrderStateReturned:  {},
	OrderStateRefunded:  {},
}

// CanTransitionOrder checks if an order can transition from current to next state
func CanTransitionOrder(current, next string) bool {
	current = strings.ToLower(current)
	next = strings.ToLower(next)
	if current == next {
		return true // No-op
	}
	allowedTransitions, exists := orderStateTransitions[current]
	if !exists {
		return false
	}
	return allowedTransitions[next]
}

// Payment State Transitions
var paymentStateTransitions = map[string]map[string]bool{
	PaymentStatePending: {
		PaymentStateSuccess:   true,
		PaymentStateFailed:    true,
		PaymentStateCancelled: true,
	},
	PaymentStateSuccess: {
		PaymentStateRefunded: true,
	},
	PaymentStateFailed:    {},
	PaymentStateCancelled: {},
	PaymentStateRefunded:  {},
}

// CanTransitionPayment checks if a payment can transition from current to next state
func CanTransitionPayment(current, next string) bool {
	current = strings.ToLower(current)
	next = strings.ToLower(next)
	if current == next {
		return true // No-op
	}
	allowedTransitions, exists := paymentStateTransitions[current]
	if !exists {
		return false
	}
	return allowedTransitions[next]
}
