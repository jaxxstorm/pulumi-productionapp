using Pulumi;
using Pulumi.Productionapp;

class ProductionApp : Stack
{
    public ProductionApp()
    {
        var app = new Pulumi.Productionapp.Deployment("example", new Pulumi.Productionapp.DeploymentArgs
        {
            Image = "gcr.io/kuar-demo/kuard-amd64:blue",
            Port = 80
        });

        this.Url = app.Url;

    }


    [Output] public Output<string> Url { get; set; }

}
