//go:build tm || ALL || functional

package vcd

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdTmContentLibrary(t *testing.T) {
	preTestChecks(t)

	skipIfNotSysAdmin(t)
	skipIfNotTm(t)

	vCenterHcl, vCenterHclRef := getVCenterHcl(t)
	nsxManagerHcl, nsxManagerHclRef := getNsxManagerHcl(t)
	regionHcl, regionHclRef := getRegionHcl(t, vCenterHclRef, nsxManagerHclRef)

	var params = StringMap{
		"Name":                t.Name(),
		"RegionId":            fmt.Sprintf("%s.id", regionHclRef),
		"RegionStoragePolicy": testConfig.Tm.RegionStoragePolicy,
		"Tags":                "tm",
	}
	testParamsNotEmpty(t, params)

	preRequisites := vCenterHcl + nsxManagerHcl + regionHcl

	// TODO: TM: There shouldn't be a need to create `preRequisites` separately, but region
	// creation fails if it is spawned instantly after adding vCenter, therefore this extra step
	// give time (with additional 'refresh' and 'refresh storage policies' operations on vCenter)
	skipBinaryTest := "# skip-binary-test: prerequisite buildup for acceptance tests"
	configText0 := templateFill(vCenterHcl+nsxManagerHcl+skipBinaryTest, params)
	params["FuncName"] = t.Name() + "-step0"

	configText1 := templateFill(preRequisites+testAccVcdTmContentLibraryStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(preRequisites+testAccVcdTmContentLibraryStep2, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resourceName := "vcd_tm_content_library.cl"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText0,
			},
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", t.Name()),
					resource.TestCheckResourceAttr(resourceName, "description", t.Name()),
					resource.TestCheckResourceAttr(resourceName, "storage_class_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_attach", "true"), // TODO: TM: Test with false
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttr(resourceName, "is_shared", "true"),        // TODO: TM: Test with false
					resource.TestCheckResourceAttr(resourceName, "is_subscribed", "false"),   // TODO: TM: Test with true
					resource.TestCheckResourceAttr(resourceName, "library_type", "PROVIDER"), // TODO: TM: Test with tenant catalog
					resource.TestMatchResourceAttr(resourceName, "owner_org_id", regexp.MustCompile("urn:vcloud:org:")),
					resource.TestCheckResourceAttr(resourceName, "subscription_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "version_number", "1"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resourceFieldsEqual(resourceName, "data.vcd_tm_content_library.cl_ds", nil),
				),
			},
			{
				ResourceName:      "vcd_tm_content_library.cl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     params["Name"].(string),
			},
		},
	})

	postTestChecks(t)
}

const testAccVcdTmContentLibraryStep1 = `
data "vcd_tm_region_storage_policy" "sp" {
  region_id = {{.RegionId}}
  name      = "{{.RegionStoragePolicy}}"
}

resource "vcd_tm_content_library" "cl" {
  name = "{{.Name}}"
  description = "{{.Name}}"
  storage_class_ids = [
    data.vcd_tm_region_storage_policy.sp.id
  ]
}
`

const testAccVcdTmContentLibraryStep2 = testAccVcdTmContentLibraryStep1 + `
data "vcd_tm_content_library" "cl_ds" {
  name = vcd_tm_content_library.cl.name
}
`
