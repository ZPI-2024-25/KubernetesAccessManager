package models

type OperationType string

const (
	Create OperationType = "create"
	Update OperationType = "update"
	Read   OperationType = "read"
	Delete OperationType = "delete"
	List   OperationType = "list"
	All    OperationType = "*"
	all    string        = "*"
)

type Operation struct {
	Resource  string        `json:"resource,omitempty"`
	Type      OperationType `json:"operation,omitempty"`
	Namespace string        `json:"namespace,omitempty"`
}

func GetAllOperationTypes() []OperationType {
    return []OperationType{
        Create,
        Read,
        Update,
        Delete,
        List,
    }
}

func (permission OperationType)ShortString() string {
	switch permission {
	case Create:
		return "c"
	case Read:
		return "r"
	case Update:
		return "u"
	case Delete:
		return "d"
	case List:
		return "l"
	default:
		return "x"
	}
}

func (o *Operation) IsSuper(operation *Operation) bool {
	return (o.Type == All || operation.Type == o.Type) &&
		(o.Resource == all || operation.Resource == o.Resource) &&
		(o.Namespace == all || operation.Namespace == o.Namespace)
}
