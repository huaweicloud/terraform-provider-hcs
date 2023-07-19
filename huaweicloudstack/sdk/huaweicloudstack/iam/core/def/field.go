package def

type LocationType int32

const (
	Header LocationType = 1 << iota
	Path
	Query
	Body
	Form
	Cname
)

type FieldDef struct {
	LocationType LocationType
	Name         string
	JsonTag      string
	KindName     string
}

func NewFieldDef() *FieldDef {
	return &FieldDef{}
}

func (field *FieldDef) WithLocationType(locationType LocationType) *FieldDef {
	field.LocationType = locationType
	return field
}

func (field *FieldDef) WithName(name string) *FieldDef {
	field.Name = name
	return field
}

func (field *FieldDef) WithJsonTag(tag string) *FieldDef {
	field.JsonTag = tag
	return field
}

func (field *FieldDef) WithKindName(kindName string) *FieldDef {
	field.KindName = kindName
	return field
}
