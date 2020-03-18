package client

import (
	"fmt"
	"net/http"

	"kubepack.dev/kubepack/apis/kubepack/v1alpha1"
)

// GetProductPlans provides the product plan list
// associated with the product id
func (c *Client) GetProductPlans(productID string) (v1alpha1.PlanList, error) {
	apiPth := fmt.Sprintf("products/product_id/%s/plans", productID)

	var plans v1alpha1.PlanList
	if err := c.getParsedResponse(http.MethodPost, apiPth, jsonHeader, nil, plans); err != nil {
		return v1alpha1.PlanList{}, err
	}

	return plans, nil
}

// GetProductPlan provides the product plan
// associated with the product id and plan id
func (c *Client) GetProductPlan(productID, planID string) (*v1alpha1.Plan, error) {
	apiPth := fmt.Sprintf("products/product_id/%s/plans/%s", productID, planID)

	var plans = new(v1alpha1.Plan)
	if err := c.getParsedResponse(http.MethodPost, apiPth, jsonHeader, nil, plans); err != nil {
		return nil, err
	}

	return plans, nil
}
