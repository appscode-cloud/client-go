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
	valid, planID := c.GetLicensePlan(clusterID, prodID, prodOwnerID)
	if valid {
		fmt.Printf("Subscribed plan: %s\n", planID)
	} else {
		fmt.Println("Not a Valid license")
	}

	// Output: Not a Valid license
}
