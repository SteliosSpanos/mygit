package objects


type Blob struct {
	Data []byte
}


func NewBlob(data []byte) *Blob {
	return &Blob{Data: data}
}


func (b *Blob) Type() ObjectType {
	return BlobObject
}


func (b *Blob) Serialize() ([]byte, error) {
	return b.Data, nil
}


func (b *Blob) Deserialize(data []byte) error {
	b.Data = data
	return nil
}
