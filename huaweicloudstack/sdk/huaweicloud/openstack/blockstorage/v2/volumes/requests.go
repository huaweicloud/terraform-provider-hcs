package volumes

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/evs/v2/cloudvolumes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVolumeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume. This object is passed to
// the volumes.Create function. For more information about these parameters,
// see the Volume object.
type CreateOpts struct {
	// The size of the volume, in GB
	Size int `json:"size" required:"true"`
	// The availability zone
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// ConsistencyGroupID is the ID of a consistency group
	ConsistencyGroupID string `json:"consistencygroup_id,omitempty"`
	// The enterprise project ID to which the volume belongs.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
	// The volume description
	Description string `json:"description,omitempty"`
	// One or more metadata key and value pairs to associate with the volume
	Metadata map[string]string `json:"metadata,omitempty"`
	// The volume name
	Name string `json:"name,omitempty"`
	// the ID of the existing volume snapshot
	SnapshotID string `json:"snapshot_id,omitempty"`
	// SourceReplica is a UUID of an existing volume to replicate with
	SourceReplica string `json:"source_replica,omitempty"`
	// the ID of the existing volume
	SourceVolID string `json:"source_volid,omitempty"`
	// Specifies the encryption configuration for the volume.
	// Required fields: cmk_id (KMS key ID), cipher (encryption algorithm, supports: AES256-XTS, SM4-XTS).
	EncryptionInfo *cloudvolumes.EncryptionInfoSpec `json:"encryption_info,omitempty"`
	// The ID of the image from which you want to create the volume.
	// Required to create a bootable volume.
	ImageID string `json:"imageRef,omitempty"`
	// The associated volume type
	VolumeType string `json:"volume_type,omitempty"`
	// Shared disk
	Multiattach bool `json:"multiattach,omitempty"`
}

// ToVolumeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToVolumeCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "volume")
}

// Create will create a new Volume based on the values in CreateOpts. To extract
// the Volume object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVolumeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// DeleteOptsBuilder is an interface by which can be able to build the query string
// of volume deletion.
type DeleteOptsBuilder interface {
	ToVolumeDeleteQuery() (string, error)
}

type DeleteOpts struct {
	//Specifies to delete all snapshots associated with the EVS disk.
	Cascade bool `q:"cascade"`
}

func (opts DeleteOpts) ToVolumeDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// Delete will delete the existing Volume with the provided ID
func Delete(client *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, id)
	if opts != nil {
		q, err := opts.ToVolumeDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += q
	}
	_, r.Err = client.Delete(url, nil)
	return
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeListQuery() (string, error)
}

// ListOpts holds options for listing Volumes. It is passed to the volumes.List
// function.
type ListOpts struct {
	// The availability zone
	AvailabilityZone string `q:"availability_zone"`

	// The Enterprise ProjectID
	EnterpriseProjectID string `q:"enterprise_project_id"`

	// Server ID
	ServerID string `q:"server_id"`

	// AllTenants will retrieve volumes of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	MetadataOrigin map[string]interface{}

	// Metadata will filter results based on specified metadata.
	Metadata map[string]string `q:"metadata"`

	// Name will filter by the specified volume name.
	Name string `q:"name"`

	// whether all snapshots of the volume are displayed in the volume list.
	WithSnapshot bool `q:"with_snapshot"`

	// Status will filter by the specified status.
	Status string `q:"status"`

	// TenantID will filter by a specific tenant/project ID.
	// Setting AllTenants is required for this.
	TenantID string `q:"project_id"`

	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// Tags will filter results based on specified tags.
	Tags map[string]string `q:"tags"`
}

// ToVolumeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeListQuery() (string, error) {
	handleMetadata(&opts)
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func handleMetadata(opts *ListOpts) {
	m := opts.MetadataOrigin
	nm := make(map[string]string, len(m))
	for k, v := range m {
		nm[k] = v.(string)
	}
	opts.Metadata = nm
}

// List returns Volumes optionally limited by the conditions provided in ListOpts.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VolumePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVolumeUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Volume. This object is passed
// to the volumes.Update function. For more information about the parameters, see
// the Volume object.
type UpdateOpts struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ToVolumeUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToVolumeUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "volume")
}

// Update will update the Volume with provided information. To extract the updated
// Volume from the response, call the Extract method on the UpdateResult.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// IDFromName is a convienience function that returns a server's ID given its name.
func IDFromName(client *golangsdk.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	listOpts := ListOpts{
		Name: name,
	}

	pages, err := List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractVolumes(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", golangsdk.ErrResourceNotFound{Name: name, ResourceType: "volume"}
	case 1:
		return id, nil
	default:
		return "", golangsdk.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "volume"}
	}
}
