package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SteliosSpanos/mygit/pkg/objects"
)

func WriteObject(gitDir string, obj objects.Object) (string, error) {
	hash, err := objects.Hash(obj)
	if err != nil {
		return "", fmt.Errorf("failed to hash object: %w", err)
	}

	data, err := obj.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}

	header := fmt.Sprintf("%s %d\x00", obj.Type(), len(data))
	fullData := append([]byte(header), data...)

	compressed, err := objects.Compress(fullData)
	if err != nil {
		return "", fmt.Errorf("failed to compress object: %w", err)
	}

	objectDir := filepath.Join(gitDir, "objects", hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create object directory: %w", err)
	}

	if err := os.WriteFile(objectPath, compressed, 0644); err != nil {
		return "", fmt.Errorf("failed to write object file: %w", err)
	}

	return hash, nil
}

func ReadObject(gitDir, hash string) (objects.ObjectType, []byte, error) {
	objectPath := filepath.Join(gitDir, "objects", hash[:2], hash[2:])

	compressed, err := os.ReadFile(objectPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read object file: %w", err)
	}

	decompressed, err := objects.Decompress(compressed)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decompress object: %w", err)
	}

	nullIdx := bytes.IndexByte(decompressed, 0)
	if nullIdx == -1 {
		return "", nil, fmt.Errorf("invalid object format: no null byte")
	}

	header := string(decompressed[:nullIdx])
	content := decompressed[nullIdx+1:]

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid object header: %s", header)
	}

	objectType := objects.ObjectType(parts[0])

	return objectType, content, nil
}

func LoadObject(gitDir, hash string) (objects.Object, error) {
	objType, content, err := ReadObject(gitDir, hash)
	if err != nil {
		return nil, err
	}

	var obj objects.Object

	switch objType {
	case objects.BlobObject:
		obj = objects.NewBlob(nil)
	case objects.TreeObject:
		obj = objects.NewTree()
	case objects.CommitObject:
		obj = objects.NewCommit("", "", "")
	default:
		return nil, fmt.Errorf("unknown object type: %s", objType)
	}

	if err := obj.Deserialize(content); err != nil {
		return nil, fmt.Errorf("failed to deserialize object: %w", err)
	}

	return obj, nil
}
