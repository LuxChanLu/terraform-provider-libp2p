package internal

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type keyResource struct{}

var keyTypes = []string{"RSA", "ED25519", "SECP256K1", "ECDSA"}

var keyTypeMapping = map[string]int{
	keyTypes[0]: crypto.RSA,
	keyTypes[1]: crypto.Ed25519,
	keyTypes[2]: crypto.Secp256k1,
	keyTypes[3]: crypto.ECDSA,
}

type keyModel struct {
	ID      types.String `tfsdk:"id"`
	Type    types.String `tfsdk:"type"`
	Bits    types.Int64  `tfsdk:"bits"`
	Private types.String `tfsdk:"private"`
	Public  types.String `tfsdk:"public"`
	PeerId  types.String `tfsdk:"peer_id"`
}

func NewKeyResource() resource.Resource {
	return &keyResource{}
}

func (r *keyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (r *keyResource) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
				Description: "Unique ID of the key (CID)",
			},
			"type": {
				Type:        types.StringType,
				Required:    true,
				Description: "Key type",
				Validators:  []tfsdk.AttributeValidator{stringvalidator.OneOf(keyTypes...)},
			},
			"bits": {
				Type:        types.Int64Type,
				Optional:    true,
				Description: "Bits count for RSA key",
				Validators:  []tfsdk.AttributeValidator{int64validator.AtLeast(1)},
			},
			"private": {
				Type:      types.StringType,
				Sensitive: true,
				Computed:  true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
				Description: "Private key encoded in base64",
			},
			"public": {
				Type:      types.StringType,
				Sensitive: true,
				Computed:  true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
				Description: "Public key encoded in base64",
			},
			"peer_id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
				Description: "PeerID encoded in base58",
			},
		},
	}, nil
}

func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	newState := keyModel{}
	res.Diagnostics.Append(req.Plan.Get(ctx, &newState)...)
	if res.Diagnostics.HasError() {
		return
	}

	priv, pub, err := crypto.GenerateKeyPair(keyTypeMapping[newState.Type.Value], int(newState.Bits.Value))
	if err != nil {
		res.Diagnostics.AddError("Unable to generate key", err.Error())
		return
	}
	peerId, err := peer.IDFromPublicKey(pub)
	if err != nil {
		res.Diagnostics.AddError("Unable to get peer ID from public key", err.Error())
		return
	}
	privRaw, err := priv.Raw()
	if err != nil {
		res.Diagnostics.AddError("Unable to get raw private key", err.Error())
		return
	}
	pubRaw, err := pub.Raw()
	if err != nil {
		res.Diagnostics.AddError("Unable to get raw public key", err.Error())
		return
	}

	newState.ID = types.String{Value: peer.ToCid(peerId).String()}

	newState.Private = types.String{Value: base64.StdEncoding.EncodeToString(privRaw)}
	newState.Public = types.String{Value: base64.StdEncoding.EncodeToString(pubRaw)}
	newState.PeerId = types.String{Value: peerId.String()}

	res.Diagnostics.Append(res.State.Set(ctx, newState)...)
}

func (r *keyResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {

}

func (r *keyResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	model := &keyModel{}
	res.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if res.Diagnostics.HasError() {
		return
	}
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}

func (r *keyResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {

}
