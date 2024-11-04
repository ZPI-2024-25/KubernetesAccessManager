package models

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Read   OperationType = "read"
	Delete OperationType = "delete"
	List   OperationType = "list"
	All    OperationType = "*"
	all    string 	 	 = "*"
)

type Operation struct {
	Resource  string `json:"resource,omitempty"`
	Type      OperationType `json:"type,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (o *Operation) IsSuper(operation *Operation) bool {
    return (o.Type == All || operation.Type == o.Type) &&
           (o.Resource == All || operation.Resource == o.Resource) &&
           (o.Namespace == All || operation.Namespace == o.Namespace)
}

