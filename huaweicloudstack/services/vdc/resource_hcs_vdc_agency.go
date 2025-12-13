package vdc

import (
	"context"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/identity/v3/projects"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/agency"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

// ResourceVdcAgency
// @API VDC POST /rest/vdc/v3.0/vdc-agencies
// @API VDC GET /rest/vdc/v3.0/tenant-agencies/agency-detail?agency_id={agency_id}
// @API VDC GET /rest/vdc/v3.0/vdc-agencies/{agency_id}/roles
// @API VDC GET /rest/vdc/v3.1/vdcs/{vdc_id}/projects
// @API VDC GET /rest/vdc/v3.0/OS-ROLE/roles/third-party/roles
// @API VDC PUT /rest/vdc/v3.0/vdc-agencies/{agency_id}/projects/{project_id}/roles/{role_id}
// @API VDC PUT /rest/vdc/v3.0/vdc-agencies/{agency_id}/domains/{domain_id}/roles/{role_id}/inherited_to_projects
// @API VDC PUT /rest/vdc/v3.0/vdc-agencies/{agency_id}/domains/{domain_id}/roles/{role_id}
// @API VDC DELETE /rest/vdc/v3.0/vdc-agencies/{agency_id}/domains/{domain_id}/roles/{role_id}/inherited_to_projects
// @API VDC DELETE /rest/vdc/v3.0/vdc-agencies/{agency_id}/projects/{project_id}/roles/{role_id}
// @API VDC DELETE /rest/vdc/v3.0/vdc-agencies/{agency_id}/domains/{domain_id}/roles/{role_id}
// @API VDC DELETE /rest/vdc/v3.0/tenant-agencies/{agency_id}
func ResourceVdcAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcAgencyCreate,
		ReadContext:   resourceVdcAgencyRead,
		UpdateContext: resourceVdcAgencyUpdate,
		DeleteContext: resourceVdcAgencyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				ForceNew: true,
			},
			"delegated_domain_name": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				ForceNew: true,
			},
			"project_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
						"roles": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
							Optional: false,
						},
					},
				},
			},
			"domain_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Required: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_resources_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Required: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVdcAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	desc := d.Get("description").(string)
	domainName := d.Get("delegated_domain_name").(string)
	cfg := meta.(*config.Config)
	hcsConfig := config.GetHcsConfig(meta)

	opts := agency.CreateAgencyOpts{
		Agency: agency.Agency{
			Name:            name,
			Description:     desc,
			DomainID:        cfg.DomainID,
			TrustDomainName: domainName,
			Duration:        "FOREVER",
		},
	}
	log.Printf("[DEBUG] The createOpt of ServiceStage environment is: %v", opts)

	vdcClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating VDC client : %s", err)
	}

	// create agency first
	r, err := agency.CreateAgency(vdcClient, opts)
	if err != nil {
		return fmtp.DiagErrorf("error creating agency : %s", err)
	}
	d.SetId(r.ID)

	// add agency roles
	if err = updateAgencyRole(cfg, hcsConfig, vdcClient, nil, d); err != nil {
		return fmtp.DiagErrorf("error creating agency authrization : %s", err)
	}

	return resourceVdcAgencyRead(ctx, d, meta)
}

func resourceVdcAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	hcsConfig := config.GetHcsConfig(meta)
	vdcClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating VDC client : %s", err)
	}

	// get agency information
	ag, err := agency.GetAgency(vdcClient, agency.GetAgencyOpts{
		AgencyId: d.Id(),
	})
	if err != nil {
		return fmtp.DiagErrorf("error querying agency : %s", err)
	}

	d.SetId(ag.ID)

	// get agency roles
	roles, err := agency.GetAgencyRole(vdcClient, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("error querying agency role : %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", ag.Name),
		d.Set("description", ag.Description),
		d.Set("delegated_domain_name", ag.TrustDomainName),
		setAgencyRole(cfg, d, roles),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVdcAgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	hcsConfig := config.GetHcsConfig(meta)
	vdcClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))

	roles, err := agency.GetAgencyRole(vdcClient, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("error querying agency role : %s", err)
	}

	if err = updateAgencyRole(cfg, hcsConfig, vdcClient, roles, d); err != nil {
		return fmtp.DiagErrorf("error creating agency authrization : %s", err)
	}

	return resourceVdcAgencyRead(ctx, d, meta)
}

func resourceVdcAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating VDC client : %s", err)
	}

	err = agency.DeleteAgency(vdcClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting agency %s: %s", d.Id(), err)
	}

	return nil
}

