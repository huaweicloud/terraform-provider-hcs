package vdc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/user"
)

func ResourceVdcUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcUserCreate,
		ReadContext:   resourceVdcUserRead,
		UpdateContext: resourceVdcUserUpdate,
		DeleteContext: resourceVdcUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVdcUserInstanceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Required:  false,
				Sensitive: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Required: false,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"vdc_id": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.StringInSlice([]string{"LOCAL_AUTH", "SAML_AUTH", "LDAP_AUTH", "MACHINE_USER"}, false),
				Default:      "LOCAL_AUTH",
			},
			"access_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.StringInSlice([]string{"default", "programmatic", "console"}, false),
				Default:      "default",
			},
		},
	}
}

func resourceVdcUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	authType := d.Get("auth_type").(string)
	authTypeVal, ok := user.ApiAuthType[authType]
	if !ok {
		return diag.Errorf(`The values supported by auth_type are: "LOCAL_AUTH", "SAML_AUTH", "LDAP_AUTH", "MACHINE_USER".`)
	}

	isLocalAuth := authTypeVal == user.ApiAuthType["LOCAL_AUTH"]
	isMachineUser := authTypeVal == user.ApiAuthType["MACHINE_USER"]
	pwd := d.Get("password").(string)
	pwdIsNull := pwd == ""
	if (isLocalAuth || isMachineUser) && pwdIsNull {
		return diag.Errorf(`When auth_type is set to LOCAL_AUTH or MACHINE_USER, you need to specify a value for password.`)
	}

	isSamlUser := authTypeVal == user.ApiAuthType["SAML_AUTH"]
	isLdapUser := authTypeVal == user.ApiAuthType["LDAP_AUTH"]
	hasPwd := pwd != ""
	if (isSamlUser || isLdapUser) && hasPwd {
		return diag.Errorf(`When auth_type is set to SAML_AUTH or LDAP_AUTH, do not specify password.`)
	}

	mode := d.Get("access_mode").(string)
	accessModeVal, ok := user.ApiAccessMode[mode]
	if !ok {
		return diag.Errorf(`The values supported by access_mode are: "default", "programmatic", "console".`)
	}

	unMachineUser := authTypeVal != user.ApiAuthType["MACHINE_USER"]
	isProgrammaticMode := accessModeVal == user.ApiAccessMode["programmatic"]
	unProgrammaticMode := accessModeVal != user.ApiAccessMode["programmatic"]
	if isMachineUser && unProgrammaticMode {
		return diag.Errorf(`When auth_type is set to MACHINE_USER, access_mode must be set to programmatic.`)
	}
	if isProgrammaticMode && unMachineUser {
		return diag.Errorf(`When access_mode is set to programmatic, auth_type must be set to MACHINE_USER.`)
	}

	hcsConfig := config.GetHcsConfig(meta)
	userClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VDC user client %s", err)
	}

	vdcId := d.Get("vdc_id").(string)
	createOpts := user.CreateOpts{
		Name:        d.Get("name").(string),
		Password:    pwd,
		DisplayName: d.Get("display_name").(string),
		AuthType:    authTypeVal,
		Enabled:     d.Get("enabled").(bool),
		Description: d.Get("description").(string),
		AccessMode:  accessModeVal,
	}

	addUser, err := user.Create(userClient, vdcId, createOpts).ToExtract()

	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack VDC user: %s", err)
	}

	d.SetId(addUser.ID)

	return resourceVdcUserRead(ctx, d, meta)
}

func resourceVdcUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	userClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user client : %s", err)
	}

	userDetail, err := user.Get(userClient, d.Id()).ToExtract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloudStack VDC user")
	}

	typeVal := user.AuthType[userDetail.AuthType]
	modeVal := user.AccessMode[userDetail.AccessMode]
	mErr := multierror.Append(nil,
		d.Set("id", userDetail.ID),
		d.Set("name", userDetail.Name),
		d.Set("display_name", userDetail.DisplayName),
		d.Set("enabled", userDetail.Enabled),
		d.Set("description", userDetail.Description),
		d.Set("auth_type", typeVal),
		d.Set("access_mode", modeVal),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting HuaweiCloudStack VDC user fields: %w", err)
	}

	return nil
}

func resourceVdcUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("vdc_id", "name") {
		return diag.Errorf(`Unsupported attribute values for modification: "vdc_id", "name".`)
	}

	authType := d.Get("auth_type").(string)
	authTypeVal, ok := user.ApiAuthType[authType]
	if !ok {
		return diag.Errorf(`The values supported by auth_type are: "LOCAL_AUTH", "SAML_AUTH", "LDAP_AUTH", "MACHINE_USER".`)
	}

	mode := d.Get("access_mode").(string)
	accessModeVal, ok := user.ApiAccessMode[mode]
	if !ok {
		return diag.Errorf(`The values supported by access_mode are: "default", "programmatic", "console".`)
	}

	isChangeAccessMode := d.HasChanges("access_mode")
	isChangeAuthType := d.HasChanges("auth_type")
	unProgrammaticMode := accessModeVal != user.ApiAccessMode["programmatic"]
	if isChangeAccessMode && isChangeAuthType && unProgrammaticMode {
		return diag.Errorf(`The attribute value "auth_type" is not supported for modification.`)
	}

	isProgrammaticMode := accessModeVal == user.ApiAccessMode["programmatic"]
	unMachineUser := authTypeVal != user.ApiAuthType["MACHINE_USER"]
	isMachineUser := authTypeVal == user.ApiAuthType["MACHINE_USER"]
	if isMachineUser && unProgrammaticMode {
		return diag.Errorf(`When auth_type is set to MACHINE_USER, access_mode must be set to programmatic.`)
	}
	if isProgrammaticMode && unMachineUser {
		return diag.Errorf(`When access_mode is set to programmatic, auth_type must be set to MACHINE_USER.`)
	}

	hcsConfig := config.GetHcsConfig(meta)
	userClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user client : %s", err)
	}

	if d.HasChanges("password") {
		pwd := d.Get("password").(string)
		pwdIsNull := pwd == ""
		isLocalAuth := authTypeVal == user.ApiAuthType["LOCAL_AUTH"]
		if (isLocalAuth || isMachineUser) && pwdIsNull {
			return diag.Errorf(`When auth_type is set to LOCAL_AUTH or MACHINE_USER, you need to specify a value for password.`)
		}

		isSamlUser := authTypeVal == user.ApiAuthType["SAML_AUTH"]
		isLdapUser := authTypeVal == user.ApiAuthType["LDAP_AUTH"]
		if isSamlUser || isLdapUser {
			return diag.Errorf(`When auth_type is set to SAML_AUTH or LDAP_AUTH, do not specify password.`)
		}

		updatePwdOpts := user.PutPwdOpts{
			Password: d.Get("password").(string),
		}

		_, err = user.UpPwd(userClient, updatePwdOpts, d.Id()).ToExtract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloudStack VDC user password: %s", err)
		}
	}

	enabledVal := d.Get("enabled").(bool)
	if !d.HasChanges("enabled") && !enabledVal {
		return diag.Errorf(`You cannot modify information about disabled users: "display_name", "description", "access_mode".`)
	}
	if d.HasChanges("display_name", "enabled", "description", "access_mode") {
		updateOpts := user.PutOpts{
			DisplayName: d.Get("display_name").(string),
			Enabled:     enabledVal,
			Description: d.Get("description").(string),
			AccessMode:  accessModeVal,
		}

		_, err = user.Update(userClient, updateOpts, d.Id(), false).ToExtract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloudStack VDC user: %s", err)
		}
	}

	return resourceVdcUserRead(ctx, d, meta)
}

func resourceVdcUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	userClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user client : %s", err)
	}

	err = user.Delete(userClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VDC user %s: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func resourceVdcUserInstanceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	userClient, err := hcsConfig.VdcClient(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VDC user client: %s", err)
	}

	userDetail, err := user.Get(userClient, d.Id()).ToExtract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "VDC user instance")
	}

	typeVal := user.AuthType[userDetail.AuthType]
	modeVal := user.AccessMode[userDetail.AccessMode]
	mErr := multierror.Append(nil,
		d.Set("vdc_id", userDetail.VdcId),
		d.Set("id", userDetail.ID),
		d.Set("name", userDetail.Name),
		d.Set("display_name", userDetail.DisplayName),
		d.Set("enabled", userDetail.Enabled),
		d.Set("description", userDetail.Description),
		d.Set("auth_type", typeVal),
		d.Set("access_mode", modeVal),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmtp.Errorf("error setting HuaweiCloudStack VDC user fields: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
