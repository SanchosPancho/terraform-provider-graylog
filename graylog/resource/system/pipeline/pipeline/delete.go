package pipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client"
)

func destroy(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	if _, err := cl.Pipeline.Delete(ctx, d.Id()); err != nil {
		return fmt.Errorf("failed to delete a pipeline %s: %w", d.Id(), err)
	}
	return nil
}
