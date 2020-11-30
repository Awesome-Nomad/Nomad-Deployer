# Template directory structure
/**
deployment
├── dev.yaml
├── base.yaml
├── prod.yaml
├── job.nomadtpl
└── stg.yaml
*/

template_dir = "./deploy-template/nomad"
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
      type = "ssh"
      ssh {
        username = "root"
        key_file = "~/.ssh/id_rsa"
        address = "10.30.83.75"
      }
    }
  }
}

env "stg" {
  git {
    default_ref = "refs/remotes/origin/master"
  }
  docker {
    registry = "registry.hub.docker.com/library"
  }
  nomad {
    address = "10.30.83.75:4646"
    acl_token = ""
    connection {
      type = "ssh"
      ssh {
        username = "root"
        key_file = "~/.ssh/id_rsa"
        address = "10.30.83.75"
      }
    }
  }
}

env "prod" {
  docker {
    registry = "registry.hub.docker.com/library"
  }
  nomad {
    address = "10.30.83.66:4646"
    acl_token = ""
    connection {
      type = "ssh"
      ssh {
        username = "root"
        key_file = "~/.ssh/id_rsa"
        address = "10.30.83.75"
      }
    }
  }
}
