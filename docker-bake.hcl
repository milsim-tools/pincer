target "docker-metadata-action" {}
target "default" {
  inherits   = ["docker-metadata-action"]
  context    = "./"
  dockerfile = "Dockerfile"
  args       = { VERSION = VERSION }

  contexts = {
    go    = "docker-image://golang:1.24.4"
    rocky = "docker-image://rockylinux/rockylinux:9-ubi-micro"
  }
}

variable "VERSION" {
  type = string
}
