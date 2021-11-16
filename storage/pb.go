package storage

import (
	"errors"
	"os"
	"sort"
	"sync"

	"github.com/bububa/visiondb/pb"
	"google.golang.org/protobuf/proto"
)

// ProtoBufStorage represents protobufer storage
type ProtoBufStorage struct {
	dbPath string
	db     *pb.Database
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
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}
	p.Lock()
	defer p.Unlock()
	if err = proto.Unmarshal(data, p.db); err != nil {
		return err
	}
	return nil
}

// Flush flush db data to protobufer file
func (p *ProtoBufStorage) Flush() error {
	fn, err := os.Create(p.dbPath)
	if err != nil {
		return err
	}
	defer fn.Close()
	p.RLock()
	defer p.RUnlock()
	data, err := proto.Marshal(p.db)
	if err != nil {
		return err
	}
	if _, err := fn.Write(data); err != nil {
		return err
	}
	return fn.Sync()
}

// Truncate reset db
func (p *ProtoBufStorage) Truncate() error {
	p.Lock()
	defer p.Unlock()
	p.db.Reset()
	return nil
}

// Shape returns the data shape
func (p *ProtoBufStorage) Shape() (*pb.Shape, error) {
	p.RLock()
	defer p.RUnlock()
	return p.db.GetShape(), nil
}

// Labels returns labels' names
func (p *ProtoBufStorage) Labels() ([]string, error) {
	p.RLock()
	defer p.RUnlock()
	labels := p.db.GetLabels()
	names := make([]string, 0, len(labels))
	for _, l := range labels {
		names = append(names, l.GetName())
	}
	return names, nil
}

// GetLabelByID returns label name by labelID
func (p *ProtoBufStorage) GetLabelByID(id int) (string, error) {
	p.RLock()
	defer p.RUnlock()
	labels := p.db.GetLabels()
	if id >= len(labels) {
		return "", errors.New("invalid labelID")
	}
	return labels[id].GetName(), nil
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
func (p *ProtoBufStorage) AddLabelItems(labelID int, items ...*pb.Item) error {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return errors.New("invalid labelID")
	}
	label := labels[labelID]
	label.Items = append(label.Items, items...)
	return nil
}

// DeleteLabelItems delete label items by item index
func (p *ProtoBufStorage) DeleteLabelItems(labelID int, itemIndexes ...int) error {
	p.Lock()
	defer p.Unlock()
	labels := p.db.GetLabels()
	if labelID >= len(labels) {
		return errors.New("invalid labelID")
	}
	label := labels[labelID]
	sort.Ints(itemIndexes)
	for _, i := range itemIndexes {
		l := len(label.GetItems())
		if i >= l {
			return errors.New("invalid index")
		}
		label.Items[i] = label.Items[l-1]
		label.Items[l-1] = nil
		label.Items = label.Items[:l-1]
	}
	return nil
}

// AddLabel add a new Label
func (p *ProtoBufStorage) AddLabel(label string) error {
	p.Lock()
	defer p.Unlock()
	p.db.Labels = append(p.db.Labels, new(pb.Label))
	return nil
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
