package engine

import "github.com/OVINC-CN/IPCity/ipcity"

var (
	IPCityClient *ipcity.Client
	err          error
)

func InitIPCity() {
	IPCityClient = ipcity.NewClient()
	err = IPCityClient.Load("data/ipv4.dat")
	if err != nil {
		panic(err.Error())
	}
	err = IPCityClient.Load("data/ipv6.dat")
	if err != nil {
		panic(err.Error())
	}
}
