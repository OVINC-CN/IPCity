package provider

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
)

var (
	newNilParamError = func(paramName string) error {
		return fmt.Errorf("work with nil %s", paramName)
	}
	newUnsupportedVersionError = func(version DataVersion) error {
		return fmt.Errorf("work with the unsupported version %d", version)
	}
)

var (
	goUntilError = func(fn ...func() error) error {
		for _, f := range fn {
			if err := f(); err != nil {
				return err
			}
		}
		return nil
	}
)

// Store defines a store stored the ipcity data.
type Store struct {
	header     *Header
	metaTable  []*Meta
	entityList []*Entity
}

// NewStore returns a new store.
func NewStore() *Store {
	return &Store{}
}

// WithHeader returns the store with meta table.
func (s *Store) WithHeader(header *Header) *Store {
	if s != nil {
		s.header = header
	}
	return s
}

// WithMetaTable returns the store with meta table.
func (s *Store) WithMetaTable(metaTable []*Meta) *Store {
	if s != nil {
		s.metaTable = metaTable
	}
	return s
}

// WithEntityList returns the store with entity list.
func (s *Store) WithEntityList(entityList []*Entity) *Store {
	if s != nil {
		s.entityList = entityList
	}
	return s
}

// Header returns the header of the store.
func (s *Store) Header() *Header {
	if s != nil {
		return s.header
	}
	return nil
}

// MetaTable returns the meta table of the store.
func (s *Store) MetaTable() []*Meta {
	if s != nil {
		return s.metaTable
	}
	return nil
}

// Meta returns the pointed index meta in the meta table.
func (s *Store) Meta(i int) *Meta {
	if s != nil && i < s.MetaRowCount() && i >= 0 {
		return s.metaTable[i]
	}
	return nil
}

// MetaRowCount returns the row count of the meta table.
func (s *Store) MetaRowCount() int {
	if s != nil {
		return len(s.metaTable)
	}
	return 0
}

// EntityList returns the entity list of the store.
func (s *Store) EntityList() []*Entity {
	if s != nil {
		return s.entityList
	}
	return nil
}

// Entity returns the pointed index entity in the entity list.
func (s *Store) Entity(i int) *Entity {
	if s != nil && i < s.EntityCount() && i >= 0 {
		return s.entityList[i]
	}
	return nil
}

// EntityCount return the length of the entity list.
func (s *Store) EntityCount() int {
	if s != nil {
		return len(s.entityList)
	}
	return 0
}

func (s *Store) searchByIPIndex(ipIndex uint64) *Meta {
	if index := sort.Search(s.EntityCount(), func(i int) bool {
		return s.Entity(i).IPIndex() >= ipIndex
	}); s.Entity(index) != nil {
		if s.Entity(index).IPIndex() != ipIndex {
			index = index - 1
		}
		return s.Meta(int(s.Entity(index).MetaRowIndex()))
	}
	return nil
}

// Search returns the meta queryed from the store.
func (s *Store) Search(addr net.IP) *Meta {
	ipIndex := uint64(0)
	switch s.Header().Mode() {
	case DataModeIPv4:
		if b := []byte(addr.To4()); b != nil {
			if s.Header().Version() == DataVersion(2) {
				ipIndex = uint64(binary.BigEndian.Uint32(b) >> 8)
			} else {
				ipIndex = uint64(binary.BigEndian.Uint32(b))
			}
		}
	case DataModeIPv6:
		if b := []byte(addr.To16()); b != nil {
			ipIndex = binary.BigEndian.Uint64(b[0:8])
		}
	default:
		// pass
	}
	return s.searchByIPIndex(ipIndex)
}

// UnmarshalFrom will unmarshal Store from a raeder.
func (s *Store) UnmarshalFrom(reader io.Reader) error {
	ireader := bufio.NewReader(reader)
	err := goUntilError(func() error {
		header := &Header{impl: &headerImpl{}}
		if err := header.UnmarshalFrom(ireader); err != nil {
			return fmt.Errorf("unmarshal header error, %s", err)
		}
		s.header = header
		return nil
	}, func() error {
		var err error
		metaTable := make([]*Meta, 0, s.Header().MetaRowCount())
		var i int
		for i = 0; i < int(s.Header().MetaRowCount()); i++ {
			var line []byte
			if line, err = ireader.ReadBytes('\n'); err != nil {
				break
			}
			meta := &Meta{}
			if err = meta.Unmarshal(line); err != nil {
				break
			}
			metaTable = append(metaTable, meta)
		}
		if err != nil {
			return fmt.Errorf("unmarshal meta table row[%d/%d] error, %s",
				i, s.Header().MetaRowCount(), err)
		}
		s.metaTable = metaTable
		return nil
	}, func() error {
		if s.MetaRowCount() != int(s.Header().MetaRowCount()) {
			return fmt.Errorf(
				"invalid meta table size %d, header.MetaRowCount is %d",
				s.MetaRowCount(), s.Header().MetaRowCount())
		}
		return nil
	}, func() error {
		unmarshaler := &EntityUnmarshaler{
			IPIndexSize:      s.Header().IPIndexSize(),
			MetaRowIndexSize: s.Header().MetaRowIndexSize(),
		}

		var entityList []*Entity
		if s.Header().EntityCount() > 0 {
			entityList = make([]*Entity, 0, s.Header().EntityCount())
		} else {
			entityList = make([]*Entity, 0, 1024)
		}

		var err error
		for {
			entity := &Entity{}
			if err = unmarshaler.UnmarshalFrom(ireader, entity); err != nil {
				break
			}
			entityList = append(entityList, entity)
		}
		if err != nil && err != io.EOF {
			return fmt.Errorf("unmarshal entity list error, %s", err)
		}
		s.entityList = entityList
		return nil
	}, func() error {
		if s.Header().EntityCount() > 0 &&
			s.EntityCount() != int(s.Header().EntityCount()) {
			return fmt.Errorf(
				"invalid entity list size %d, header.EntityCount is %d",
				s.EntityCount(), s.Header().EntityCount())
		}
		return nil
	})
	return err
}

// MarshalTo will marshal Store to a writer.
func (s *Store) MarshalTo(writer io.Writer) (int, error) {
	return 0, fmt.Errorf("unsupported interface! %s", writer)
}
