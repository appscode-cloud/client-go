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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.bytebuilders.dev/client/api"
)

// VerifyLicense returns the verified license
func (c *Client) VerifyLicense(licenseData string) (*api.License, error) {
	apiPth := "/user/licenses/verify"
	licenseToken := api.LicenseVerificationParams{
		Raw: licenseData,
	}
	jsonBytes, err := json.Marshal(licenseToken)
	if err != nil {
		return nil, err
	}

	license := new(api.License)
	if err := c.getParsedResponse(http.MethodPost, apiPth, jsonHeader, bytes.NewReader(jsonBytes), license); err != nil {
		return nil, err
	}

	return license, nil
}

// GetLicensePlan provides the plan corresponding to
// product id and owner id if it's still valid
func (c *Client) GetLicensePlan(licenseData, clusterID, productID string, productOwnerID int64) (string, error) {
	license, err := c.VerifyLicense(licenseData)
	if err != nil {
		return "", err
	}

	if license.Audience[0] != clusterID {
		return "", fmt.Errorf("license isn't issued for this cluster")
	} else if license.Status != "active" {
		return "", fmt.Errorf("license status isn't active, current status: %s", license.Status)
	} else if license.NotBefore.Time.Unix() > jwt.NewNumericDate(time.Now()).Time.Unix() {
		return "", fmt.Errorf("license isn't active yet. It will be activated on %v", license.NotBefore.Time.UTC())
	} else if license.Expiry.Time.Unix() < jwt.NewNumericDate(time.Now()).Time.Unix() {
		return "", fmt.Errorf("license expired on: %v", license.Expiry.Time.UTC())
	}

	prod, err := c.GetProductByID(productID)
	if err != nil {
		return "", err
	}

	if prod.Spec.Owner != productOwnerID {
		return "", fmt.Errorf("product doesn't belong to provided ownber")
	}

	plans, err := c.GetProductPlans(productID)
	if err != nil {
		return "", err
	}

	planList := make(map[string]struct{})
	for _, plan := range plans.Items {
		planList[plan.Spec.StripeID] = struct{}{}
	}
	for _, plan := range license.SubscribedPlans {
		if _, ok := planList[plan]; ok {
			return plan, nil
		}
	}
	return "", fmt.Errorf("provided license doesn't include this product")
}
