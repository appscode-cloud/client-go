/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"time"

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
	SubscribedPlans  []string         `json:"subscribed_plans"`
	SubscriptionID   string           `json:"subscription_id"`
	SubscriptionName string           `json:"subscription_name"`
	JWT              string           `json:"jwt"`
	Status           string           `json:"status"`
	CanceledAt       *int64           `json:"canceled_at"`
	IpAddress        *string          `json:"ip_address"`
	CancelerID       *string          `json:"canceler_id"`
}

type User struct {
	// the user's id
	ID int64 `json:"id"`
	// the user's username
	UserName string `json:"login"`
	// the user's full name
	FullName string `json:"full_name"`
	// swagger:strfmt email
	Email string `json:"email"`
	// URL to the user's avatar
	AvatarURL string `json:"avatar_url"`
	// User locale
	Language string `json:"language"`
	// Is the user an administrator
	IsAdmin bool `json:"is_admin"`
	// swagger:strfmt date-time
	LastLogin time.Time `json:"last_login,omitempty"`
	// swagger:strfmt date-time
	Created time.Time `json:"created,omitempty"`
	// define individual user or organization
	Type int `json:"type"`
	// Is user active
	IsActive bool `json:"active"`
	// Is user login prohibited
	ProhibitLogin bool `json:"prohibit_login"`
	// the user's location
	Location string `json:"location"`
	// the user's website
	Website string `json:"website"`
	// the user's description
	Description string `json:"description"`
}
