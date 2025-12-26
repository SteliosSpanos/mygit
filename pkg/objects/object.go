package objects
import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
)


type ObjectType string

const (
	BlobObject ObjectType = "blob"
	TreeObject ObjectType = "tree"
	CommitObject ObjectType = "commit"
)


type Object interface {
	Type() ObjectType
	Serialize() ([]byte, error) //Convert object to bytes
	Deserialize(data []byte) error //Parse bytes into object
}


func Hash(obj Object) (string, error) {
	data, err := obj.Serialize()
	if err != nil {
		return "", err
	}

	//Create the git object format: <type> <size>\0<content>
	header := fmt.Sprintf("%s %d\x00", obj.Type(), len(data))
	store := append([]byte(header), data...)

	hash := sha1.Sum(store)
	return fmt.Sprintf("%x", hash), nil
}



func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}



func Decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	return io.ReadAll(reader)
}
