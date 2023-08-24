package provider

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

var (
	// DataMagicNumber defines the magic number in the ipcity date file header.
	DataMagicNumber = []byte("ipCT")
)

// DataVersion is the data version type.
type DataVersion byte

const (
	// DataVersionUnknown is the unknown ipcity date version.
	DataVersionUnknown = DataVersion(0)
	// DataVersionLatest is latest version of the ipcity data.
	DataVersionLatest = DataVersion(3)
)

// DataMode is the data mode type.
type DataMode byte

const (
	// DataModeUnknown is the unknown data mode.
	DataModeUnknown = DataMode(0)
	// DataModeIPv4 means this file is a IPv4 data file.
	DataModeIPv4 = DataMode(1)
	// DataModeIPv6 means this file is a IPv6 data file.
	DataModeIPv6 = DataMode(2)
)

// DataModeName is a mapping for data mode name.
var DataModeName = map[DataMode]string{
	DataModeUnknown: "Unknown",
	DataModeIPv4:    "IPv4",
	DataModeIPv6:    "IPv6",
}

var (
	headerBytesLength = 24
)

type headerImpl struct {
	Version           DataVersion
	Mode              DataMode
	IPIndexSize       byte
	MetaRowIndexSize  byte
	MetaRowCount      uint32
	EntityCount       uint32
	SourceUpdatedTime uint32
	UpdatedTime       uint32
}

func (h *headerImpl) ReadFrom(r io.Reader) error {
	if h == nil {
		return fmt.Errorf("init <nil> header")
	}

	buffer := make([]byte, headerBytesLength)
	if _, err := io.ReadFull(r, buffer); err != nil {
		return err
	}
	if mn := buffer[0:4]; bytes.Compare(DataMagicNumber, mn) != 0 {
		return fmt.Errorf("unrecognized magic number %X", string(mn))
	}
	h.Version = DataVersion(buffer[4])
	h.Mode = DataMode(buffer[5])
	h.IPIndexSize = buffer[6]
	h.MetaRowIndexSize = buffer[7]
	h.MetaRowCount = binary.BigEndian.Uint32(buffer[8:12])
	h.EntityCount = binary.BigEndian.Uint32(buffer[12:16])
	h.SourceUpdatedTime = binary.BigEndian.Uint32(buffer[16:20])
	h.UpdatedTime = binary.BigEndian.Uint32(buffer[20:])
	return nil
}

func (h *headerImpl) WriteTo(w io.Writer) (int64, error) {
	if h == nil {
		return 0, fmt.Errorf("dump <nil> header")
	}

	buffer := bytes.NewBuffer(make([]byte, 0, headerBytesLength))

	if err := func(processors ...func() error) error {
		for _, processor := range processors {
			if err := processor(); err != nil {
				return err
			}
		}
		return nil
	}(
		func() (err error) {
			_, err = buffer.Write(DataMagicNumber)
			return err
		},
		func() error { return buffer.WriteByte(byte(h.Version)) },
		func() error { return buffer.WriteByte(byte(h.Mode)) },
		func() error { return buffer.WriteByte(h.IPIndexSize) },
		func() error { return buffer.WriteByte(h.MetaRowIndexSize) },
		func() error { return binary.Write(buffer, binary.BigEndian, h.MetaRowCount) },
		func() error { return binary.Write(buffer, binary.BigEndian, h.EntityCount) },
		func() error { return binary.Write(buffer, binary.BigEndian, h.SourceUpdatedTime) },
		func() error { return binary.Write(buffer, binary.BigEndian, h.UpdatedTime) },
	); err != nil {
		return 0, err
	}

	return buffer.WriteTo(w)
}

// Header defines the ipcity data file header.
type Header struct {
	impl *headerImpl
}

// NewHeader returns a new data header with version and mode.
func NewHeader(version DataVersion, mode DataMode) *Header {
	return &Header{impl: &headerImpl{Version: version, Mode: mode}}
}

// WithMetaRowCount returns the data header with city count.
func (h *Header) WithMetaRowCount(metaRowCount uint32) *Header {
	if h != nil && h.impl != nil {
		h.impl.MetaRowCount = metaRowCount
	}
	return h
}

// WithEntityCount returns the data header with entity count.
func (h *Header) WithEntityCount(entityCount uint32) *Header {
	if h != nil && h.impl != nil {
		h.impl.EntityCount = entityCount
	}
	return h
}

// WithSourceUpdatedTime returns the data header with source updated time.
func (h *Header) WithSourceUpdatedTime(updatedTime int64) *Header {
	if h != nil && h.impl != nil && updatedTime >= 0 {
		h.impl.SourceUpdatedTime = uint32(updatedTime)
	}
	return h
}

// WithUpdatedTime returns the data header with updated time.
func (h *Header) WithUpdatedTime(updatedTime int64) *Header {
	if h != nil && h.impl != nil && updatedTime >= 0 {
		h.impl.UpdatedTime = uint32(updatedTime)
	}
	return h
}

// SetMetaRowCount set city count of the data header.
func (h *Header) SetMetaRowCount(metaRowCount uint32) {
	if h != nil && h.impl != nil {
		h.impl.MetaRowCount = metaRowCount
	}
}

