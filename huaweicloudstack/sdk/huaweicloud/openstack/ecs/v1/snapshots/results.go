package snapshots

import (
	"encoding/json"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// SnapshotPage stores a single page of all Snapshots results from a List call.
type SnapshotPage struct {
	pagination.SinglePageBase
}

type QueryImage struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	Status               string   `json:"status"`
	SnapshotFromInstance string   `json:"__snapshot_from_instance"`
	BaseImageRef         string   `json:"base_image_ref"`
	BlockDeviceMapping   string   `json:"block_device_mapping"`
	BootedVolume         string   `json:"booted_volume"`
	RootDeviceName       string   `json:"root_device_name"`
	ShutdownAt           string   `json:"shutdown_at"`
	Tags                 []string `json:"tags"`
	ContainerFormat      string   `json:"container_format"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
	DiskFormat           string   `json:"disk_format"`
	Visibility           string   `json:"visibility"`
	Self                 string   `json:"self"`
	MinDisk              uint16   `json:"min_disk"`
	Protected            bool     `json:"protected"`
	File                 string   `json:"file"`
	Owner                string   `json:"owner"`
	MinRam               uint16   `json:"min_ram"`
	Schema               string   `json:"schema"`
	Architecture         string   `json:"architecture"`

	BDM []BlockDeviceMapping
}

type BlockDeviceMapping struct {
	DeviceType          string `json:"device_type"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	SnapshotId          string `json:"snapshot_id"`
	DeviceName          string `json:"device_name"`
	DiskBus             string `json:"disk_bus"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type"`
	VolumeId            string `json:"volume_id"`
	VolumeSize          uint16 `json:"volume_size"`
	VolumeProjectId     string `json:"volume_project_id"`
	ExtendedInfo        string `json:"extended_info"`
}

type InstanceSnapshots struct {
	// ECS snapshot array.
	Images []QueryImage `json:"images"`

	// ECS snapshot view
	Schema string `json:"schema"`

	// Interface URI
	First string `json:"first"`
}

// ExtractSnapshots interprets a page of results as a slice of InstanceSnapshots.
func ExtractSnapshots(r pagination.Page) ([]QueryImage, error) {
	var s InstanceSnapshots
	err := (r.(SnapshotPage)).ExtractInto(&s)
	if err != nil {
		return nil, err

	}

	for i, _ := range s.Images {
		err = json.Unmarshal([]byte(s.Images[i].BlockDeviceMapping), &(s.Images[i].BDM))
		if err != nil {
			return s.Images, err
		}
	}

	return s.Images, err
}
