package provider

import (
	"fmt"
	"io"
)

// Entity defines a Entity of the ipcity data.
type Entity struct {
	ipIndex      uint64
	metaRowIndex uint32
}

// NewEntity returns a new entity with IP index and metaRow Index.
func NewEntity(pi uint64, mi uint32) *Entity {
	return &Entity{
		ipIndex:      pi,
		metaRowIndex: mi,
	}
}

// IPIndex returns the IP index in the entity as uint64.
func (e *Entity) IPIndex() uint64 {
	if e != nil {
		return e.ipIndex
	}
	return 0
}

// MetaRowIndex returns the meta row index in the entity.
func (e *Entity) MetaRowIndex() uint32 {
	if e != nil {
		return e.metaRowIndex
	}
	return 0
}

func readScalableBigEndianOrderBytesToUint32(bytes []byte, value *uint32) {
	l := len(bytes) - 1
	if l >= 0 && value != nil {
		for i := range bytes {
			*value |= uint32(bytes[i]) << (8 * (l - i))
		}
	}
}

func readScalableBigEndianOrderBytesToUint64(bytes []byte, value *uint64) {
	l := len(bytes) - 1
	if l >= 0 && value != nil {
		for i := range bytes {
			*value |= uint64(bytes[i]) << (8 * (l - i))
		}
	}
}

func writeUint32ToScalableBigEndianOrderBytes(value uint32, bytes []byte) {
	l := len(bytes) - 1
	for i := range bytes {
		bytes[i] = byte(value >> (8 * (l - i)))
	}
}

func writeUint64ToScalableBigEndianOrderBytes(value uint64, bytes []byte) {
	l := len(bytes) - 1
	for i := range bytes {
		bytes[i] = byte(value >> (8 * (l - i)))
	}
}

var (
	ipIndexSizeSelector = map[uint32]bool{
		4: true,
		8: true,
	}
	metaRowIndexSizeSelector = map[uint32]bool{
		1: true,
		2: true,
		3: true,
		4: true,
	}
)

func validateIPIndexSize(size uint32) bool {
	_, ok := ipIndexSizeSelector[size]
	return ok
}

func validateMetaRowIndexSize(size uint32) bool {
	_, ok := metaRowIndexSizeSelector[size]
	return ok
}

// EntityUnmarshaler define a unmarshaler for the entity.
type EntityUnmarshaler struct {
	IPIndexSize      uint32
	MetaRowIndexSize uint32
}

// UnmarshalFrom will unmarshal Entity from a reader.
func (u *EntityUnmarshaler) UnmarshalFrom(reader io.Reader, entity *Entity) error {
	if u == nil {
		return newNilParamError("EntityUnmarshaler")
	}
	if !validateIPIndexSize(u.IPIndexSize) || !validateMetaRowIndexSize(u.MetaRowIndexSize) {
		return fmt.Errorf("entity unmarshaler with invalid index size")
	}

	if reader == nil {
		return newNilParamError("Reader")
	}
	if entity == nil {
		return newNilParamError("Entity")
	}

	buffer := make([]byte, u.IPIndexSize+u.MetaRowIndexSize)
	n, err := io.ReadFull(reader, buffer)
	if n == 1 && buffer[0] == 0 {
		return io.EOF
	}
	if err != nil {
		return err
	}
	readScalableBigEndianOrderBytesToUint64(buffer[0:u.IPIndexSize], &(entity.ipIndex))
	readScalableBigEndianOrderBytesToUint32(buffer[u.IPIndexSize:], &(entity.metaRowIndex))
	return nil
}

// EntityMarshaler defines a marshaler for the entity.
type EntityMarshaler struct {
	DataVersion      DataVersion
	IPIndexSize      uint32
	MetaRowIndexSize uint32
}

// MarshalTo will marshal Entity to a bytes buffer.
func (m *EntityMarshaler) MarshalTo(entity *Entity, writer io.Writer) (int, error) {
	if m == nil {
		return 0, newNilParamError("EntityMarshaler")
	}
	if !validateIPIndexSize(m.IPIndexSize) || !validateMetaRowIndexSize(m.MetaRowIndexSize) {
		return 0, fmt.Errorf("entity marshaler with invalid index size")
	}

	if entity == nil {
		return 0, newNilParamError("Entity")
	}
	if writer == nil {
		return 0, newNilParamError("Writer")
	}

	bytes := make([]byte, m.IPIndexSize+m.MetaRowIndexSize)
	ipIndex := entity.ipIndex
	writeUint64ToScalableBigEndianOrderBytes(ipIndex, bytes[0:m.IPIndexSize])
	writeUint32ToScalableBigEndianOrderBytes(entity.metaRowIndex, bytes[m.IPIndexSize:])
	return writer.Write(bytes)
}
