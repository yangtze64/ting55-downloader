package request

import (
	"fmt"
	"math/rand"
	"time"
)

// GenIpaddr 生成IP
func GenIpaddr() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
