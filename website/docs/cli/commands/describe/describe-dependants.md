---
title: atmos describe dependants
sidebar_label: dependants
sidebar_class_name: command
id: dependants
description: This command produces a list of Atmos components in Atmos stacks that depend on the provided Atmos component.
---

:::note Purpose
Use this command to show a list of Atmos components in Atmos stacks that depend on the provided Atmos component.
:::

## Description

In Atmos, you can define component dependencies by using the `settings.depends_on` section. The section used to define all the Atmos components (in
the same or different stacks) that the current component depends on.

The `settings.depends_on` section is a map of objects. The map keys are just the descriptions of dependencies and can be strings or numbers.
Provide meaningful descriptions so that people can understand what the dependencies are about.

Each object in the `settings.depends_on` section has the following schema:

- `component` (required) - an Atmos component that the current component depends on
- `namespace` (optional) - the `namespace` where the Atmos component is provisioned
- `tenant` (optional) - the `tenant` where the Atmos component is provisioned
- `environment` (optional) - the `environment` where the Atmos component is provisioned
- `stage` (optional) - the `stage` where the Atmos component is provisioned

<br/>

The `component` attribute is required. The rest are the context variables and are used to define Atmos stacks other than the current stack.
For example, you can specify:

- `namespace` if the `component` is from a different Organization
- `tenant` if the `component` is from a different Organizational Unit
- `environment` if the `component` is from a different region
- `stage` if the `component` is from a different account
- `tenant`, `environment` and `stage` if the component is from a different Atmos stack (e.g. `tenant1-ue2-dev`)

<br/>

In the following example, we define that the `top-level-component1` component depends on the following:

- The `test/test-component-override` component in the same Atmos stack
- The `test/test-component` component in Atmos stacks identified by the `dev` stage
- The `my-component` component from the `tenant1-ue2-staging` Atmos stack

```yaml title="examples/complete/stacks/catalog/terraform/top-level-component1.yaml"
components:
  terraform:
    top-level-component1:
      settings:
        depends_on:
          1:
            # If the `context` (namespace, tenant, environment, stage) is not provided, 
            # the `component` is from the same Atmos stack as this component
            component: "test/test-component-override"
          2:
            # This component (in any stage) depends on `test/test-component` 
            # from the `dev` stage (in any `environment` and any `tenant`)
            component: "test/test-component"
            stage: "dev"
          3:
            # This component depends on `my-component` 
            # from the `tenant1-ue2-staging` Atmos stack
            component: "my-component"
            tenant: "tenant1"
            environment: "ue2"
            stage: "staging"
      vars:
        enabled: true
```

In the following example, we specify that the `top-level-component2` component depends on the following:

- The `test/test-component` component in the same Atmos stack
- The `test/test2/test-component-2` component in the same Atmos stack

```yaml title="examples/complete/stacks/catalog/terraform/top-level-component2.yaml"
components:
  terraform:
    top-level-component2:
      metadata:
        # Point to Terraform component
        component: "top-level-component1"
      settings:
        depends_on:
          1:
            # If the `context` (namespace, tenant, environment, stage) is not provided, 
            # the `component` is from the same Atmos stack as this component
            component: "test/test-component"
          2:
            # If the `context` (namespace, tenant, environment, stage) is not provided, 
            # the `component` is from the same Atmos stack as this component
            component: "test/test2/test-component-2"
      vars:
        enabled: true
```

<br/>

Having the `top-level-component` and `top-level-component2` components configured as shown above, we can now execute the following Atmos command
to show all the components that depend on the `test/test-component` component in the `tenant1-ue2-dev` stack:

```shell
atmos describe dependants test/test-component -s tenant1-ue2-dev
```

```json
[
  {
    "component": "top-level-component1",
    "component_type": "terraform",
    "component_path": "examples/complete/components/terraform/top-level-component1",
    "namespace": "cp",
    "tenant": "tenant1",
    "environment": "ue2",
    "stage": "dev",
    "stack": "tenant1-ue2-dev",
    "stack_slug": "tenant1-ue2-dev-top-level-component1",
    "spacelift_stack": "tenant1-ue2-dev-top-level-component1",
    "atlantis_project": "tenant1-ue2-dev-top-level-component1"
  }
]
```

Similarly, the following Atmos command shows all the components that depend on the `test/test-component` component in
the `tenant1-ue2-test-1` stack:

```shell
atmos describe dependants test/test-component -s tenant1-ue2-test-1
```

```json
[
  {
    "component": "top-level-component1",
    "component_type": "terraform",
    "component_path": "examples/complete/components/terraform/top-level-component1",
    "namespace": "cp",
    "tenant": "tenant1",
    "environment": "ue2",
    "stage": "test-1",
    "stack": "tenant1-ue2-test-1",
    "stack_slug": "tenant1-ue2-dev-top-level-component1",
    "spacelift_stack": "tenant1-ue2-test-1-top-level-component1",
    "atlantis_project": "tenant1-ue2-test-1-top-level-component1"
  },
  {
    "component": "top-level-component2",
    "component_type": "terraform",
    "component_path": "examples/complete/components/terraform/top-level-component1",
    "namespace": "cp",
    "tenant": "tenant1",
    "environment": "ue2",
    "stage": "test-1",
    "stack": "tenant1-ue2-test-1",
    "stack_slug": "tenant1-ue2-test-1-top-level-component2",
    "atlantis_project": "tenant1-ue2-test-1-top-level-component2"
  }
]
```

