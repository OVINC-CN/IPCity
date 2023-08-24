package ipcity

import (
	"github.com/OVINC-CN/IPCity/ipcity/provider"
	"net"
	"os"
	"sync"
)

var (
	stores = make([]*Store, 0, 16)
	mutex  sync.Mutex
)

// Store exports provider.Store.
type Store = provider.Store

// Header exports provider.Header.
type Header = provider.Header

// Meta exports provider.Meta.
type Meta = provider.Meta

// Entity exports provider.Entity.
type Entity = provider.Entity

// Load IPCity data file.
func Load(filename string) error {
	var err error

	// open file by filename
	var fileReader *os.File
	if fileReader, err = os.Open(filename); err != nil {
		return err
	}
	defer func() { _ = fileReader.Close() }()

	// load store from file reader
	store := &Store{}
	if err = store.UnmarshalFrom(fileReader); err != nil {
		return err
	}

	// append store list
	func() {
		mutex.Lock()
		defer mutex.Unlock()
		stores = append(stores, store)
	}()

	return nil
}

// Search meta by address .
func Search(addr string) *Meta {
	var meta *Meta
	for _, v := range stores {
		meta = v.Search(net.ParseIP(addr))
		if meta != nil && !meta.IsEmpty() {
			break
		}
	}
	return meta
}

// ClientInterface 用于查询ip归属地信息的接口
type ClientInterface interface {
	Load(filename string) error
	Search(addr string) *Meta
}

type Client struct {
	stores []*Store
}

// Load 加载ip信息库
func (c *Client) Load(filename string) error {
	var err error
	// open file by filename
	var fileReader *os.File
	if fileReader, err = os.Open(filename); err != nil {
		return err
	}
	defer func() { _ = fileReader.Close() }()
	// load store from file reader
	store := &Store{}
	if err = store.UnmarshalFrom(fileReader); err != nil {
		return err
	}
	// append store list
	c.stores = append(c.stores, store)
	return nil
}

// Search 查询ip信息
func (c *Client) Search(addr string) *Meta {
	var meta *Meta
	for _, v := range c.stores {
		meta = v.Search(net.ParseIP(addr))
		if meta != nil && !meta.IsEmpty() {
			break
		}
	}
	return meta
}

// NewClient 生成client对象
func NewClient() *Client {
	return &Client{}
}
