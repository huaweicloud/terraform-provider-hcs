package agency

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CreateAgencyOpts struct {
	Agency Agency `json:"agency"`
}

type Agency struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DomainID        string `json:"domain_id"`
	TrustDomainID   string `json:"trust_domain_id"`
	TrustDomainName string `json:"trust_domain_name"`
	Description     string `json:"description"`
	Duration        string `json:"duration"`
	ExpireTime      string `json:"expire_time"`
	CreateTime      string `json:"create_time"`
}

func CreateAgency(httpClient *golangsdk.ServiceClient, opts CreateAgencyOpts) (*Agency, error) {
	url := createAgencyURL(httpClient)
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	var r CreateAgencyResponse

	_, err = httpClient.Post(url, b, &r, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return &r.Agency, nil
}

type GetAgencyOpts struct {
	AgencyId   string `q:"agency_id"`
	AgencyName string `q:"agency_name"`
}

func GetAgency(c *golangsdk.ServiceClient, opts GetAgencyOpts) (*AgencyDetail, error) {
	url := getAgencyURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r AgencyDetail

	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type ListAgencyOpts struct {
	Name string `q:"name"`
}

func ListAgency(c *golangsdk.ServiceClient, opts ListAgencyOpts) ([]AgencyDetail, error) {
	url := listAgencyURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()
	var r struct {
		IamAgencies []AgencyDetail `json:"iam_agencies"`
	}

	_, err = c.Get(url, &r, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return r.IamAgencies, nil

}

func GetAgencyRole(c *golangsdk.ServiceClient, id string) ([]AgencyRole, error) {
	url := getAgencyRoleURL(c, id)
	var r AgencyRoleResponse

	_, err := c.Get(url, &r, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return r.Roles, nil
}

func DeleteAgency(c *golangsdk.ServiceClient, id string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult

	_, r.Err = c.Delete(deleteAgencyURL(c, id), &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return &r
}

func CreateAgencyProjectRole(c *golangsdk.ServiceClient, agencyId, projectId, roleId string) error {
	_, err := c.Put(createAgencyProjectRoleURL(c, agencyId, projectId, roleId), nil, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}

func DeleteAgencyProjectRole(c *golangsdk.ServiceClient, agencyId, projectId, roleId string) error {
	_, err := c.Delete(createAgencyProjectRoleURL(c, agencyId, projectId, roleId), &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}

func CreateAgencyDomainRole(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	_, err := c.Put(createAgencyDomainRoleURL(c, agencyId, domainId, roleId), nil, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}

func DeleteAgencyDomainRole(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	_, err := c.Delete(createAgencyDomainRoleURL(c, agencyId, domainId, roleId), &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}

func CreateAgencyDomainInheritedRole(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	_, err := c.Put(createAgencyDomainRoleInheritedURL(c, agencyId, domainId, roleId), nil, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}

func DeleteAgencyDomainInheritedRole(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	_, err := c.Delete(createAgencyDomainRoleInheritedURL(c, agencyId, domainId, roleId), &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201, 204},
		MoreHeaders: golangsdk.MoreHeaders,
	})

	return err
}
