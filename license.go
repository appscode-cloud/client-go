package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go.bytebuilder.dev/client/api"

	"gopkg.in/square/go-jose.v2/jwt"
)

// GetLicensePlan provides the plan corresponding to
// product id and owner id if it's still valid
func GetLicensePlan(token, clusterID, productID string, productOwnerID int64) (bool, string) {
	license, err := VerifyLicense(token)
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

// VerifyLicense returns the verified license
func VerifyLicense(token string) (*api.License, error) {
	url := "https://byte.builders/api/v1/user/licenses/verify"
	licenseToken := api.LicenseVerificationParams{
		Raw: token,
	}
	jsonBytes, err := json.Marshal(licenseToken)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "JWT "+token)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	license := new(api.License)
	if err = json.NewDecoder(resp.Body).Decode(license); err != nil {
		return nil, err
	}

	return license, nil
}
