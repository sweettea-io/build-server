package redis

import r "github.com/gomodule/redigo/redis"

// NewPool creates and returns a configured Redis pool instance.
func NewPool(address string, password string, maxActive int, maxIdle int, wait bool) *r.Pool {
  pool := &r.Pool{
    // Max number of connections allocated by the pool at any given time (0 = unlimited)
    MaxActive: maxActive,

    // Max number of idle connections in the pool.
    MaxIdle: maxIdle,

    // Whether to wait for newly available connections if the pool is at its MaxActive limit.
    Wait: wait,

    Dial: func() (r.Conn, error) {
      // Connect to Redis address (hostname:port)
      conn, err := r.Dial("tcp", address)

      if err != nil {
        panic(err)
      }

      // If Redis password is provided, use that for auth.
      if password != "" {
        conn.Do("AUTH", password)
      }

      return conn, nil
    },
  }

  return pool
}