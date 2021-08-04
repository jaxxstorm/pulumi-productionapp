import * as pulumi from "@pulumi/pulumi";
import * as prodapp from "@jaxxstorm/pulumi-productionapp";

const app = new prodapp.Deployment("example", {
    image: "gcr.io/kuar-demo/kuard-amd64:blue",
    port: 80,
})

export const url = app.url