// SetEntityCount set entity count of the data header.
func (h *Header) SetEntityCount(entityCount uint32) {
	if h != nil && h.impl != nil {
		h.impl.EntityCount = entityCount
	}
}

// Version returns the version of the ipcity data.
func (h *Header) Version() DataVersion {
	if h != nil && h.impl != nil {
		return h.impl.Version
	}
	return DataVersionUnknown
}

// Mode returns the mode of the ipcity data.
func (h *Header) Mode() DataMode {
	if h != nil && h.impl != nil {
		return h.impl.Mode
	}
	return DataModeUnknown
}

// ModeName returns the name of mode of the ipcity data.
func (h *Header) ModeName() string {
	if h != nil && h.impl != nil {
		if modeName, ok := DataModeName[h.impl.Mode]; ok {
			return modeName
		}
	}
	return DataModeName[DataModeUnknown]
}

// MetaRowCount returns the city count in the ipcity data.
func (h *Header) MetaRowCount() uint32 {
	if h != nil && h.impl != nil {
		return h.impl.MetaRowCount
	}
	return 0
}

// EntityCount returns the entity count in the ipcity data.
func (h *Header) EntityCount() uint32 {
	if h != nil && h.impl != nil {
		return h.impl.EntityCount
	}
	return 0
}

var (
	defaultIPIndexSizeMap = map[DataVersion]map[DataMode]uint32{
		DataVersion(3): map[DataMode]uint32{
			DataModeIPv4: 4,
			DataModeIPv6: 8,
		},
	}
)

// IPIndexSize returns the length of the IP index of a entity.
func (h *Header) IPIndexSize() uint32 {
	if h != nil && h.impl != nil {
		if h.impl.IPIndexSize > 0 {
			return uint32(h.impl.IPIndexSize)
		}
		// default IPIndex size
		if m, ok := defaultIPIndexSizeMap[h.Version()]; ok {
			if ipIndex, ok := m[h.Mode()]; ok {
				return ipIndex
			}
		}
	}
	return 0
}

func belong(value, min, max uint32) bool {
	return value >= min && value <= max
}

// MetaRowIndexSize returns the length of the meta row index of a entity.
func (h *Header) MetaRowIndexSize() uint32 {
	if h != nil && h.impl != nil {
		if h.impl.MetaRowIndexSize > 0 {
			return uint32(h.impl.MetaRowIndexSize)
		}

		// default MetaRowIndex size
		switch h.Version() {
		case DataVersion(3):
			switch {
			case belong(h.MetaRowCount(), 0, 0x000000FF):
				return 1
			case belong(h.MetaRowCount(), 0x00000100, 0x0000FFFF):
				return 2
			case belong(h.MetaRowCount(), 0x00010000, 0x00FFFFFF):
				return 3
			case belong(h.MetaRowCount(), 0x01000000, 0xFFFFFFFF):
				return 4
			default:
				return 0
			}
		default:
			return 0
		}
	}
	return 0
}

// SourceUpdatedTime returns the updated time of the ipcity source data.
func (h *Header) SourceUpdatedTime() time.Time {
	if h != nil && h.impl != nil {
		return time.Unix(int64(h.impl.SourceUpdatedTime), 0)
	}
	return time.Unix(0, 0)
}

// UpdatedTime returns the updated time of the ipcity data.
func (h *Header) UpdatedTime() time.Time {
	if h != nil && h.impl != nil {
		return time.Unix(int64(h.impl.UpdatedTime), 0)
	}
	return time.Unix(0, 0)
}

func (h *Header) String() string {
	var dummy string
	if h != nil {
		dummy = fmt.Sprintf(
			"{version:%d mode:%s "+
				"ipIndexSize:%d metaRowIndexSize:%d "+
				"metaRowCount:%d entityCount:%d "+
				"sourceUpdatedTime:%q updatedTime:%q}",
			h.Version(),
			h.ModeName(),
			h.IPIndexSize(),
			h.MetaRowIndexSize(),
			h.MetaRowCount(),
			h.EntityCount(),
			func() string {
				if h.SourceUpdatedTime().Unix() != 0 {
					return h.SourceUpdatedTime().Format("2006-01-02 15:04:05")
				}
				return ""
			}(),
			h.UpdatedTime().Format("2006-01-02 15:04:05"))
	}
	return dummy
}

// UnmarshalFrom will unmarshal Header from a reader.
func (h *Header) UnmarshalFrom(r io.Reader) error {
	if h != nil && h.impl != nil {
		return h.impl.ReadFrom(r)
	}
	return newNilParamError("Header")

}

// Unmarshal Header from a bytes buffer.
func (h *Header) Unmarshal(data []byte) error {
	return h.UnmarshalFrom(bytes.NewBuffer(data))
}

// MarshalTo will marshal Header to a writer.
func (h *Header) MarshalTo(w io.Writer) (int, error) {
	if h != nil && h.impl != nil {
		n, err := h.impl.WriteTo(w)
		return int(n), err
	}
	return 0, newNilParamError("Header")
}

// Marshal Header to a bytes buffer.
func (h *Header) Marshal() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, headerBytesLength))
	_, err := h.MarshalTo(buffer)
	return buffer.Bytes(), err
}
