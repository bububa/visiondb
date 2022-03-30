package storage

import (
	"github.com/bububa/visiondb/pb"
)

type Storage interface {
	Flush() error
	Reload() error
	Truncate() error
	Shape() (*pb.Shape, error)
	SetShape(*pb.Shape) error
	Labels() ([]string, []int, error)
	AddLabel(label string) (int, error)
	DeleteLabel(labelID int) error
	GetLabelByID(id int) (string, int, error)
	GetLabelItems(labelID int) ([]*pb.Item, error)
	GetLabelItem(labelID int, itemID int) (*pb.Item, error)
	AddLabelItems(labelID int, items ...*pb.Item) (int, int, error)
	DeleteLabelItems(labelID int, itemIndexes ...int) (int, error)
	ResetLabelItems(labelID int) error
	ChangeLabelName(labelID int, name string) error
}
