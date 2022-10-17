package sesv2_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfsesv2 "github.com/hashicorp/terraform-provider-aws/internal/service/sesv2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESV2DedicatedIPPool_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sesv2_dedicated_ip_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDedicatedIPPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedIPPoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedIPPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "pool_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ses", regexp.MustCompile(`dedicated-ip-pool/.+`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSESV2DedicatedIPPool_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sesv2_dedicated_ip_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDedicatedIPPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedIPPoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedIPPoolExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfsesv2.ResourceDedicatedIPPool(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDedicatedIPPoolDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Conn
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sesv2_dedicated_ip_pool" {
			continue
		}

		_, err := tfsesv2.FindDedicatedIPPoolByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			var nfe *types.NotFoundException
			if errors.As(err, &nfe) {
				return nil
			}
			return err
		}

		return create.Error(names.SESV2, create.ErrActionCheckingDestroyed, tfsesv2.ResNameDedicatedIPPool, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckDedicatedIPPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPPool, name, errors.New("not found"))
		}
		if rs.Primary.ID == "" {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPPool, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Conn
		ctx := context.Background()

		_, err := tfsesv2.FindDedicatedIPPoolByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPPool, rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Conn
	ctx := context.Background()

	_, err := conn.ListDedicatedIpPools(ctx, &sesv2.ListDedicatedIpPoolsInput{})
	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccDedicatedIPPoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_dedicated_ip_pool" "test" {
  pool_name = %[1]q
}
`, rName)
}
