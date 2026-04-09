package font

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
)

// IsWOFF reports whether data begins with the WOFF magic bytes.
// IsWOFF reports whether data is a WOFF1 font.
func IsWOFF(data []byte) bool {
	return len(data) >= 4 && data[0] == 0x77 && data[1] == 0x4F && data[2] == 0x46 && data[3] == 0x46
}

// IsWOFF2 reports whether data begins with the WOFF2 magic bytes.
// IsWOFF2 reports whether data is a WOFF2 font.
func IsWOFF2(data []byte) bool {
	return len(data) >= 4 && data[0] == 0x77 && data[1] == 0x4F && data[2] == 0x46 && data[3] == 0x32
}

const woffHeaderSize = 44
const woffTableEntrySize = 20

type woffHeader struct {
	Signature      uint32
	Flavor         uint32
	Length         uint32
	NumTables      uint16
	Reserved       uint16
	TotalSfntSize  uint32
	MajorVersion   uint16
	MinorVersion   uint16
	MetaOffset     uint32
	MetaLength     uint32
	MetaOrigLength uint32
	PrivOffset     uint32
	PrivLength     uint32
}

type woffTableEntry struct {
	Tag          uint32
	Offset       uint32
	CompLength   uint32
	OrigLength   uint32
	OrigChecksum uint32
}

// DecompressWOFF converts WOFF data to standard OTF/TTF format.
// DecompressWOFF converts WOFF1 data to raw TTF/OTF.
func DecompressWOFF(data []byte) ([]byte, error) {
	if !IsWOFF(data) {
		return nil, fmt.Errorf("not a WOFF file")
	}
	if len(data) < woffHeaderSize {
		return nil, fmt.Errorf("WOFF data too short for header")
	}

	r := bytes.NewReader(data)
	var hdr woffHeader
	if err := binary.Read(r, binary.BigEndian, &hdr); err != nil {
		return nil, fmt.Errorf("read WOFF header: %w", err)
	}

	numTables := int(hdr.NumTables)
	if len(data) < woffHeaderSize+numTables*woffTableEntrySize {
		return nil, fmt.Errorf("WOFF data too short for table directory")
	}

	entries := make([]woffTableEntry, numTables)
	for i := range entries {
		if err := binary.Read(r, binary.BigEndian, &entries[i]); err != nil {
			return nil, fmt.Errorf("read WOFF table entry %d: %w", i, err)
		}
	}

	// OTF/TTF header: 12 bytes offset table + 16 bytes per table record
	sfntHeaderSize := 12 + 16*numTables
	sfntOffset := uint32(sfntHeaderSize)

	// Align each table to 4-byte boundary and compute total size
	totalSize := uint32(sfntHeaderSize)
	tableOffsets := make([]uint32, numTables)
	for i, e := range entries {
		tableOffsets[i] = totalSize
		totalSize += e.OrigLength
		if totalSize%4 != 0 {
			totalSize += 4 - (totalSize % 4)
		}
	}
	_ = sfntOffset

	out := make([]byte, totalSize)

	// Write OTF/TTF offset table header
	binary.BigEndian.PutUint32(out[0:4], hdr.Flavor)
	binary.BigEndian.PutUint16(out[4:6], hdr.NumTables)

	searchRange := uint16(1)
	entrySelector := uint16(0)
	for searchRange*2 <= hdr.NumTables {
		searchRange *= 2
		entrySelector++
	}
	searchRange *= 16
	rangeShift := hdr.NumTables*16 - searchRange

	binary.BigEndian.PutUint16(out[6:8], searchRange)
	binary.BigEndian.PutUint16(out[8:10], entrySelector)
	binary.BigEndian.PutUint16(out[10:12], rangeShift)

	// Write table records and decompress table data
	for i, e := range entries {
		recOff := 12 + i*16
		binary.BigEndian.PutUint32(out[recOff:recOff+4], e.Tag)
		binary.BigEndian.PutUint32(out[recOff+4:recOff+8], e.OrigChecksum)
		binary.BigEndian.PutUint32(out[recOff+8:recOff+12], tableOffsets[i])
		binary.BigEndian.PutUint32(out[recOff+12:recOff+16], e.OrigLength)

		tableData := data[e.Offset : e.Offset+e.CompLength]

		if e.CompLength < e.OrigLength {
			fr := flate.NewReader(bytes.NewReader(tableData))
			decompressed, err := io.ReadAll(fr)
			fr.Close()
			if err != nil {
				return nil, fmt.Errorf("decompress WOFF table %d: %w", i, err)
			}
			if uint32(len(decompressed)) != e.OrigLength {
				return nil, fmt.Errorf("WOFF table %d: decompressed size %d != expected %d", i, len(decompressed), e.OrigLength)
			}
			copy(out[tableOffsets[i]:], decompressed)
		} else {
			copy(out[tableOffsets[i]:], tableData[:e.OrigLength])
		}
	}

	return out, nil
}
