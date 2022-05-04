package com.jaxxstorm.example.productionapp;

import com.pulumi.Pulumi;
import com.pulumi.productionapp.Deployment;
import com.pulumi.productionapp.DeploymentArgs;


public class App {

    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var app = new Deployment("example",
                    DeploymentArgs.builder()
                            .image("gcr.io/kuar-demo/kuard-amd64:blue")
                            .port(80)
                            .build());
            ctx.export("url", app.url());
        });
    }


}
