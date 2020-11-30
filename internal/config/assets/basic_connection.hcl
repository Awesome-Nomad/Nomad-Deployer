template_dir = "./example"
git_project_dir = "./projects"
env "basic" {
  nomad {
    address = "localhost:4646"
    acl_token = ""
    connection {
      type = "direct"
    }
  }
  consul {
    address = "localhost:8500"
    acl_token = ""
    connection {
      type = "ssh"
      ssh {
        username = "root"
        key_file = "/root/.ssh/id_rsa"
      }
    }
  }
}