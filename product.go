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

// GetProductByID provides the product associated with the product id
func (c *Client) GetProductByID(productID string) (*v1alpha1.Product, error) {
	apiPth := fmt.Sprintf("products/product_id/%s", productID)

	product := new(v1alpha1.Product)
	if err := c.getParsedResponse(http.MethodGet, apiPth, jsonHeader, nil, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByOwnerAndKey provides the product
// associated with the owner name and product key
func (c *Client) GetProductByOwnerAndKey(owner, key string) (*v1alpha1.Product, error) {
	apiPth := fmt.Sprintf("products/%s/%s", owner, key)

	product := new(v1alpha1.Product)
	if err := c.getParsedResponse(http.MethodGet, apiPth, jsonHeader, nil, product); err != nil {
		return nil, err
	}

	return product, nil
}
