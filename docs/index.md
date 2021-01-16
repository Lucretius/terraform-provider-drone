# Drone Provider

A [Terraform](https://www.terraform.io) provider for configuring the 
[Drone](https://drone.io) continuous delivery platform.

## Example Usage

```hcl
provider "drone" {
  server = "https://ci.example.com/"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoib2N0b2NhdCIsInR5cGUiOiJ1c2VyIn0.Fg0eYxO9x2CfGIvIHDZKhQbCGbRAsSB_iRDJlDEW6vc"
}
```

## Argument reference

* `server` - (Optional) The Drone servers url, It must be provided, but can also be sourced from the `DRONE_SERVER` environment variable.

* `token` - (Optional) The Drone servers api token, It must be provided, but can also be sourced from the `DRONE_TOKEN` environment variable.
