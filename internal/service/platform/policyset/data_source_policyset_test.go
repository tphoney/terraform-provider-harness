package policyset_test

import (
	"fmt"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyset(t *testing.T) {
	id := t.Name() + utils.RandStringBytes(6)
	resourceName := "data.harness_platform_policyset.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyset(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "name", id),
				),
			},
		},
	})
}

func testAccDataSourcePolicyset(id string) string {
	return fmt.Sprintf(`
		resource "harness_platform_policyset" "test" {
			identifier = "%[1]s"
			name = "%[1]s"
		}

		data "harness_platform_policyset" "test" {
			identifier = harness_platform_policyset.test.identifier
			name = harness_platform_policyset.test.name
		}
	`, id)
}
