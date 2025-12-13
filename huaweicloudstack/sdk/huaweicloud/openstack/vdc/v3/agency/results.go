package agency

type CreateAgencyResponse struct {
	Agency Agency `json:"agency"`
}

type AgencyDetail struct {
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

type AgencyRoleResponse struct {
	Total int64        `json:"total"`
	Roles []AgencyRole `json:"roles"`
}

type AgencyRole struct {
	ID            string     `json:"id"`
	DomainID      string     `json:"domain_id"`
	DomainName    string     `json:"domain_name"`
	Name          string     `json:"name"`
	DisplayName   string     `json:"display_name"`
	Flag          string     `json:"flag"`
	Catalog       string     `json:"catalog"`
	Type          string     `json:"type"`
	DescriptionCN string     `json:"description_cn"`
	Description   string     `json:"description"`
	CloudPlatform string     `json:"cloud_platform"`
	Policy        Policy     `json:"policy"`
	Tag           string     `json:"tag"`
	AppName       string     `json:"app_name"`
	DisplayType   string     `json:"display_type"`
	Projects      []AuthItem `json:"projects"`
	Inherit       bool       `json:"inherit"`
	Eps           []AuthItem `json:"eps"`
}

type AuthItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type Policy struct {
	Version string   `json:"Version"`
	Depends []Depend `json:"Depends"`
}

type Depend struct {
	Catalog     string `json:"Catalog"`
	DisplayName string `json:"Display_name"`
}
