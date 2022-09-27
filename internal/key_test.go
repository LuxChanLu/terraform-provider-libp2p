package internal_test

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceKey(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: `
					resource "libp2p_key" "test" {
						type = "ED25519"
					}
				`,
				Check: r.ComposeAggregateTestCheckFunc(
					r.TestCheckResourceAttrSet("libp2p_key.test", "private"),
					r.TestCheckResourceAttrSet("libp2p_key.test", "public"),
					r.TestCheckResourceAttrSet("libp2p_key.test", "peer_id"),
					r.TestCheckResourceAttrSet("libp2p_key.test", "id"),
				),
			},
		},
	})
}
