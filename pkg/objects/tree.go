package objects

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string
}

type Tree struct {
	Entries []TreeEntry
}

func NewTree() *Tree {
	return &Tree{Entries: make([]TreeEntry, 0)}
}

func (t *Tree) AddEntry(mode, name, hash string) {
	t.Entries = append(t.Entries, TreeEntry{
		Mode: mode,
		Name: name,
		Hash: hash,
	})
}

func (t *Tree) Type() ObjectType {
	return TreeObject
}

func (t *Tree) Serialize() ([]byte, error) {
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})

	var buf bytes.Buffer

	for _, entry := range t.Entries {
		buf.WriteString(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name))
		hashBytes, err := HexToBytes(entry.Hash)
		if err != nil {
			return nil, fmt.Errorf("invalid hash for %s: %w", entry.Name, err)
		}
		buf.Write(hashBytes)
	}

	return buf.Bytes(), nil
}

func (t *Tree) Deserialize(data []byte) error {
	t.Entries = make([]TreeEntry, 0)

	for len(data) > 0 {
		nullIdx := bytes.IndexByte(data, 0) //First occurance of 0
		if nullIdx == -1 {
			return fmt.Errorf("invalid tree format: no null byte")
		}

		header := string(data[:nullIdx])
		parts := strings.SplitN(header, " ", 2)
		if len(parts) > 2 {
			return fmt.Errorf("invalid tree format: %s", header)
		}

		mode := parts[0]
		name := parts[1]

		hashStart := nullIdx + 1
		if len(data) < hashStart+20 {
			return fmt.Errorf("invalid tree format: not enough data for hash")
		}

		hashBytes := data[hashStart : hashStart+20]
		hash := fmt.Sprintf("%x", hashBytes)

		t.AddEntry(mode, name, hash)

		data = data[hashStart+20:] //Move to next extry
	}
	return nil
}

func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr) != 40 {
		return nil, fmt.Errorf("hash must be 40 characters, got %d", len(hexStr))
	}

	bytes := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b, err := strconv.ParseUint(hexStr[i*2:i*2+2], 16, 8) //Each slice is one byte (2 hex)
		if err != nil {
			return nil, err
		}

		bytes[i] = byte(b)
	}

	return bytes, nil
}
