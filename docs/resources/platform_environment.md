---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "harness_platform_environment Resource - terraform-provider-harness"
subcategory: "Next Gen"
description: |-
  Resource for creating a Harness environment.
---

# harness_platform_environment (Resource)

Resource for creating a Harness environment.

## Example Usage

```terraform
resource "harness_platform_environment" "test" {
  identifier = "identifier"
  name       = "name"
  org_id     = "org_id"
  project_id = "project_id"
  tags       = ["foo:bar", "baz"]
  type       = "PreProduction"
  yaml       = <<-EOT
			   environment:
         name: name
         identifier: identifier
         orgIdentifier: org_id
         projectIdentifier: project_id
         type: PreProduction
         tags:
           foo: bar
           baz: ""
         variables:
           - name: envVar1
             type: String
             value: v1
             description: ""
           - name: envVar2
             type: String
             value: v2
             description: ""
         overrides:
           manifests:
             - manifest:
                 identifier: manifestEnv
                 type: Values
                 spec:
                   store:
                     type: Git
                     spec:
                       connectorRef: <+input>
                       gitFetchType: Branch
                       paths:
                         - file1
                       repoName: <+input>
                       branch: master
           configFiles:
             - configFile:
                 identifier: configFileEnv
                 spec:
                   store:
                     type: Harness
                     spec:
                       files:
                         - account:/Add-ons/svcOverrideTest
                       secretFiles: []
      EOT
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `identifier` (String) Unique identifier of the resource.
- `name` (String) Name of the resource.
- `org_id` (String) Unique identifier of the Organization.
- `project_id` (String) Unique identifier of the Project.
- `type` (String) The type of environment. Valid values are PreProduction, Production

### Optional

- `color` (String) Color of the environment.
- `description` (String) Description of the resource.
- `tags` (Set of String) Tags to associate with the resource. Tags should be in the form `name:value`.
- `yaml` (String) Environment YAML

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Import using environment id
terraform import harness_platform_environment.example <environment_id>
```
