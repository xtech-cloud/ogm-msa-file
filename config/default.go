package config

const defaultYAML string = `
service:
  name: omo.msa.file
  address: :9608
  ttl: 15
  interval: 10
logger:
  level: info
  dir: /var/log/msa/
database:
  lite: true
  timeout: 10
  mysql:
    address: 127.0.0.1:3306
    user: root
    password: mysql@OMO
    db: msa_file
  sqlite:
    path: /tmp/msa-file.db
publisher:
- /bucket/make
- /bucket/updateengine
- /bucket/updatecapacity
- /bucket/resettoken
- /bucket/remove
- /object/prepare
- /object/flush
- /object/remove
`
