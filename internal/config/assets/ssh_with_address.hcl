env "ssh_with_address" {
  consul {
    address = "localhost:8500"
    acl_token = ""
    connection {
      type = "ssh"
      ssh {
        username = "root"
        key_file = "/root/.ssh/id_rsa"
        address = "localhost:22"
      }
    }
  }
}