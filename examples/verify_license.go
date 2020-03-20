/*
Copyright 2019 AppsCode Inc.

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

package examples_test

import (
	"fmt"
	"log"

	"go.bytebuilders.dev/client-go"
)

func Example() {
	license := "valid-license-string"
	c := client.NewClient("", license)
	lic, err := c.VerifyLicense()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(lic.SubscribedPlans)

	clusterID, prodID, prodOwnerID := "a-cluster-id", "product-id", int64(2)
	planID, err := c.GetLicensePlan(clusterID, prodID, prodOwnerID)
	if err != nil {
		fmt.Printf("Subscribed plan: %s\n", planID)
	} else {
		fmt.Println(err)
	}

	// Output: Not a Valid license
}
