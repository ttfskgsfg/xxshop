package redis_lib

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var Conn redis.Conn
var Err error
func init()  {

	Conn,Err = redis.Dial("tcp","192.168.0.105:6379")
	if Err != nil {
		fmt.Println(Err)
		panic(Err)
	}


}
