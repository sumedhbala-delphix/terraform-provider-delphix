# Resource: delphix_vdb_group

A VDB group is a collection of virtual databases and datasets. The `delphix_vdb_group` resource allows you to create and manage such collections.

## Example Usage

### Basic VDB Group
```hcl
resource "delphix_vdb" "my_vdb" {
  auto_select_repository = true
  source_data_id         = "1-ORACLE_DB_CONTAINER-2"
}

resource "delphix_vdb_group" "my_group" {
  name    = "my-vdb-group"
  vdb_ids = [delphix_vdb.my_vdb.id]
}
```

### VDB Group with Tags
```hcl
resource "delphix_vdb" "my_vdb" {
  auto_select_repository = true
  source_data_id         = "1-ORACLE_DB_CONTAINER-2"
}

resource "delphix_vdb_group" "my_group" {
  name    = "my-vdb-group"
  vdb_ids = [delphix_vdb.my_vdb.id]
  
  tags {
    key   = "environment"
    value = "production"
  }
  tags {
    key   = "team"
    value = "platform"
  }
  tags {
    key   = "description"
    value = "Production VDB group for Oracle database"
  }
}
```

### VDB Group with Multiple VDBs
```hcl
resource "delphix_vdb" "vdb1" {
  auto_select_repository = true
  source_data_id         = "1-ORACLE_DB_CONTAINER-2"
}

resource "delphix_vdb" "vdb2" {
  auto_select_repository = true
  source_data_id         = "1-ORACLE_DB_CONTAINER-3"
}

resource "delphix_vdb_group" "my_group" {
  name    = "my-vdb-group"
  vdb_ids = [delphix_vdb.vdb1.id, delphix_vdb.vdb2.id]
  
  tags {
    key   = "environment"
    value = "production"
  }
}
```

## Argument Reference

* `id` - (Computed) A unique identifier for the VDB group.

* `name` - (Required) A unique name for the VDB group.

* `vdb_ids` - (Optional) The list of VDB IDs to include in this VDB group.

* `tags` - (Optional) A list of tags to associate with the VDB group. Each tag has the following attributes:
  * `key` - (Required) The tag key.
  * `value` - (Required) The tag value. Must be between 1 and 4000 characters.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier of the VDB group.
* `name` - The name of the VDB group.
* `vdb_ids` - The list of VDB IDs in this VDB group.
* `tags` - The list of tags associated with the VDB group.

## Notes

* Tag values must be between 1 and 4000 characters.
* Empty tag values are not supported.
* Tags can be updated independently of other VDB group properties.
* Multiple tags with the same key are supported.
* Special characters and Unicode are supported in tag values.
