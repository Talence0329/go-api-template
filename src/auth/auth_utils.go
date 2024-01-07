package auth

import (
	"math/rand"
	"strconv"
	"time"
)

func genRandomNum(len int) string {
	t := time.Now().UnixNano()
	r := rand.New(rand.NewSource(t))
	ri := r.Int()
	ret := ""
	for i := 0; i < len; i++ {
		ret = ret + strconv.Itoa(ri%10)
		ri = ri / 10
	}

	return ret
}
