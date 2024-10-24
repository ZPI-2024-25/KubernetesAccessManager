package models

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Read   OperationType = "read"
	Delete OperationType = "delete"
	List   OperationType = "list"
)

type Operation struct {
	Resource  string
	Type      OperationType
	Namespace string
}
