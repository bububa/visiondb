package pb

import "github.com/bububa/visiondb/utils"

// Hash generate hash for item
func (i *Item) GenHash() string {
	var bs []byte
	raw := i.GetRaw()
	if raw != nil {
		bs = raw
	} else {
		bs = utils.Float32SliceToBytes(i.GetEmbedding())
	}
	i.Hash = utils.Sha1(bs)
	return i.Hash
}
