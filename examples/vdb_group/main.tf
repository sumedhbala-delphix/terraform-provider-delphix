/**
* Summary: This example demonstrates how to:
* 1) Create multiple VDBs with different configurations
* 2) Group VDBs into a VDB group
* 3) Manage tags for VDB groups
* 4) Handle multiple environments
*/

terraform {
  required_providers {
    delphix = {
      version = ">=3.3.2"
      source  = "delphix-integrations/delphix"
    }
  }
}

// Configure the Delphix provider
provider "delphix" {
  tls_insecure_skip = true
  key               = "1.XXXX"  # Replace with your API key
  host              = "HOSTNAME" # Replace with your Delphix host
}

// Define VDB configurations for different environments
locals {
  vdbs = {
    "prod" = {
      snapshot_id = "6-ORACLE_DB_CONTAINER-7"
      name        = "prod-db"
      environment = "production"
    }
    "staging" = {
      snapshot_id = "6-ORACLE_DB_CONTAINER-1"
      name        = "staging-db"
      environment = "staging"
    }
    "dev" = {
      snapshot_id = "6-ORACLE_DB_CONTAINER-5"
      name        = "dev-db"
      environment = "development"
    }
  }

  // Common tags for all resources
  common_tags = {
    "managed_by" = "terraform"
    "team"       = "platform"
  }
}

// Create VDBs for each environment
resource "delphix_vdb" "vdb_provision" {
  for_each               = local.vdbs
  name                   = each.value.name
  source_data_id         = each.value.snapshot_id
  auto_select_repository = true

  tags {
    key   = "environment"
    value = each.value.environment
  }
  tags {
    key   = "managed_by"
    value = local.common_tags["managed_by"]
  }
  tags {
    key   = "team"
    value = local.common_tags["team"]
  }
}

// Create a VDB group for production databases
resource "delphix_vdb_group" "production" {
  name    = "production-group"
  vdb_ids = [delphix_vdb.vdb_provision["prod"].id]

  tags {
    key   = "environment"
    value = "production"
  }
  tags {
    key   = "purpose"
    value = "production-databases"
  }
  tags {
    key   = "managed_by"
    value = local.common_tags["managed_by"]
  }
  tags {
    key   = "team"
    value = local.common_tags["team"]
  }
}

// Create a VDB group for non-production databases
resource "delphix_vdb_group" "non_production" {
  name    = "non-production-group"
  vdb_ids = [
    delphix_vdb.vdb_provision["staging"].id,
    delphix_vdb.vdb_provision["dev"].id
  ]

  tags {
    key   = "environment"
    value = "non-production"
  }
  tags {
    key   = "purpose"
    value = "development-and-staging"
  }
  tags {
    key   = "managed_by"
    value = local.common_tags["managed_by"]
  }
  tags {
    key   = "team"
    value = local.common_tags["team"]
  }
}

// Output the VDB group IDs
output "production_group_id" {
  description = "ID of the production VDB group"
  value       = delphix_vdb_group.production.id
}

output "non_production_group_id" {
  description = "ID of the non-production VDB group"
  value       = delphix_vdb_group.non_production.id
}
