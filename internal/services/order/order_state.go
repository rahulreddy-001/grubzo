package order

//go:generate go run ../../../cmd/injecttrace -file order_state.go -receiver orderStateMachine -service OrderStateMachine

import (
	"grubzo/internal/models/entity"
	"grubzo/internal/router/ext"
)

const (
	orderStatusPending    = "pending"
	orderStatusPreparing  = "preparing"
	orderStatusReady      = "ready"
	orderStatusDelivered  = "delivered"
	orderStatusCancelled  = "cancelled"
	paymentStatusPending  = "pending"
	paymentStatusPaid     = "paid"
	paymentStatusRefunded = "refunded"
	paymentStatusVoided   = "voided"
	paymentModeWallet     = "wallet"
	paymentModePOS        = "pos"
)

type transitionSet map[string]map[string]bool

type orderStateMachine struct {
	orderTransitions   transitionSet
	paymentTransitions map[string]transitionSet
}

func newOrderStateMachine() *orderStateMachine {
	return &orderStateMachine{
		orderTransitions: transitionSet{
			orderStatusPending: {
				orderStatusPreparing: true,
				orderStatusReady:     true,
				orderStatusDelivered: true,
				orderStatusCancelled: true,
			},
			orderStatusPreparing: {
				orderStatusReady:     true,
				orderStatusDelivered: true,
				orderStatusCancelled: true,
			},
			orderStatusReady: {
				orderStatusDelivered: true,
				orderStatusCancelled: true,
			},
			orderStatusDelivered: {},
			orderStatusCancelled: {},
		},
		paymentTransitions: map[string]transitionSet{
			paymentModeWallet: {
				paymentStatusPending: {
					paymentStatusPaid:   true,
					paymentStatusVoided: true,
				},
				paymentStatusPaid: {
					paymentStatusRefunded: true,
				},
				paymentStatusRefunded: {},
				paymentStatusVoided:   {},
			},
			paymentModePOS: {
				paymentStatusPending: {
					paymentStatusPaid: true,
				},
				paymentStatusPaid: {
					paymentStatusRefunded: true,
				},
				paymentStatusRefunded: {},
				paymentStatusVoided:   {},
			},
		},
	}
}

func (sm *orderStateMachine) Validate(order *entity.Order, nextOrderStatus, nextPaymentStatus string) error {
	if nextOrderStatus != "" && !sm.canTransition(sm.orderTransitions, order.Status, nextOrderStatus) {
		return ext.Error("Invalid order state transition")
	}

	if nextPaymentStatus != "" {
		transitions, ok := sm.paymentTransitions[order.PaymentMode]
		if !ok || !sm.canTransition(transitions, order.PaymentStatus, nextPaymentStatus) {
			return ext.Error("Invalid payment state transition")
		}
	}

	return nil
}

func (sm *orderStateMachine) canTransition(transitions transitionSet, current, next string) bool {
	allowedTransitions, ok := transitions[current]
	if !ok {
		return false
	}

	return allowedTransitions[next]
}
