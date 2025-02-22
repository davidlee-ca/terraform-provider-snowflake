package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_View(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_View2(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%';"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_ViewWithCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func viewConfig(n string, copyGrants bool, q string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
}

resource "snowflake_view" "test" {
	name        = "%v"
	comment     = "Terraform test resource"
	database    = snowflake_database.test.name
	schema      = "PUBLIC"
	is_secure   = true
	or_replace  = %t
	copy_grants = %t
	statement   = "%s"
}
`, n, n, copyGrants, copyGrants, q)
}
