package pipeline

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/convert"
)

const (
	keyID = "id"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(data[keyID].(string))
	return nil
}
