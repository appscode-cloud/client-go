/*
Copyright AppsCode Inc..

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
	"net/http"
	"time"

	"go.bytebuilders.dev/client/api"

	"gopkg.in/square/go-jose.v2/jwt"
)

// VerifyLicense returns the verified license
func (c *Client) VerifyLicense() (*api.License, error) {
	apiPth := "/user/licenses/verify"
	licenseToken := api.LicenseVerificationParams{
		Raw: c.license,
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
func (c *Client) GetLicensePlan(clusterID, productID string, productOwnerID int64) (bool, string) {
	license, err := c.VerifyLicense()
	if err != nil {
		return false, ""
	}

	if license.Audience[0] != clusterID || license.Status != "active" ||
		license.NotBefore.Time().Unix() > jwt.NewNumericDate(time.Now()).Time().Unix() ||
		license.Expiry.Time().Unix() < jwt.NewNumericDate(time.Now()).Time().Unix() {
		return false, ""
	}
	for _, plans := range license.SubscribedPlans {
		if plans.ProductID == productID && plans.OwnerID == productOwnerID {
			return true, plans.PlanID
		}
	}
	return false, ""
}
