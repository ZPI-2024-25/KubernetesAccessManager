package models

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Read   OperationType = "read"
	Delete OperationType = "delete"
	List   OperationType = "list"
	All    OperationType = "all"
	all    string 	 	 = "all"
)

type Operation struct {
	Resource  string `json:"resource,omitempty"`
	Type      OperationType `json:"type,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (o *Operation) IsSuper(operation *Operation) bool {
	if (o.Type == All || operation.Type == o.Type) &&  
	(o.Resource == all || operation.Resource == o.Resource) &&
	(o.Namespace == all || operation.Namespace == o.Namespace) {
		return true
	}

	return false
}

