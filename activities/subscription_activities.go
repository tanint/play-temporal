package activities

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// SubscriptionDetails contains information about a subscription
type SubscriptionDetails struct {
	ID              string
	CustomerID      string
	PlanID          string
	PricePerMonth   float64
	StartDate       time.Time
	BillingDay      int
	Status          string
	PaymentMethodID string
}

// InvoiceDetails contains information about an invoice
type InvoiceDetails struct {
	ID             string
	SubscriptionID string
	Amount         float64
	Currency       string
	Status         string
	DueDate        time.Time
	Items          []InvoiceItem
}

// InvoiceItem represents a line item in an invoice
type InvoiceItem struct {
	Description string
	Amount      float64
	Quantity    int
}

// PaymentDetails contains information about a payment
type PaymentDetails struct {
	ID              string
	InvoiceID       string
	Amount          float64
	Currency        string
	Status          string
	PaymentMethodID string
	ProcessedAt     time.Time
}

// CreateSubscriptionActivity simulates creating a new subscription
func CreateSubscriptionActivity(ctx context.Context, customerID string, planID string) (SubscriptionDetails, error) {
	fmt.Printf("[Subscription Activity] Creating subscription for customer %s on plan %s\n", customerID, planID)

	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	// Generate a random subscription ID
	subscriptionID := fmt.Sprintf("sub_%d", rand.Intn(1000000))

	// Create subscription details
	subscription := SubscriptionDetails{
		ID:              subscriptionID,
		CustomerID:      customerID,
		PlanID:          planID,
		PricePerMonth:   rand.Float64() * 100, // Random price between 0 and 100
		StartDate:       time.Now(),
		BillingDay:      time.Now().Day(),
		Status:          "active",
		PaymentMethodID: fmt.Sprintf("pm_%d", rand.Intn(1000000)),
	}

	fmt.Printf("[Subscription Activity] Created subscription %s with monthly price %.2f\n",
		subscription.ID, subscription.PricePerMonth)

	return subscription, nil
}

// CalculateChargesActivity simulates calculating charges for a billing period
func CalculateChargesActivity(ctx context.Context, subscription SubscriptionDetails) (float64, error) {
	fmt.Printf("[Subscription Activity] Calculating charges for subscription %s\n", subscription.ID)

	// Simulate processing time
	time.Sleep(300 * time.Millisecond)

	// Base charge is the subscription price
	baseCharge := subscription.PricePerMonth

	// Simulate some usage-based charges (random)
	usageCharge := rand.Float64() * 20 // Random usage between 0 and 20

	// Calculate total
	totalCharge := baseCharge + usageCharge

	fmt.Printf("[Subscription Activity] Calculated charges for subscription %s: base=%.2f, usage=%.2f, total=%.2f\n",
		subscription.ID, baseCharge, usageCharge, totalCharge)

	return totalCharge, nil
}

// GenerateInvoiceActivity simulates generating an invoice
func GenerateInvoiceActivity(ctx context.Context, subscription SubscriptionDetails, amount float64) (InvoiceDetails, error) {
	fmt.Printf("[Subscription Activity] Generating invoice for subscription %s\n", subscription.ID)

	// Simulate processing time
	time.Sleep(400 * time.Millisecond)

	// Generate a random invoice ID
	invoiceID := fmt.Sprintf("inv_%d", rand.Intn(1000000))

	// Create invoice details
	invoice := InvoiceDetails{
		ID:             invoiceID,
		SubscriptionID: subscription.ID,
		Amount:         amount,
		Currency:       "USD",
		Status:         "pending",
		DueDate:        time.Now().Add(7 * 24 * time.Hour), // Due in 7 days
		Items: []InvoiceItem{
			{
				Description: fmt.Sprintf("Subscription to %s", subscription.PlanID),
				Amount:      subscription.PricePerMonth,
				Quantity:    1,
			},
			{
				Description: "Usage charges",
				Amount:      amount - subscription.PricePerMonth,
				Quantity:    1,
			},
		},
	}

	fmt.Printf("[Subscription Activity] Generated invoice %s for subscription %s with amount %.2f USD\n",
		invoice.ID, subscription.ID, invoice.Amount)

	return invoice, nil
}

// ProcessPaymentActivity simulates processing a payment for an invoice
func ProcessPaymentActivity(ctx context.Context, invoice InvoiceDetails, subscription SubscriptionDetails) (PaymentDetails, error) {
	fmt.Printf("[Subscription Activity] Processing payment for invoice %s\n", invoice.ID)

	// Simulate processing time
	time.Sleep(600 * time.Millisecond)

	// Generate a random payment ID
	paymentID := fmt.Sprintf("py_%d", rand.Intn(1000000))

	// Simulate payment success (90% chance)
	paymentStatus := "succeeded"
	if rand.Float64() < 0.1 {
		paymentStatus = "failed"
	}

	// Create payment details
	payment := PaymentDetails{
		ID:              paymentID,
		InvoiceID:       invoice.ID,
		Amount:          invoice.Amount,
		Currency:        invoice.Currency,
		Status:          paymentStatus,
		PaymentMethodID: subscription.PaymentMethodID,
		ProcessedAt:     time.Now(),
	}

	fmt.Printf("[Subscription Activity] Processed payment %s for invoice %s with status: %s\n",
		payment.ID, invoice.ID, payment.Status)

	return payment, nil
}

// SendInvoiceEmailActivity simulates sending an invoice email
func SendInvoiceEmailActivity(ctx context.Context, invoice InvoiceDetails, customerID string) error {
	fmt.Printf("[Subscription Activity] Sending invoice email for invoice %s to customer %s\n",
		invoice.ID, customerID)

	// Simulate processing time
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("[Subscription Activity] Invoice email sent successfully for invoice %s\n", invoice.ID)

	return nil
}

// UpdateSubscriptionStatusActivity simulates updating a subscription status
func UpdateSubscriptionStatusActivity(ctx context.Context, subscriptionID string, status string) error {
	fmt.Printf("[Subscription Activity] Updating subscription %s status to: %s\n",
		subscriptionID, status)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("[Subscription Activity] Updated subscription %s status to: %s\n",
		subscriptionID, status)

	return nil
}