<br/>

After the `test/test-component` has been provisioned, you can use the outputs to perform the following actions:

- Provision the dependent components by executing the Atmos commands `atmos terraform apply top-level-component1 -s tenant1-ue2-test-1` and
  `atmos terraform apply top-level-component2 -s tenant1-ue2-test-1` (on the command line or from a GitHub Action)

- Trigger the dependent Spacelift stack (from a GitHub Action by using the [spacectl](https://github.com/spacelift-io/spacectl) CLI, or by using an
  OPA [Trigger](https://docs.spacelift.io/concepts/policy/trigger-policy)
  policy, or by using
  the [spacelift_stack_dependency](https://registry.terraform.io/providers/spacelift-io/spacelift/latest/docs/resources/stack_dependency) resource)

- Trigger the dependent Atlantis project

## Usage

```shell
atmos describe dependants [options]
```

<br/>

:::tip
Run `atmos describe dependants --help` to see all the available options
:::

## Examples

```shell
atmos describe dependants test/test-component -s tenant1-ue2-test-1
atmos describe dependants test/test-component -s tenant1-ue2-dev --format yaml
atmos describe dependants test/test-component -s tenant1-ue2-test-1 -f yaml
atmos describe dependants test/test-component -s tenant1-ue2-test-1 --file dependants.json
atmos describe dependants test/test-component -s tenant1-ue2-test-1 --format yaml --file dependants.yaml
```

## Arguments

| Argument    | Description     | Required |
|:------------|:----------------|:---------|
| `component` | Atmos component | yes      |

## Flags

| Flag       | Description                                         | Alias | Required |
|:-----------|:----------------------------------------------------|:------|:---------|
| `--stack`  | Atmos stack                                         | `-s`  | yes      |
| `--format` | Output format: `json` or `yaml` (`json` is default) | `-f`  | no       |
| `--file`   | If specified, write the result to the file          |       | no       |

## Output

The command outputs a list of objects (in JSON or YAML format).

Each object has the following schema:

```json
{
  "component": "....",
  "component_type": "....",
  "component_path": "....",
  "namespace": "....",
  "tenant": "....",
  "environment": "....",
  "stage": "....",
  "stack": "....",
  "stack_slug": "",
  "spacelift_stack": ".....",
  "atlantis_project": "....."
}
```

where:

- `component` - the dependant Atmos component

- `component_type` - the type of the dependant component (`terraform` or `helmfile`)

- `component_path` - the filesystem path to the `terraform` or `helmfile` component

- `namespace` - the `namespace` where the dependant Atmos component is provisioned

- `tenant` - the `tenant` where the dependant Atmos component is provisioned

- `environment` - the `environment` where the dependant Atmos component is provisioned

- `stage` - the `stage` where the dependant Atmos component is provisioned

- `stack` - the Atmos stack where the dependant Atmos component is provisioned

- `stack_slug` - the Atmos stack slug (concatenation of the Atmos stack and Atmos component)

- `spacelift_stack` - the dependant Spacelift stack. It will be included only if the Spacelift workspace is enabled for the dependant Atmos component
  in the Atmos stack in the `settings.spacelift.workspace_enabled` section (either directly in the component's `settings.spacelift.workspace_enabled`
  section or via inheritance)

- `atlantis_project` - the dependant Atlantis project name. It will be included only if the Atlantis integration is configured in
  the `settings.atlantis` section in the stack config. Refer to [Atlantis Integration](/integrations/atlantis.md) for more details

<br/>

:::note

Abstract Atmos components (`metadata.type` is set to `abstract`) are not included in the output since they serve as blueprints for other
Atmos components and are not meant to be provisioned.

:::

## Output Example

```shell
atmos describe dependants test/test-component -s tenant1-ue2-test-1
```

```json
[
  {
    "component": "top-level-component2",
    "component_type": "terraform",
    "component_path": "examples/complete/components/terraform/top-level-component1",
    "namespace": "cp",
    "tenant": "tenant1",
    "environment": "ue2",
    "stage": "test-1",
    "stack": "tenant1-ue2-test-1",
    "stack_slug": "tenant1-ue2-dev-top-level-component2",
    "atlantis_project": "tenant1-ue2-test-1-top-level-component2"
  },
  {
    "component": "top-level-component1",
    "component_type": "terraform",
    "component_path": "examples/complete/components/terraform/top-level-component1",
    "namespace": "cp",
    "tenant": "tenant1",
    "environment": "ue2",
    "stage": "dev",
    "stack": "tenant1-ue2-dev",
    "stack_slug": "tenant1-ue2-test-1-top-level-component1",
    "spacelift_stack": "tenant1-ue2-dev-top-level-component1",
    "atlantis_project": "tenant1-ue2-dev-top-level-component1"
  }
]
```
