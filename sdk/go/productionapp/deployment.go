// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package productionapp

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Deployment struct {
	pulumi.ResourceState

	// The URL from the generated service
	Url pulumi.StringOutput `pulumi:"url"`
}

// NewDeployment registers a new resource with the given unique name, arguments, and options.
func NewDeployment(ctx *pulumi.Context,
	name string, args *DeploymentArgs, opts ...pulumi.ResourceOption) (*Deployment, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Image == nil {
		return nil, errors.New("invalid value for required argument 'Image'")
	}
	if args.Port == nil {
		return nil, errors.New("invalid value for required argument 'Port'")
	}
	var resource Deployment
	err := ctx.RegisterRemoteComponentResource("productionapp:index:Deployment", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type deploymentArgs struct {
	// The image to deploy in your production application
	Image string `pulumi:"image"`
	// The port your container listens on
	Port int `pulumi:"port"`
}

// The set of arguments for constructing a Deployment resource.
type DeploymentArgs struct {
	// The image to deploy in your production application
	Image pulumi.StringInput
	// The port your container listens on
	Port pulumi.IntInput
}

func (DeploymentArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*deploymentArgs)(nil)).Elem()
}

type DeploymentInput interface {
	pulumi.Input

	ToDeploymentOutput() DeploymentOutput
	ToDeploymentOutputWithContext(ctx context.Context) DeploymentOutput
}

func (*Deployment) ElementType() reflect.Type {
	return reflect.TypeOf((*Deployment)(nil))
}

func (i *Deployment) ToDeploymentOutput() DeploymentOutput {
	return i.ToDeploymentOutputWithContext(context.Background())
}

func (i *Deployment) ToDeploymentOutputWithContext(ctx context.Context) DeploymentOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DeploymentOutput)
}

func (i *Deployment) ToDeploymentPtrOutput() DeploymentPtrOutput {
	return i.ToDeploymentPtrOutputWithContext(context.Background())
}

func (i *Deployment) ToDeploymentPtrOutputWithContext(ctx context.Context) DeploymentPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DeploymentPtrOutput)
}

type DeploymentPtrInput interface {
	pulumi.Input

	ToDeploymentPtrOutput() DeploymentPtrOutput
	ToDeploymentPtrOutputWithContext(ctx context.Context) DeploymentPtrOutput
}

type deploymentPtrType DeploymentArgs

func (*deploymentPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Deployment)(nil))
}

func (i *deploymentPtrType) ToDeploymentPtrOutput() DeploymentPtrOutput {
	return i.ToDeploymentPtrOutputWithContext(context.Background())
}

func (i *deploymentPtrType) ToDeploymentPtrOutputWithContext(ctx context.Context) DeploymentPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DeploymentPtrOutput)
}

// DeploymentArrayInput is an input type that accepts DeploymentArray and DeploymentArrayOutput values.
// You can construct a concrete instance of `DeploymentArrayInput` via:
//
//          DeploymentArray{ DeploymentArgs{...} }
type DeploymentArrayInput interface {
	pulumi.Input

	ToDeploymentArrayOutput() DeploymentArrayOutput
	ToDeploymentArrayOutputWithContext(context.Context) DeploymentArrayOutput
}

type DeploymentArray []DeploymentInput

func (DeploymentArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Deployment)(nil)).Elem()
}

func (i DeploymentArray) ToDeploymentArrayOutput() DeploymentArrayOutput {
	return i.ToDeploymentArrayOutputWithContext(context.Background())
}

func (i DeploymentArray) ToDeploymentArrayOutputWithContext(ctx context.Context) DeploymentArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DeploymentArrayOutput)
}

// DeploymentMapInput is an input type that accepts DeploymentMap and DeploymentMapOutput values.
// You can construct a concrete instance of `DeploymentMapInput` via:
//
//          DeploymentMap{ "key": DeploymentArgs{...} }
type DeploymentMapInput interface {
	pulumi.Input

	ToDeploymentMapOutput() DeploymentMapOutput
	ToDeploymentMapOutputWithContext(context.Context) DeploymentMapOutput
}

type DeploymentMap map[string]DeploymentInput

func (DeploymentMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Deployment)(nil)).Elem()
}

func (i DeploymentMap) ToDeploymentMapOutput() DeploymentMapOutput {
	return i.ToDeploymentMapOutputWithContext(context.Background())
}

func (i DeploymentMap) ToDeploymentMapOutputWithContext(ctx context.Context) DeploymentMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(DeploymentMapOutput)
}

type DeploymentOutput struct {
	*pulumi.OutputState
}

func (DeploymentOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Deployment)(nil))
}

func (o DeploymentOutput) ToDeploymentOutput() DeploymentOutput {
	return o
}

func (o DeploymentOutput) ToDeploymentOutputWithContext(ctx context.Context) DeploymentOutput {
	return o
}

func (o DeploymentOutput) ToDeploymentPtrOutput() DeploymentPtrOutput {
	return o.ToDeploymentPtrOutputWithContext(context.Background())
}

func (o DeploymentOutput) ToDeploymentPtrOutputWithContext(ctx context.Context) DeploymentPtrOutput {
	return o.ApplyT(func(v Deployment) *Deployment {
		return &v
	}).(DeploymentPtrOutput)
}

type DeploymentPtrOutput struct {
	*pulumi.OutputState
}

func (DeploymentPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Deployment)(nil))
}

func (o DeploymentPtrOutput) ToDeploymentPtrOutput() DeploymentPtrOutput {
	return o
}

func (o DeploymentPtrOutput) ToDeploymentPtrOutputWithContext(ctx context.Context) DeploymentPtrOutput {
	return o
}

type DeploymentArrayOutput struct{ *pulumi.OutputState }

func (DeploymentArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Deployment)(nil))
}

func (o DeploymentArrayOutput) ToDeploymentArrayOutput() DeploymentArrayOutput {
	return o
}

func (o DeploymentArrayOutput) ToDeploymentArrayOutputWithContext(ctx context.Context) DeploymentArrayOutput {
	return o
}

func (o DeploymentArrayOutput) Index(i pulumi.IntInput) DeploymentOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Deployment {
		return vs[0].([]Deployment)[vs[1].(int)]
	}).(DeploymentOutput)
}

type DeploymentMapOutput struct{ *pulumi.OutputState }

func (DeploymentMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]Deployment)(nil))
}

func (o DeploymentMapOutput) ToDeploymentMapOutput() DeploymentMapOutput {
	return o
}

func (o DeploymentMapOutput) ToDeploymentMapOutputWithContext(ctx context.Context) DeploymentMapOutput {
	return o
}

func (o DeploymentMapOutput) MapIndex(k pulumi.StringInput) DeploymentOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) Deployment {
		return vs[0].(map[string]Deployment)[vs[1].(string)]
	}).(DeploymentOutput)
}

func init() {
	pulumi.RegisterOutputType(DeploymentOutput{})
	pulumi.RegisterOutputType(DeploymentPtrOutput{})
	pulumi.RegisterOutputType(DeploymentArrayOutput{})
	pulumi.RegisterOutputType(DeploymentMapOutput{})
}