func updateAgencyRole(cfg *config.Config, hcsConfig *config.HcsConfig, c *golangsdk.ServiceClient, existingRoles []agency.AgencyRole, newRoles *schema.ResourceData) error {
	epr, edr, edri := buildAgencyRoleFromResponse(cfg.DomainID, existingRoles)
	pr, dr, dri := buildAgencyRoleFromSchema(cfg, newRoles)

	pjm, err := getProjectMap(cfg, hcsConfig, projects.ListOpts{DomainID: cfg.DomainID})
	if err != nil {
		return fmtp.Errorf("error listing projects: %s", err)
	}
	rom, err := getRoleMap(c, role.ListOpts{DomainId: cfg.DomainID, Limit: 100})
	if err != nil {
		return fmtp.Errorf("error listing roles: %s", err)
	}

	// set project name and role name
	for i, e := range pr {
		v, ok := pjm[e.itemName]
		if !ok {
			return fmtp.Errorf("project not found : %s", e.itemName)
		}
		pr[i].itemId = v
		v, ok = rom[e.roleName]
		if !ok {
			return fmtp.Errorf("role not found : %s", e.roleName)
		}
		pr[i].roleId = v
	}

	for i, e := range dr {
		v, ok := rom[e.roleName]
		if !ok {
			return fmtp.Errorf("role not found : %s", e.roleName)
		}
		dr[i].roleId = v
	}

	for i, e := range dri {
		v, ok := rom[e.roleName]
		if !ok {
			return fmtp.Errorf("role not found : %s", e.roleName)
		}
		dri[i].roleId = v
	}

	// update project role
	agencyId := newRoles.Id()
	toCreate, toDelete := calcDiff(epr, pr)
	log.Printf("[DEBUG] projects toCreate: %v, toDelete: %v", toCreate, toDelete)
	if err := createAgencyRole(c, newRoles.Id(), toCreate, agency.CreateAgencyProjectRole); err != nil {
		return err
	}
	if err := deleteAgencyRole(c, agencyId, toCreate, agency.DeleteAgencyProjectRole); err != nil {
		return err
	}

	// update domain role
	toCreate, toDelete = calcDiff(edr, dr)
	log.Printf("[DEBUG] projects toCreate: %v, toDelete: %v", toCreate, toDelete)
	if err := createAgencyRole(c, agencyId, toCreate, agency.CreateAgencyDomainRole); err != nil {
		return err
	}
	if err := deleteAgencyRole(c, agencyId, toCreate, agency.DeleteAgencyDomainRole); err != nil {
		return err
	}

	// update all resources role
	toCreate, toDelete = calcDiff(edri, dri)
	log.Printf("[DEBUG] projects toCreate: %v, toDelete: %v", toCreate, toDelete)
	if err := createAgencyRole(c, agencyId, toCreate, agency.CreateAgencyDomainInheritedRole); err != nil {
		return err
	}
	if err := deleteAgencyRole(c, agencyId, toCreate, agency.DeleteAgencyDomainInheritedRole); err != nil {
		return err
	}
	return nil
}

func createAgencyRole(c *golangsdk.ServiceClient, agencyId string, roles []RoleItem, fun func(*golangsdk.ServiceClient, string, string, string) error) error {
	if len(roles) == 0 {
		return nil
	}
	for _, r := range roles {
		if err := fun(c, agencyId, r.itemId, r.roleId); err != nil {
			return err
		}
	}
	return nil
}

func deleteAgencyRole(c *golangsdk.ServiceClient, agencyId string, roles []RoleItem, fun func(*golangsdk.ServiceClient, string, string, string) error) error {
	//for _, r := range roles {
	//	if err := fun(c, agencyId, r.itemId, r.roleId); err != nil {
	//		return err
	//	}
	//}
	return nil
}

