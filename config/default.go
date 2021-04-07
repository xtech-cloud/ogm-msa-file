package config

const defaultYAML string = `
service:
  name: xtc.api.ogm.file
  address: :9608
  ttl: 15
  interval: 10
logger:
  level: info
  dir: /var/log/ogm/
database:
  lite: true
  timeout: 10
  mysql:
    address: 127.0.0.1:3306
    user: root
    password: mysql@OMO
    db: ogm_file
  sqlite:
    path: /tmp/ogm-file.db
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
