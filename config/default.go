package config

const defaultYAML string = `
service: 
    name: xtc.ogm.file
    address: :18808
    ttl: 15
    interval: 10
logger:
    level: info
    dir: /var/log/ogm/
database:
    # 驱动类型，可选值为 [sqlite,mysql]
    driver: sqlite
    mysql:
        address: localhost:3306
        user: root
        password: mysql@XTC
        db: ogm
        # Set the maximum number of connections in the idle connection pool.
        # If MaxOpenConns is greater than 0 but less than the new MaxIdleConns, then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
        # If n <= 0, no idle connections are retained.
        maxIdleConns: 10
        # Sets the maximum number of open connections to the database.
        # If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit.
        # If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
        maxOpenConns: 100
        # Set the maximum amount of time a connection may be reused(Minute).
        # Expired connections may be closed lazily before reuse.
        # If d <= 0, connections are not closed due to a connection's age.
        maxLiftTime: 60 
        # Set the maximum amount of time a connection may be idle(Minute).
        # Expired connections may be closed lazily before reuse.
        # If d <= 0, connections are not closed due to a connection's idle time.
        maxIdleTime: 60 
    sqlite:
        path: /tmp/ogm-file.db
`
