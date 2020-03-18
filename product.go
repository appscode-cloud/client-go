package client

import (
	"fmt"
	"net/http"

	"kubepack.dev/kubepack/apis/kubepack/v1alpha1"
)

// GetProductByID provides the product associated with the product id
func (c *Client) GetProductByID(productID string) (*v1alpha1.Product, error) {
	apiPth := fmt.Sprintf("products/product_id/%s", productID)

	product := new(v1alpha1.Product)
	if err := c.getParsedResponse(http.MethodPost, apiPth, jsonHeader, nil, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByOwnerAndKey provides the product
// associated with the owner name and product key
func (c *Client) GetProductByOwnerAndKey(owner, key string) (*v1alpha1.Product, error) {
	apiPth := fmt.Sprintf("products/%s/%s", owner, key)

	product := new(v1alpha1.Product)
	if err := c.getParsedResponse(http.MethodPost, apiPth, jsonHeader, nil, product); err != nil {
		return nil, err
	}

	return product, nil
}
