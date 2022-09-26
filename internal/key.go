package internal

import (
	"fmt"

	b64 "encoding/base64"

	"github.com/hashicorp/terraform/helper/schema"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

var keyTypeMapping = map[string]int{
	"RSA":       crypto.RSA,
	"ED25519":   crypto.Ed25519,
	"SECP256K1": crypto.Secp256k1,
	"ECDSA":     crypto.ECDSA,
}

func keyItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key type",
				ValidateFunc: func(v interface{}, key string) ([]string, []error) {
					var errors []error
					var warnings []string
					value, ok := v.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("Expected %s to be string", key))
						return warnings, errors
					}
					_, ok = keyTypeMapping[value]
					if !ok {
						errors = append(errors, fmt.Errorf("Key type %s not found, available : (RSA, ED25519, SECP256K1, ECDSA)", key))
						return warnings, errors
					}
					return warnings, errors
				},
			},
			"bits": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Bits count for RSA key",
				ValidateFunc: func(v interface{}, key string) ([]string, []error) {
					var errors []error
					var warnings []string
					value, ok := v.(int)
					if !ok {
						return warnings, append(errors, fmt.Errorf("Expected %s to be int", key))
					}
					if value <= 0 {
						return warnings, append(errors, fmt.Errorf("Key size need to be more than 0"))
					}
					return warnings, errors
				},
			},
			"private": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "Private key encoded in base64",
			},
			"public": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "Public key encoded in base64",
			},
			"peerId": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "PeerID encoded in base58",
			},
		},
		Create: keyCreate,
		Read:   keyRead,
		Update: keyUpdate,
		Delete: keyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func keyCreate(d *schema.ResourceData, m interface{}) error {
	keyType := d.Get("type").(string)
	bits := d.Get("bits").(int)

	priv, pub, err := crypto.GenerateKeyPair(keyTypeMapping[keyType], bits)
	if err != nil {
		return err
	}
	peerId, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return err
	}
	privRaw, err := priv.Raw()
	if err != nil {
		return err
	}
	pubRaw, err := pub.Raw()
	if err != nil {
		return err
	}

	d.SetId(peerId.String())

	d.Set("private", b64.StdEncoding.EncodeToString(privRaw))
	d.Set("public", b64.StdEncoding.EncodeToString(pubRaw))
	d.Set("peerId", peerId.String())

	return keyRead(d, m)
}

func keyRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func keyUpdate(d *schema.ResourceData, m interface{}) error {
	return keyRead(d, m)
}

func keyDelete(d *schema.ResourceData, m interface{}) error {
	return keyRead(d, m)
}
