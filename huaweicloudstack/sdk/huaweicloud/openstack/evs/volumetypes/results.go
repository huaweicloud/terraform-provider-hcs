package volumetypes

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// VolumeType contains all the information associated with an OpenStack Volume Type.
type VolumeType struct {
	// Unique identifier for the volume type.
	ID string `json:"id"`
	// Human-readable display name for the volume type.
	Name string `json:"name"`
	// Human-readable description for the volume type.
	Description string `json:"description"`
	// Arbitrary key-value pairs defined by the user.
	ExtraSpecs map[string]string `json:"extra_specs"`
	// Whether the volume type is publicly visible.
	IsPublic bool `json:"is_public"`
	// Qos Spec ID
	QosSpecID string `json:"qos_specs_id"`
}

// VolumeTypePage is a pagination.pager that is returned from a call to the List function.
type VolumeTypePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Volume Types.
func (r VolumeTypePage) IsEmpty() (bool, error) {
	volumetypes, err := ExtractVolumeTypes(r)
	return len(volumetypes) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (page VolumeTypePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"volume_type_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// ExtractVolumeTypes extracts and returns Volumes. It is used while iterating over a volumetypes.List call.
func ExtractVolumeTypes(r pagination.Page) ([]VolumeType, error) {
	var s []VolumeType
	err := ExtractVolumeTypesInto(r, &s)
	return s, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Volume Type object out of the commonResult object.
func (r commonResult) Extract() (*VolumeType, error) {
	var s VolumeType
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a volume type struct
func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "volume_type")
}

// ExtractVolumeTypesInto similar to ExtractInto but operates on a `list` of volume types
func ExtractVolumeTypesInto(r pagination.Page, v interface{}) error {
	return r.(VolumeTypePage).Result.ExtractIntoSlicePtr(v, "volume_types")
}
