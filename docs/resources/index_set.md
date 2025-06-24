# Resource: graylog_index_set

* [Example](https://github.com/terraform-provider-graylog/terraform-provider-graylog/blob/master/examples/v0.12/index_set.tf)
* [Source Code](https://github.com/terraform-provider-graylog/terraform-provider-graylog/blob/master/graylog/resource/system/indices/indexset/resource.go)

## Argument Reference

* `title` - (Required) the title of the Index Set. The data type is `string`.
* `index_prefix` - (Required, Forces new resource) the index prefix of the Index Set. The data type is `string`.
* `rotation_strategy_class` - (Required) the rotation strategy class of the Index Set. The data type is `string`.
* `rotation_strategy` - (Required) the rotation strategy of the Index Set. The data type is `JSON string`.
* `retention_strategy_class` - (Required) the retention strategy class of the Index Set. The data type is `string`.
* `retention_strategy` - (Required) the retention strategy of the Index Set. The data type is `JSON string`.
* `index_analyzer` - (Required) the Index Analyzer of the Index Set. The data type is `string`.
* `shards` - (Required) the number of shards of the Index Set. The data type is `int`.
* `description` - (Optional) the description of the Index Set. The data type is `string`.
* `replicas` - (Optional) the number of the replicas of the Index Set. The data type is `int`.
* `index_optimization_disabled` - (Optional) The data type is `bool`.
* `index_optimization_max_num_segments` - (Required) The data type is `int`.
* `default` - (Optional) The data type is `bool`.
* `field_type_refresh_interval` - (Optional) The data type is `int`.
* `use_legacy_rotation` - (Optional) The legacy rotation flag, added in graylog 6+, used with old(legacy) rotation\retention strategy(default false) `bool`.
* `can_be_default` - (Optional) whether this index set can be made the default. The data type is `bool`.
* `index_template_type` - (Optional) the index template type to apply (e.g. legacy or v2). The data type is `string`.
* `field_type_profile` - (Optional) the field type profile ID to attach to this index set. The data type is `string`.
* `data_tiering` - (Optional) the data-tiering configuration for this index set in JSON format (type, index_lifetime_min, index_lifetime_max). The data type is `JSON string`.
* `data_tiering_status` - (Computed) the current status of data-tiering for this index set. The data type is `string`.

## Attributes Reference

* `writable` - The data type is `bool`.
* `creation_date` - The date time when the Index Set is created. The data type is `string`.

## Import

`graylog_index_set` can be imported using the Index Set id, e.g.

```console
$ terraform import graylog_index_set.test 5c4acaefc9e77bbbbbbbbbbb
```
