template_dir = "/deploy-template/nomad"
job {
  template = "job.nomadtpl"
}

env "dev" {
  git {
    default_ref = "refs/remotes/origin/develop"
  }
  docker {
    registry = "registry.hub.docker.com/library"
  }
  nomad {
    address = "10.30.83.2:4646"
    acl_token = ""
    connection {
      type = "direct"
    }
  }
}

