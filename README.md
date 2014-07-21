## bosh-diego-cpi (WIP)

This is an experimental BOSH Diego CPI. Currently it only supports stemcell and VM CPI actions work. 

### Resource pool cloud_properties

```
resource_pools:
- name: default
  network: default
  stemcell:
    name: bosh-diego
    version: latest
  cloud_properties:
    memory_mb: 10
    disk_mb: 1000 # Root disk size
```

### What's missing

- Networking: 
  there is no way to define "security group"

- Persistent disks: 
  Diego will schedule containers on any cell

- Health management via Diego: 
  Diego might move the container to a different cell and that will reset BOSH job to its empty VM state

- Ability to gracefully roll Diego cells: 
  when deploying Diego cells BOSH job VMs deployed with Diego CPI will be shutdown immediately
