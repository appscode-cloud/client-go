package api

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

// LicenseVerificationParams represents the license token for verification
type LicenseVerificationParams struct {
	Raw string `json:"raw"`
}

// License represents the product license for user
type License struct {
	Issuer           string           `json:"issuer,omitempty"`
	Subject          string           `json:"subject,omitempty"`
	Audience         jwt.Audience     `json:"audience,omitempty"`
	Expiry           *jwt.NumericDate `json:"expiry,omitempty"`
	NotBefore        *jwt.NumericDate `json:"not_before,omitempty"`
	IssuedAt         *jwt.NumericDate `json:"issued_at,omitempty"`
	ID               string           `json:"id,omitempty"`
	SubscribedPlans  []SubscribedPlan `json:"subscribed_plans"`
	SubscriptionID   string           `json:"subscription_id"`
	SubscriptionName string           `json:"subscription_name"`
	JWT              string           `json:"jwt"`
	Status           string           `json:"status"`
	CanceledAt       *int64           `json:"canceled_at"`
	IpAddress        *string          `json:"ip_address"`
	CancelerID       *string          `json:"canceler_id"`
}

// SubscribedPlan represents included plans in the license
type SubscribedPlan struct {
	PlanID    string `json:"plan"`
	ProductID string `json:"product"`
	OwnerID   int64  `json:"owner"`
}
