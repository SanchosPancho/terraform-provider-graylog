package indexset

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/convert"
)

const (
	keyID                = "id"
	keyRotationStrategy  = "rotation_strategy"
	keyRetentionStrategy = "retention_strategy"
	keyDataTiering       = "data_tiering"
	keyDefault           = "default"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	// An omitted optional JSON value is represented by an empty string in the
	// Terraform SDK. It must become JSON null before decoding the request.
	if data[keyDataTiering] == "" {
		data[keyDataTiering] = "null"
	}
	if err := convert.JSONToData(data, keyRotationStrategy, keyRetentionStrategy, keyDataTiering); err != nil {
		return nil, err
	}

	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.DataToJSON(data, keyRotationStrategy, keyRetentionStrategy, keyDataTiering); err != nil {
		return err
	}

	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(data[keyID].(string))
	return nil
}