// query and construct project map: project.name => project.id
func getProjectMap(cfg *config.Config, hcsConfig *config.HcsConfig, opts projects.ListOpts) (map[string]string, error) {
	c, err := hcsConfig.IdentityV3Client(cfg.Region)
	if err != nil {
		return nil, err
	}
	pjs := projects.List(c, opts)
	pjm := make(map[string]string)
	err = pjs.EachPage(func(page pagination.Page) (bool, error) {
		ps, err := projects.ExtractProjects(page)
		if err != nil {
			return false, err
		}
		for _, p := range ps {
			pjm[p.Name] = p.ID
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return pjm, nil
}

// query and construct role map: role.display_name => role.id
func getRoleMap(c *golangsdk.ServiceClient, opts role.ListOpts) (map[string]string, error) {
	rm := make(map[string]string)
	for {
		resp := role.List(c, opts)
		if resp.Err != nil {
			return nil, resp.Err
		}
		rs, t, err := resp.Extract()
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			rm[r.DisplayName] = r.ID
		}
		opts.Start += len(rs)
		if opts.Start >= t {
			break
		}
	}
	return rm, nil
}

// calculate symmetric difference
func calcDiff(existingRoles []RoleItem, newRoles []RoleItem) (toCreate []RoleItem, toDelete []RoleItem) {
	ex, ne := make(map[string]bool), make(map[string]bool)
	f := func(r RoleItem) string {
		return r.roleName + ":" + r.itemName
	}
	for _, e := range existingRoles {
		ex[f(e)] = true
	}
	for _, e := range newRoles {
		ne[f(e)] = true
	}
	for _, e := range existingRoles {
		if !ne[f(e)] {
			toDelete = append(toDelete, e)
		}
	}
	for _, e := range newRoles {
		if !ex[f(e)] {
			toCreate = append(toCreate, e)
		}
	}
	return
}

func setAgencyRole(cfg *config.Config, d *schema.ResourceData, roles []agency.AgencyRole) error {
	epr, edr, edri := buildAgencyRoleFromResponse(cfg.DomainID, roles)

	// construct project_role list according to provided configuration, including:
	// 1. duplicate project name in different blocks
	// 2. duplicate role name in different blocks
	dpr := make([]interface{}, 0)
	if d.Get("project_role") != nil {
		dpr = d.Get("project_role").(*schema.Set).List()
	}
	eprs := make(map[string]map[string]bool, len(epr))
	for _, i := range epr {
		if eprs[i.itemName] == nil {
			eprs[i.itemName] = make(map[string]bool)
		}
		eprs[i.itemName][i.roleName] = true
	}
	prd := make([]map[string]interface{}, 0, len(eprs))
	for _, r := range dpr {
		m := r.(map[string]interface{})
		p := m["project"].(string)
		s := make([]string, 0)
		for _, i := range m["roles"].(*schema.Set).List() {
			if eprs[p][i.(string)] {
				s = append(s, i.(string))
			}
		}
		slices.Sort(s)
		prd = append(prd, map[string]interface{}{"project": p, "roles": s})
	}
	slices.SortFunc(prd, func(a, b map[string]interface{}) int {
		return strings.Compare(a["project"].(string), b["project"].(string))
	})

	// construct domain_roles
	drd := make([]string, len(edr))
	for i, d := range edr {
		drd[i] = d.roleName
	}
	slices.Sort(drd)

	// construct all_resources_roles
	drid := make([]string, len(edri))
	for i, d := range edri {
		drid[i] = d.roleName
	}
	slices.Sort(drid)

	return multierror.Append(nil,
		d.Set("project_role", prd),
		d.Set("domain_roles", drd),
		d.Set("all_resources_roles", drid),
	)
}

func buildAgencyRoleFromResponse(domainId string, roles []agency.AgencyRole) (projectRole []RoleItem, domainRole []RoleItem, domainRoleInherited []RoleItem) {
	for _, r := range roles {
		for _, p := range r.Projects {
			if p.ID == domainId {
				domainRole = append(domainRole, RoleItem{roleId: r.ID, roleName: r.DisplayName, itemId: "", itemName: ""})
			} else {
				projectRole = append(projectRole, RoleItem{roleId: r.ID, roleName: r.DisplayName, itemId: p.ID, itemName: p.Name})
			}
		}

	}
	return
}

func buildAgencyRoleFromSchema(cfg *config.Config, d *schema.ResourceData) (projectRole []RoleItem, domainRole []RoleItem, domainRoleInherited []RoleItem) {
	domainId := cfg.DomainID
	domainName := cfg.DomainName

	if d.Get("project_role") != nil {
		pr := d.Get("project_role").(*schema.Set)
		for _, i := range pr.List() {
			e := i.(map[string]interface{})
			projectName := e["project"].(string)
			roles := e["roles"].(*schema.Set)
			for _, r := range roles.List() {
				projectRole = append(projectRole, RoleItem{roleId: "", roleName: r.(string), itemId: "", itemName: projectName})
			}
		}
	}

	if d.Get("domain_roles") != nil {
		dr := d.Get("domain_roles").(*schema.Set)
		for _, i := range dr.List() {
			domainRole = append(domainRole, RoleItem{roleId: "", roleName: i.(string), itemId: domainId, itemName: domainName})
		}
	}

	if d.Get("all_resources_roles") != nil {
		dri := d.Get("all_resources_roles").(*schema.Set)
		for _, i := range dri.List() {
			domainRoleInherited = append(domainRoleInherited, RoleItem{roleId: "", roleName: i.(string), itemId: domainId, itemName: domainName})
		}
	}

	return
}

type RoleItem struct {
	roleId   string
	roleName string
	itemId   string
	itemName string
}
