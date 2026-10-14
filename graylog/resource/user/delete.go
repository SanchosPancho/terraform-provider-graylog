package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client"
	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client/api"
)

func destroy(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	var deleted_id string
	if api.UsesUserIDs(cl.APIVersion) {
		deleted_id = "id/" + d.Id()
	} else {
		deleted_id = d.Id()
	}
	if _, err := cl.User.Delete(ctx, deleted_id); err != nil {
		return err
	}
	return nil
}
