package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/jaxxstorm/pulumi-productionapp/sdk/go/productionapp"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		application, err := productionapp.NewDeployment(ctx, "example", &productionapp.DeploymentArgs{
			Image: pulumi.String("gcr.io/kuar-demo/kuard-amd64:blue"),
			Port:  pulumi.Int(80),
		})
		if err != nil {
			return fmt.Errorf("error creating application: %v", err)
		}
	
		ctx.Export("url", application.Url)
	
		return nil
	})
}
