package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/SanchosPancho/terraform-provider-graylog/graylog/client"
)

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, err := getDataFromResourceData(d)
	if cl.APIVersion == "v4" {
		delete(data, "full_name")
	} else {
		delete(data, "first_name")
		delete(data, "last_name")
	}
	if err != nil {
		return err
	}
	if _, ok := data[keyPermissions]; !ok {
		data[keyPermissions] = []string{}
	}

	_, err = cl.User.Create(ctx, data)
	if err != nil {
		return err
	}
	if cl.APIVersion == "v4" {
		// Here we use v3 to retrieve the user information since the id is not returned by the API on creation
		new_data, _, _ := cl.User.Get(ctx, data[keyUsername].(string), "v3")
		err = d.Set("full_name", new_data["full_name"])
		if err != nil {
			return err
		}
		err = d.Set("external", new_data["external"])
		if err != nil {
			return err
		}
		err = d.Set("permissions", new_data["permissions"])
		if err != nil {
			return err
		}
		err = d.Set("read_only", new_data["read_only"])
		if err != nil {
			return err
		}
		err = d.Set("session_active", new_data["session_active"])
		if err != nil {
			return err
		}
		d.SetId(new_data["id"].(string))
	} else {
		d.SetId(data[keyUsername].(string))
	}
	return nil
}
