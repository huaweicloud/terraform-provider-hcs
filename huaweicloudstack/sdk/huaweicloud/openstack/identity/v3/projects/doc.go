/*
Package projects manages and retrieves Projects in the OpenStack Identity
Service.

Example to List Projects

	listOpts := projects.ListOpts{
		Enabled: golangsdk.Enabled,
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		panic(err)
	}

	for _, project := range allProjects {
		fmt.Printf("%+v\n", project)
	}
*/
package projects
