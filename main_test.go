package main

import (
	"fmt"
	"github.com/OVINC-CN/IPCity/engine"
	"math/rand"
	"testing"
)

const maxIP = 1000 * 1000

func BenchmarkIPCityV4(b *testing.B) {
	var ipList [maxIP]string
	for i := int64(0); i < maxIP; i++ {
		ipList[i] = randomIPV4()
	}
	engine.InitIPCity()
	var index int64 = 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ip := ipList[index%maxIP]
			engine.IPCityClient.Search(ip)
			index++
		}
	})
}

func BenchmarkIPCityV6(b *testing.B) {
	var ipList [maxIP]string
	for i := int64(0); i < maxIP; i++ {
		ipList[i] = randomIPV6()
	}
	engine.InitIPCity()
	var index int64 = 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ip := ipList[index%maxIP]
			engine.IPCityClient.Search(ip)
			index++
		}
	})
}

func randomIPV4() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func randomIPV6() string {
	return fmt.Sprintf(
		"%s%s%s%s:%s%s%s%s:%s%s%s%s:%s%s%s%s:%s%s%s%s:%s%s%s%s:%s%s%s%s:%s%s%s%s",
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
		getRandomChar(16),
	)
}

const randomData = "0123456789abcdef"

func getRandomChar(n int) string {
	index := rand.Intn(n)
	return randomData[index : index+1]
}
