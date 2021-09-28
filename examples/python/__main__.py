"""A Python Pulumi program"""

import pulumi
import jaxxstorm_pulumi_productionapp as prod_app

app = prod_app.Deployment("example", 
    image="gcr.io/kuar-demo/kuard-amd64:blue",
    port=80,
)

pulumi.export("url", app.url)