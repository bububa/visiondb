package storage

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/bububa/visiondb/logger"
	"github.com/bububa/visiondb/pb"
	"google.golang.org/protobuf/proto"
)

type labelHashMap map[string]struct{}

func (h labelHashMap) Add(k string) {
	h[k] = struct{}{}
}

func (h labelHashMap) Has(k string) bool {
	_, found := h[k]
	return found
}

func (h labelHashMap) Del(k string) {
	delete(h, k)
}

// ProtoBufStorage represents protobufer storage
type ProtoBufStorage struct {
	dbPath string
	db     *pb.Database
	hashes []labelHashMap
	mutex  *sync.RWMutex
}

// NewProtoBufStorage returns a new ProtoBufStorage
func NewProtoBufStorage(dbPath string) *ProtoBufStorage {
	return &ProtoBufStorage{
		dbPath: dbPath,
		db:     new(pb.Database),
		mutex:  new(sync.RWMutex),
	}
}

// Lock lock db
func (p *ProtoBufStorage) Lock() {
	p.mutex.Lock()
}

// Unlock unlock db
func (p *ProtoBufStorage) Unlock() {
	p.mutex.Unlock()
}

// RLock read lock
func (p *ProtoBufStorage) RLock() {
	p.mutex.RLock()
}

// RUnlock read unlock
func (p *ProtoBufStorage) RUnlock() {
	p.mutex.RUnlock()
}

// Reload read db data from protobufer file
func (p *ProtoBufStorage) Reload() error {
	data, err := os.ReadFile(p.dbPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		logger.Error().Err(err).Interface("err", err).Send()
		return err
	}
	p.Lock()
	defer p.Unlock()
	if err = proto.Unmarshal(data, p.db); err != nil {
		return err
	}
	labels := p.db.GetLabels()
	p.hashes = make([]labelHashMap, 0, len(labels))
	for _, label := range labels {
		items := label.GetItems()
		mp := make(labelHashMap, len(items))
		for _, itm := range items {
			hash := itm.GetHash()
			if hash == "" {
				hash = itm.GenHash()
			}
			mp.Add(hash)
		}
		p.hashes = append(p.hashes, mp)
	}
	return nil
}

// Flush flush db data to protobufer file
func (p *ProtoBufStorage) Flush() error {
	fn, err := os.Create(p.dbPath)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer fn.Close()
	p.RLock()
	defer p.RUnlock()
	data, err := proto.Marshal(p.db)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if _, err := fn.Write(data); err != nil {
		log.Fatalln(err)
		return err
	}
	return fn.Sync()
}

// Truncate reset db
func (p *ProtoBufStorage) Truncate() error {
	p.Lock()
	defer p.Unlock()
	p.db.Reset()
	p.hashes = p.hashes[:0]
	return nil
}

// Shape returns the data shape
func (p *ProtoBufStorage) Shape() (*pb.Shape, error) {
	p.RLock()
	defer p.RUnlock()
	return p.db.GetShape(), nil
}

// SetShape update data shape
func (p *ProtoBufStorage) SetShape(shape *pb.Shape) error {
	p.RLock()
	defer p.RUnlock()
	p.db.Shape = shape
	return nil
}

// Labels returns labels' names
func (p *ProtoBufStorage) Labels() ([]string, []int, error) {
	p.RLock()
	defer p.RUnlock()
	labels := p.db.GetLabels()
	names := make([]string, 0, len(labels))
	counts := make([]int, 0, len(labels))
	for _, l := range labels {
		names = append(names, l.GetName())
		counts = append(counts, len(l.GetItems()))
	}
	return names, counts, nil
}

// GetLabelByID returns label name by labelID
func (p *ProtoBufStorage) GetLabelByID(id int) (string, int, error) {
	p.RLock()
	defer p.RUnlock()
	labels := p.db.GetLabels()
	if id >= len(labels) {
		return "", 0, errors.New("invalid labelID")
	}
	label := labels[id]
	return label.GetName(), len(label.GetItems()), nil
}

// GetLabelItem returns Item by labelID and ItemID
func (p *ProtoBufStorage) GetLabelItem(labelID int, itemID int) (*pb.Item, error) {
	items, err := p.GetLabelItems(labelID)
	if err != nil {
		return nil, err
	}
	if itemID >= len(items) {
		return nil, errors.New("invalid itemID")
	}
	return items[itemID], nil
}

// GetLabelItems returns []Item by labelID
func (p *ProtoBufStorage) GetLabelItems(labelID int) ([]*pb.Item, error) {
	p.RLock()
	defer p.RUnlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return nil, errors.New("invalid labelID")
	}
	return labels[labelID].GetItems(), nil
}

// AddLabelItems add label items
func (p *ProtoBufStorage) AddLabelItems(labelID int, items ...*pb.Item) (int, int, error) {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return 0, 0, errors.New("invalid labelID")
	}
	label := labels[labelID]
	retID := len(label.Items)
	var count int
	for _, itm := range items {
		hash := itm.GetHash()
		if p.hashes[labelID].Has(hash) {
			continue
		}
		label.Items = append(label.Items, itm)
		p.hashes[labelID].Add(hash)
		count++
	}
	return retID, count, nil
}

// DeleteLabelItems delete label items by item index
func (p *ProtoBufStorage) DeleteLabelItems(labelID int, itemIndexes ...int) (int, error) {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return 0, errors.New("invalid labelID")
	}
	label := labels[labelID]
	sort.Ints(itemIndexes)
	var deletes int
	for _, i := range itemIndexes {
		l := len(label.GetItems())
		idx := i - deletes
		if idx >= l {
			return deletes, errors.New("invalid index")
		}
		p.hashes[labelID].Del(label.Items[idx].GetHash())
		label.Items[idx] = label.Items[l-1]
		label.Items[l-1] = nil
		label.Items = label.Items[:l-1]
		deletes++
	}
	return deletes, nil
}

// AddLabel add a new Label
func (p *ProtoBufStorage) AddLabel(name string) (int, error) {
	p.Lock()
	defer p.Unlock()
	label := new(pb.Label)
	label.Name = name
	p.db.Labels = append(p.db.Labels, label)
	p.hashes = append(p.hashes, make(labelHashMap))
	return len(p.db.Labels) - 1, nil
}

// DeleteLabel remove a Label by index
func (p *ProtoBufStorage) DeleteLabel(labelID int) error {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	l := len(labels)
	if labelID >= l {
		return errors.New("invalid labelID")
	}
	p.db.Labels[labelID] = p.db.Labels[l-1]
	p.db.Labels[l-1] = nil
	p.db.Labels = p.db.Labels[:l-1]
	p.hashes[labelID] = p.hashes[l-1]
	p.hashes[l-1] = nil
	p.hashes = p.hashes[:l-1]
	return nil
}

// ResetLabelItems clear items in label
func (p *ProtoBufStorage) ResetLabelItems(labelID int) error {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return errors.New("invalid labelID")
	}
	label := labels[labelID]
	labelName := label.GetName()
	label.Reset()
	label.Name = labelName
	p.hashes[labelID] = make(labelHashMap)
	return nil
}

// ChangeLabelName update label's name
func (p *ProtoBufStorage) ChangeLabelName(labelID int, name string) error {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return errors.New("invalid labelID")
	}
	labels[labelID].Name = name
	return nil
}
