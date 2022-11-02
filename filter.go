package clerk

type Filter interface {
	Left() Filter
	Right() Filter
	Key() string
	Value() any
}

type And struct {
	left  Filter
	right Filter
}

func NewAnd(left Filter, right Filter) *And {
	return &And{
		left:  left,
		right: right,
	}
}

func (l *And) Left() Filter {
	return l.left
}

func (l *And) Right() Filter {
	return l.right
}

func (l *And) Key() string {
	return ""
}

func (l *And) Value() any {
	return nil
}

type Or struct {
	left  Filter
	right Filter
}

func NewOr(left Filter, right Filter) *Or {
	return &Or{
		left:  left,
		right: right,
	}
}

func (l *Or) Left() Filter {
	return l.left
}

func (l *Or) Right() Filter {
	return l.right
}

func (l *Or) Key() string {
	return ""
}

func (l *Or) Value() any {
	return nil
}

type Equals struct {
	key   string
	value any
}

func NewEquals(key string, value any) *Equals {
	return &Equals{
		key:   key,
		value: value,
	}
}

func (l *Equals) Left() Filter {
	return nil
}

func (l *Equals) Right() Filter {
	return nil
}

func (l *Equals) Key() string {
	return l.key
}

func (l *Equals) Value() any {
	return l.value
}

type GreaterThan struct {
	key   string
	value any
}

func NewGreaterThan(key string, value any) *GreaterThan {
	return &GreaterThan{
		key:   key,
		value: value,
	}
}

func (l *GreaterThan) Left() Filter {
	return nil
}

func (l *GreaterThan) Right() Filter {
	return nil
}

func (l *GreaterThan) Key() string {
	return l.key
}

func (l *GreaterThan) Value() any {
	return l.value
}

type GreaterThanOrEqual struct {
	key   string
	value any
}

func NewGreaterThanOrEqual(key string, value any) *GreaterThanOrEqual {
	return &GreaterThanOrEqual{
		key:   key,
		value: value,
	}
}

func (l *GreaterThanOrEqual) Left() Filter {
	return nil
}

func (l *GreaterThanOrEqual) Right() Filter {
	return nil
}

func (l *GreaterThanOrEqual) Key() string {
	return l.key
}

func (l *GreaterThanOrEqual) Value() any {
	return l.value
}

type LessThan struct {
	key   string
	value any
}

func NewLessThan(key string, value any) *LessThan {
	return &LessThan{
		key:   key,
		value: value,
	}
}

func (l *LessThan) Left() Filter {
	return nil
}

func (l *LessThan) Right() Filter {
	return nil
}

func (l *LessThan) Key() string {
	return l.key
}

func (l *LessThan) Value() any {
	return l.value
}

type LessThanOrEqual struct {
	key   string
	value any
}

func NewLessThanOrEqual(key string, value any) *LessThanOrEqual {
	return &LessThanOrEqual{
		key:   key,
		value: value,
	}
}

func (l *LessThanOrEqual) Left() Filter {
	return nil
}

func (l *LessThanOrEqual) Right() Filter {
	return nil
}

func (l *LessThanOrEqual) Key() string {
	return l.key
}

func (l *LessThanOrEqual) Value() any {
	return l.value
}

type Exists struct {
	key   string
	value any
}

func NewExists(value any) *Exists {
	return &Exists{
		value: value,
	}
}

func (l *Exists) Left() Filter {
	return nil
}

func (l *Exists) Right() Filter {
	return nil
}

func (l *Exists) Key() string {
	return l.key
}

func (l *Exists) Value() any {
	return l.value
}

type Regex struct {
	key   string
	value any
}

func NewRegex(key string, value any) *Regex {
	return &Regex{
		key:   key,
		value: value,
	}
}

func (l *Regex) Left() Filter {
	return nil
}

func (l *Regex) Right() Filter {
	return nil
}

func (l *Regex) Key() string {
	return l.key
}

func (l *Regex) Value() any {
	return l.value
}
