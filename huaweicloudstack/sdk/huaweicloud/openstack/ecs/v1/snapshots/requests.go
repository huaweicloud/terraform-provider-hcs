package snapshots

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

type CreateInstanceSnapshotOpts struct {
	// Indicates the snapshot name.
	Name string `json:"name" required:"true"`
	// Specifies whether to create an ECS snapshot. Only true or false of the string type.
	InstanceSnapshot string `json:"instance_snapshot" required:"true"`
	// Specifies the ECS ID to create an ECS snapshot.
	ServerId string `json:"-" required:"true"`
}

type RollBackInstanceSnapshotOpts struct {
	ImageRef string `json:"imageRef" required:"true"`
}

func Create(c *golangsdk.ServiceClient, opts CreateInstanceSnapshotOpts) (r JobResult) {
	b, err := golangsdk.BuildRequestBody(opts, "createImage")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(actionURL(c, opts.ServerId), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{202}})
	return
}

// List returns a Pager that allows you to iterate over a collection of ECS snapshots.
func List(client *golangsdk.ServiceClient, serverId string) pagination.Pager {
	return pagination.NewPager(client, getURL(client, serverId), func(r pagination.PageResult) pagination.Page {
		return SnapshotPage{pagination.SinglePageBase(r)}
	})
}

func Get(client *golangsdk.ServiceClient, serverId string, imageId string) (image QueryImage, err error) {
	pages, err := List(client, serverId).AllPages()
	if err != nil {
		return image, err
	}

	images, err := ExtractSnapshots(pages)
	if err != nil {
		return image, err
	}

	for _, image := range images {
		if image.Id == imageId {
			return image, nil
		}
	}

	return image, fmt.Errorf("snpashot %s not found", imageId)
}

func Rollback(c *golangsdk.ServiceClient, serverId string, opts RollBackInstanceSnapshotOpts) (r JobResult) {
	b, err := golangsdk.BuildRequestBody(opts, "rebuild")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(actionURL(c, serverId), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
