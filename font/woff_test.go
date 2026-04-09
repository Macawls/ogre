package font

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestIsWOFF(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"woff magic", []byte{0x77, 0x4F, 0x46, 0x46, 0x00}, true},
		{"woff2 magic", []byte{0x77, 0x4F, 0x46, 0x32, 0x00}, false},
		{"ttf", []byte{0x00, 0x01, 0x00, 0x00}, false},
		{"short", []byte{0x77, 0x4F}, false},
		{"empty", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWOFF(tt.data); got != tt.want {
				t.Errorf("IsWOFF = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWOFF2(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"woff2 magic", []byte{0x77, 0x4F, 0x46, 0x32, 0x00}, true},
		{"woff1 magic", []byte{0x77, 0x4F, 0x46, 0x46, 0x00}, false},
		{"short", []byte{0x77}, false},
		{"empty", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWOFF2(tt.data); got != tt.want {
				t.Errorf("IsWOFF2 = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildTestWOFF(t *testing.T, tables []struct {
	tag  string
	data []byte
}, flavor uint32) []byte {
	t.Helper()

	numTables := len(tables)

	type compTable struct {
		tag      uint32
		origData []byte
		compData []byte
	}

	var cts []compTable
	for _, tbl := range tables {
		tag := binary.BigEndian.Uint32([]byte(tbl.tag))
		var buf bytes.Buffer
		w, _ := flate.NewWriter(&buf, flate.DefaultCompression)
		w.Write(tbl.data)
		w.Close()

		comp := buf.Bytes()
		if len(comp) >= len(tbl.data) {
			comp = tbl.data
		}
		cts = append(cts, compTable{tag: tag, origData: tbl.data, compData: comp})
	}

	dataStart := uint32(woffHeaderSize + numTables*woffTableEntrySize)
	offset := dataStart

	var entries []woffTableEntry
	for _, ct := range cts {
		e := woffTableEntry{
			Tag:          ct.tag,
			Offset:       offset,
			CompLength:   uint32(len(ct.compData)),
			OrigLength:   uint32(len(ct.origData)),
			OrigChecksum: 0,
		}
		entries = append(entries, e)
		offset += uint32(len(ct.compData))
		if offset%4 != 0 {
			offset += 4 - (offset % 4)
		}
	}

	totalSfntSize := uint32(12 + 16*numTables)
	for _, ct := range cts {
		totalSfntSize += uint32(len(ct.origData))
		if totalSfntSize%4 != 0 {
			totalSfntSize += 4 - (totalSfntSize % 4)
		}
	}

	hdr := woffHeader{
		Signature:     0x774F4646,
		Flavor:        flavor,
		Length:        offset,
		NumTables:     uint16(numTables),
		TotalSfntSize: totalSfntSize,
	}

	var out bytes.Buffer
	binary.Write(&out, binary.BigEndian, hdr)
	for _, e := range entries {
		binary.Write(&out, binary.BigEndian, e)
	}
	for i, ct := range cts {
		out.Write(ct.compData)
		padLen := entries[i].CompLength
		if padLen%4 != 0 {
			pad := 4 - (padLen % 4)
			out.Write(make([]byte, pad))
		}
	}

	return out.Bytes()
}

func TestDecompressWOFF_Roundtrip(t *testing.T) {
	origData := []byte("Hello, this is table data that is long enough to actually compress well with flate. " +
		"Adding more text here to ensure compression actually kicks in and compLength < origLength.")

	tables := []struct {
		tag  string
		data []byte
	}{
		{"test", origData},
	}

	woffData := buildTestWOFF(t, tables, 0x00010000)

	if !IsWOFF(woffData) {
		t.Fatal("built data not detected as WOFF")
	}

	ttfData, err := DecompressWOFF(woffData)
	if err != nil {
		t.Fatalf("DecompressWOFF: %v", err)
	}

	if binary.BigEndian.Uint32(ttfData[0:4]) != 0x00010000 {
		t.Error("flavor not preserved in output")
	}

	numTables := binary.BigEndian.Uint16(ttfData[4:6])
	if numTables != 1 {
		t.Fatalf("numTables = %d, want 1", numTables)
	}

	tableOffset := binary.BigEndian.Uint32(ttfData[12+8 : 12+12])
	tableLen := binary.BigEndian.Uint32(ttfData[12+12 : 12+16])

	recovered := ttfData[tableOffset : tableOffset+tableLen]
	if !bytes.Equal(recovered, origData) {
		t.Errorf("table data mismatch:\ngot  %q\nwant %q", recovered, origData)
	}
}

func TestDecompressWOFF_UncompressedTable(t *testing.T) {
	origData := []byte("XY")

	tables := []struct {
		tag  string
		data []byte
	}{
		{"head", origData},
	}

	woffData := buildTestWOFF(t, tables, 0x00010000)
	ttfData, err := DecompressWOFF(woffData)
	if err != nil {
		t.Fatalf("DecompressWOFF: %v", err)
	}

	tableOffset := binary.BigEndian.Uint32(ttfData[12+8 : 12+12])
	tableLen := binary.BigEndian.Uint32(ttfData[12+12 : 12+16])
	recovered := ttfData[tableOffset : tableOffset+tableLen]
	if !bytes.Equal(recovered, origData) {
		t.Errorf("uncompressed table mismatch: got %q, want %q", recovered, origData)
	}
}

func TestDecompressWOFF_NotWOFF(t *testing.T) {
	_, err := DecompressWOFF([]byte{0x00, 0x01, 0x00, 0x00})
	if err == nil {
		t.Fatal("expected error for non-WOFF data")
	}
}

func TestDecompressWOFF_TooShort(t *testing.T) {
	_, err := DecompressWOFF([]byte{0x77, 0x4F, 0x46, 0x46})
	if err == nil {
		t.Fatal("expected error for truncated WOFF")
	}
}

func TestLoadFont_RejectsWOFF2(t *testing.T) {
	woff2Data := []byte{0x77, 0x4F, 0x46, 0x32, 0x00, 0x00, 0x00, 0x00}
	m := NewManager()
	err := m.LoadFont(FontSource{
		Name: "test",
		Data: woff2Data,
	})
	if err == nil {
		t.Fatal("expected error loading WOFF2")
	}
}

func TestLoadFont_AcceptsTTF(t *testing.T) {
	m := NewManager()
	err := m.LoadFont(FontSource{
		Name: "goregular",
		Data: goregular.TTF,
	})
	if err != nil {
		t.Fatalf("LoadFont TTF: %v", err)
	}
}
