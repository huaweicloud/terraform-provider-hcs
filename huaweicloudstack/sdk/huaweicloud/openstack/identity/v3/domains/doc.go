/*
Package domains manages and retrieves Domains in the OpenStack Identity Service.

Example to List Domains

	var iTrue bool = true
	listOpts := domains.ListOpts{
		Enabled: &iTrue,
	}

	allPages, err := domains.List(identityClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		panic(err)
	}

	for _, domain := range allDomains {
		fmt.Printf("%+v\n", domain)
	}
*/
package domains
