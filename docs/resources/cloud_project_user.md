---
subcategory : "Cloud Project"
---

# ovh_cloud_project_user

Creates a user in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_user" "user1" {
  service_name = "XXX"
}
```

```terraform
# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_user" "user_with_rotation" {
  service_name   = "XXX"
  description    = "Service User created by Terraform with password rotation"
  password_reset = "2025-04-30"
}
```

-> **NOTE** To reset password of the user previously created, update the `password_reset` attribute. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password. This attribute can be an arbitrary string but we recommend 2 formats:
- a datetime to keep a trace of the last reset
- a md5 of other variables to automatically trigger it based on this variable update
## Argument Reference

The following arguments are supported:

* `description` - A description associated with the user.

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `role_name` - The name of a role. See `role_names`.

* `role_names` - A list of role names. Values can be:
  - admin
  - administrator
  - ai_training_operator
  - ai_training_read
  - authentication
  - backup_operator
  - compute_operator
  - image_operator
  - infrastructure_supervisor
  - key-manager_operator
  - key-manager_read
  - load-balancer_operator
  - network_operator
  - network_security_operator
  - objectstore_operator
  - volume_operator

* `password_reset` - (Optional) Arbitrary string to change to trigger a password update. Use the `terraform refresh` command after executing `terraform apply` to update the output with the new password.

## Attributes Reference

The following attributes are exported:

* `creation_date` - the date the user was created.
* `description` - See Argument Reference above.
* `openstack_rc` - a convenient map representing an openstack_rc file. Note: no password nor sensitive token is set in this map.
* `password` - (Sensitive) the password generated for the user. The password can be used with the Openstack API. This attribute is sensitive and will only be retrieve once during creation.
* `password_reset` - See Argument Reference above.
* `roles` - A list of roles associated with the user.
  * `description` - description of the role
  * `id` - id of the role
  * `name` - name of the role
  * `permissions` - list of permissions associated with the role
* `service_name` - See Argument Reference above.
* `status` - the status of the user. should be normally set to 'ok'.
* `username` - the username generated for the user. This username can be used with the Openstack API.
