# How to test this provider
1. Start the provider via `go run main.go -debug`
2. Create a "test" folder
3. Create a dev.tfrc file inside this folder and adjust the path of the following snippet accordingly to your system (terraform-provider-gosoline is the repo root)
```hcl
provider_installation {
  dev_overrides {
    "justtrackio/gosoline" = "/home/username/projects/go/oss/terraform-provider-gosoline/bin/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
4. Create a main.tf with the needed config to call resources of the provider, e.g.:
```terraform
provider "gosoline" {
  metadata = {
    domain    = "my.zone"
    use_https = false
    port      = 1234
  }
  name_patterns = {
    hostname                         = "{scheme}://{app}.{group}.{env}.{metadata_domain}:{port}"
    cloudwatch_namespace             = "{project}/{env}/{family}/{group}-{app}"
    ecs_cluster                      = "{project}-{env}-{family}"
    ecs_service                      = "{group}-{app}"
    grafana_cloudwatch_datasource    = "cloudwatch-{family}"
    grafana_elasticsearch_datasource = "elasticsearch-{env}-logs-{project}-{family}-{group}-{app}"
    kubernetes_namespace             = ""
    kubernetes_pod                   = ""
    traefik_service_name             = ""
  }
}

data "gosoline_application_dashboard_definition" "test" {
  project     = "prj"
  environment = "env"
  family      = "fam"
  group       = "grp"
  application = "app"
  containers  = ["app", "log_router"]
}

output "dashboard" {
  value = data.gosoline_application_dashboard_definition.test.body
}

terraform {
  required_providers {
    gosoline = {
      source  = "justtrackio/gosoline"
      version = "1.1.0"
    }
  }
}
```
5. Export the vars needed for local testing and adjust the output of TF_REATTACH_PROVIDERS accordingly to the output of your `go run main.go -debug` output
```shell
export TF_CLI_CONFIG_FILE=/Users/sebastian/Projects/go/oss/terraform-provider-gosoline/test/dev.tfrc
export TF_REATTACH_PROVIDERS='{"registry.terraform.io/justtrackio/gosoline":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/your-custom-folder"}}}'
```
