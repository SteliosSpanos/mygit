package index

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Mode string
	Hash string
	Path string
}

type Index struct {
	Entries []Entry
}

func NewIndex() *Index {
	return &Index{
		Entries: make([]Entry, 0),
	}
}

func (idx *Index) Add(mode, hash, path string) {
	for i, entry := range idx.Entries {
		if entry.Path == path {
			idx.Entries[i] = Entry{
				Mode: mode,
				Hash: hash,
				Path: path,
			}

			return
		}
	}

	idx.Entries = append(idx.Entries, Entry{
		Mode: mode,
		Hash: hash,
		Path: path,
	})
}

func (idx *Index) Remove(path string) bool {
	for i, entry := range idx.Entries {
		if entry.Path == path {
			idx.Entries = append(idx.Entries[:i], idx.Entries[i+1:]...)
			return true
		}
	}

	return false
}

func (idx *Index) Get(path string) (*Entry, bool) {
	for _, entry := range idx.Entries {
		if entry.Path == path {
			return &entry, true
		}
	}

	return nil, false
}

func ReadIndex(gitDir string) (*Index, error) {
	indexPath := filepath.Join(gitDir, "index")

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return NewIndex(), nil
	}

	file, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}
	defer file.Close()

	idx := NewIndex()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid index line: %s", line)
		}

		idx.Add(parts[0], parts[1], parts[2])
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading index: %w", err)
	}

	return idx, nil
}

func WriteIndex(gitDir string, idx *Index) error {
	indexPath := filepath.Join(gitDir, "index")

	file, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, entry := range idx.Entries {
		line := fmt.Sprintf("%s %s %s\n", entry.Mode, entry.Hash, entry.Path)
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write index entry: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush index: %w", err)
	}

	return nil
}
