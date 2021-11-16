package storage

import "github.com/bububa/visiondb/pb"

type Storage interface {
	Flush() error
	Reload() error
	Truncate() error
	Shape() (*pb.Shape, error)
	Labels() ([]string, error)
	AddLabel(label string) error
	DeleteLabel(labelID int) error
	GetLabelByID(id int) (string, error)
	GetLabelItems(labelID int) ([]*pb.Item, error)
	AddLabelItems(labelID int, items ...*pb.Item) error
	DeleteLabelItems(labelID int, itemIndexes ...int) error
	ResetLabelItems(labelID int) error
	ChangeLabelName(labelID int, name string) error
}
