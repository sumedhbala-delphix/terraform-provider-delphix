package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVdb_create_vdb_group_positive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVdbPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVdbGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDctVDBGroupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group")),
			},
		},
	})
}

func TestAccVdb_create_vdb_group_with_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVdbPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVdbGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDctVDBGroupConfigWithTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "2"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "test"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.key", "org"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.value", "dev"),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigWithUpdatedTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "3"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "prod"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.key", "org"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.value", "dev"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.key", "team"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.value", "platform"),
				),
			},
		},
	})
}

func TestAccVdb_create_vdb_group_with_complex_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVdbPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVdbGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDctVDBGroupConfigWithComplexTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "5"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "test environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.key", "org"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.value", "dev"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.key", "description"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.value", "This is a test VDB group with special chars: !@#$%^&*()"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.3.key", "unicode"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.3.value", "测试标签"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.4.key", "min_length"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.4.value", "a"),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigWithDuplicateKeys(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "3"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "prod"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.1.value", "staging"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.key", "environment"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.2.value", "dev"),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigWithLongTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "1"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "long_value"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", strings.Repeat("a", 1000)),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccVdb_create_vdb_group_with_many_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVdbPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVdbGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDctVDBGroupConfigWithManyTags(50), // Test with 50 tags
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "50"),
					// Verify first and last tags
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "tag_0"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "value_0"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.49.key", "tag_49"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.49.value", "value_49"),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigWithManyTags(100), // Update to 100 tags
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "100"),
					// Verify first and last tags
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.key", "tag_0"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.0.value", "value_0"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.99.key", "tag_99"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.99.value", "value_99"),
				),
			},
			{
				Config: testAccCheckDctVDBGroupConfigBasic(), // Remove all tags
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDctVdbGroupResourceExists("delphix_vdb.new", "delphix_vdb_group.new_group"),
					resource.TestCheckResourceAttr("delphix_vdb_group.new_group", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccCheckDctVDBGroupConfigBasic() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name                   = "my-vdb-group-name"
		vdb_ids                = [delphix_vdb.new.id]
	}
	`, datasource_id)
}

func testAccCheckDctVDBGroupConfigWithTags() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]
		tags {
			key   = "environment"
			value = "test"
		}
		tags {
			key   = "org"
			value = "dev"
		}
	}
	`, datasource_id)
}

func testAccCheckDctVDBGroupConfigWithUpdatedTags() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]
		tags {
			key   = "environment"
			value = "prod"
		}
		tags {
			key   = "org"
			value = "dev"
		}
		tags {
			key   = "team"
			value = "platform"
		}
	}
	`, datasource_id)
}

func testAccCheckDctVDBGroupConfigWithComplexTags() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]
		tags {
			key   = "environment"
			value = "test environment"
		}
		tags {
			key   = "org"
			value = "dev"
		}
		tags {
			key   = "description"
			value = "This is a test VDB group with special chars: !@#$%%^&*()"
		}
		tags {
			key   = "unicode"
			value = "测试标签"
		}
		tags {
			key   = "min_length"
			value = "a"
		}
	}
	`, datasource_id)
}

func testAccCheckDctVDBGroupConfigWithDuplicateKeys() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]
		tags {
			key   = "environment"
			value = "prod"
		}
		tags {
			key   = "environment"
			value = "staging"
		}
		tags {
			key   = "environment"
			value = "dev"
		}
	}
	`, datasource_id)
}

func testAccCheckDctVDBGroupConfigWithLongTags() string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]
		tags {
			key   = "long_value"
			value = "%s"
		}
	}
	`, datasource_id, strings.Repeat("a", 1000))
}

func testAccCheckDctVDBGroupConfigWithManyTags(count int) string {
	datasource_id := os.Getenv("DATASOURCE_ID")
	tags := ""
	for i := 0; i < count; i++ {
		tags += fmt.Sprintf(`
		tags {
			key   = "tag_%d"
			value = "value_%d"
		}`, i, i)
	}
	return fmt.Sprintf(`
	resource "delphix_vdb" "new" {
		auto_select_repository = true
		source_data_id         = "%s"
	}
	resource "delphix_vdb_group" "new_group" {
		name    = "my-vdb-group-name"
		vdb_ids = [delphix_vdb.new.id]%s
	}
	`, datasource_id, tags)
}

func testAccCheckDctVdbGroupResourceExists(vdbResourceName string, vdbGroupResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vdbGroupResource, ok := s.RootModule().Resources[vdbGroupResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", vdbGroupResourceName)
		}

		vdbGroupId := vdbGroupResource.Primary.ID
		if vdbGroupId == "" {
			return fmt.Errorf("No VdbGroupID set")
		}

		vdbResource, ok := s.RootModule().Resources[vdbResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", vdbResourceName)
		}

		vdbId := vdbResource.Primary.ID
		if vdbGroupId == "" {
			return fmt.Errorf("No VDbID set")
		}

		client := testAccProvider.Meta().(*apiClient).client

		res, _, err := client.VDBGroupsAPI.GetVdbGroup(context.Background(), vdbGroupId).Execute()
		if err != nil {
			return err
		}

		vdbIds := res.GetVdbIds()
		if !reflect.DeepEqual(vdbIds, []string{vdbId}) {
			return fmt.Errorf("Expected the vdb_id in VDB Group vdb_ids property")
		}

		return nil
	}
}

func testAccCheckVdbGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "delphix_vdb_group" {
			continue
		}

		vdbGroupId := rs.Primary.ID

		_, httpResp, _ := client.VDBGroupsAPI.GetVdbGroup(context.Background(), vdbGroupId).Execute()
		if httpResp == nil {
			return fmt.Errorf("VDB Group has not been deleted")
		}

		if httpResp.StatusCode != 404 {
			return fmt.Errorf("Exepcted a 404 Not Found for a deleted VDB Group but got %d", httpResp.StatusCode)
		}
	}

	return nil
}
