resource "maas_block_device" "vdb" {
  machine        = "abc123"
  name           = "vdb"
  id_path        = "/dev/vdb"
  size_gigabytes = 27
  tags = [
    "ssd",
  ]

  partitions {
    size_gigabytes = 10
    fs_type        = "ext4"
    label          = "media"
    mount_point    = "/media"
  }

  partitions {
    size_gigabytes = 15
    fs_type        = "ext4"
    mount_point    = "/storage"
  }
}
