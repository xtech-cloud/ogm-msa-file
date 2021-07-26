package config

const defaultYAML string = `
service:
  name: xtc.ogm.file
  address: :18808
  ttl: 15
  interval: 10
logger:
  level: trace
  dir: /var/log/ogm/
database:
  driver: sqlite
  mysql:
    address: localhost:3306
    user: root
    password: mysql@XTC
    db: ogm
  sqlite:
    path: /tmp/ogm-file.db
`
