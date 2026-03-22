package main

import (
	"fmt"

	"github.com/caik/spec2go/pkg/spec"
)

// OrderFailureReason represents reasons an order may be rejected.
type OrderFailureReason string

const (
	OrderTooSmall        OrderFailureReason = "ORDER_TOTAL_TOO_SMALL"
	ItemOutOfStock       OrderFailureReason = "ITEM_OUT_OF_STOCK"
	InvalidShippingAddr  OrderFailureReason = "INVALID_SHIPPING_ADDRESS"
	RestrictedCountry    OrderFailureReason = "SHIPPING_TO_RESTRICTED_COUNTRY"
	PaymentMethodInvalid OrderFailureReason = "PAYMENT_METHOD_INVALID"
)

type OrderItem struct {
	Name    string
	InStock bool
}

type ShippingAddress struct {
	Country string
	ZipCode string
}

type PaymentMethod struct {
	Type  string
	Valid bool
}

type OrderContext struct {
	Total           float64
	Items           []OrderItem
	ShippingAddress ShippingAddress
	Payment         PaymentMethod
}

// --- Custom Specification: checks all items are in stock ---

type allItemsInStockSpec struct {
	spec.NamedSpec[OrderContext, OrderFailureReason]
}

func (s allItemsInStockSpec) Evaluate(ctx OrderContext) spec.SpecificationResult[OrderFailureReason] {
	for _, item := range ctx.Items {
		if !item.InStock {
			// Return fail on first out-of-stock item (one reason per call)
			return spec.Fail(s.Name(), ItemOutOfStock)
		}
	}

	return spec.Pass[OrderFailureReason](s.Name())
}

// --- Shared specifications ---

var (
	minimumOrderValue = spec.New("MinimumOrderValue",
		func(c OrderContext) bool { return c.Total >= 10.0 },
		OrderTooSmall,
	)

	validZipCode = spec.New("ValidZipCode",
		func(c OrderContext) bool { return len(c.ShippingAddress.ZipCode) == 5 },
		InvalidShippingAddr,
	)

	blockedCountry = spec.New("BlockedCountry",
		func(c OrderContext) bool {
			restricted := map[string]bool{"XX": true, "YY": true}
			return restricted[c.ShippingAddress.Country]
		},
		RestrictedCountry,
	)

	validPayment = spec.New("ValidPaymentMethod",
		func(c OrderContext) bool { return c.Payment.Valid },
		PaymentMethodInvalid,
	)

	allItemsInStock = allItemsInStockSpec{
		spec.NamedSpec[OrderContext, OrderFailureReason]{N: "AllItemsInStock"},
	}
)

func main() {
	// NOT: fail if shipping to a blocked country
	notBlockedCountry := spec.Not("NotBlockedCountry", RestrictedCountry, blockedCountry)

	orderPolicy := spec.NewPolicy[OrderContext, OrderFailureReason]().
		With(minimumOrderValue).
		With(allItemsInStock).
		With(validZipCode).
		With(notBlockedCountry).
		With(validPayment)

	fmt.Printf("Policy structure: %s\n\n", orderPolicy)

	orders := []struct {
		name string
		ctx  OrderContext
	}{
		{
			"Valid order",
			OrderContext{
				Total:           99.99,
				Items:           []OrderItem{{"Widget", true}, {"Gadget", true}},
				ShippingAddress: ShippingAddress{"US", "90210"},
				Payment:         PaymentMethod{"credit_card", true},
			},
		},
		{
			"Out-of-stock item + low total",
			OrderContext{
				Total:           5.00,
				Items:           []OrderItem{{"Widget", false}},
				ShippingAddress: ShippingAddress{"US", "90210"},
				Payment:         PaymentMethod{"credit_card", true},
			},
		},
		{
			"Restricted country",
			OrderContext{
				Total:           50.00,
				Items:           []OrderItem{{"Widget", true}},
				ShippingAddress: ShippingAddress{"XX", "12345"},
				Payment:         PaymentMethod{"credit_card", true},
			},
		},
	}

	for _, order := range orders {
		fmt.Printf("Order: %s\n", order.name)
		result := orderPolicy.EvaluateAll(order.ctx)

		if result.AllPassed() {
			fmt.Println("  Status: ACCEPTED")
		} else {
			fmt.Println("  Status: REJECTED")

			for _, failed := range result.FailedResults() {
				fmt.Printf("    - %s: %v\n", failed.Name(), failed.FailureReasons())
			}
		}

		fmt.Println()
	}
}
