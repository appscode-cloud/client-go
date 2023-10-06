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
	"fmt"
	"net/http"

	"x-helm.dev/apimachinery/apis/products/v1alpha1"
)

// GetProductPlans provides the product plan list
// associated with the product id
func (c *Client) GetProductPlans(productID string) (v1alpha1.PlanList, error) {
	apiPth := fmt.Sprintf("products/product_id/%s/plans", productID)

	var plans v1alpha1.PlanList
	if err := c.getParsedResponse(http.MethodGet, apiPth, jsonHeader, nil, &plans); err != nil {
		return v1alpha1.PlanList{}, err
	}

	return plans, nil
}

// GetProductPlan provides the product plan
// associated with the product id and plan id
func (c *Client) GetProductPlan(productID, planID string) (*v1alpha1.Plan, error) {
	apiPth := fmt.Sprintf("products/product_id/%s/plans/plan_id/%s", productID, planID)

	plans := new(v1alpha1.Plan)
	if err := c.getParsedResponse(http.MethodGet, apiPth, jsonHeader, nil, plans); err != nil {
		return nil, err
	}

	return plans, nil
}
