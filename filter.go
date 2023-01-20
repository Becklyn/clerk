package clerk

import "errors"

var (
	ErrorInvalidFilter = errors.New("invalid filter")
)

type Filter interface {
	Left() Filter
	Right() Filter
	Key() string
	Value() any
	Values() []any
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

func (l *And) Values() []any {
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

func (l *Or) Values() []any {
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

func (l *Equals) Values() []any {
	return nil
}

type In struct {
	key    string
	values []any
}

func NewIn(key string, values ...any) *In {
	return &In{
		key:    key,
		values: values,
	}
}

func (i *In) Left() Filter {
	return nil
}

func (i *In) Right() Filter {
	return nil
}

func (i *In) Key() string {
	return i.key
}

func (i *In) Value() any {
	return nil
}

func (i *In) Values() []any {
	return i.values
}

type NotIn struct {
	key    string
	values []any
}

func NewNotIn(key string, values ...any) *NotIn {
	return &NotIn{
		key:    key,
		values: values,
	}
}

func (i *NotIn) Left() Filter {
	return nil
}

func (i *NotIn) Right() Filter {
	return nil
}

func (i *NotIn) Key() string {
	return i.key
}

func (i *NotIn) Value() any {
	return nil
}

func (i *NotIn) Values() []any {
	return i.values
}

type InArray struct {
	key    string
	values []any
}

func NewInArray(key string, values ...any) *InArray {
	return &InArray{
		key:    key,
		values: values,
	}
}

func (l *InArray) Left() Filter {
	return nil
}

func (l *InArray) Right() Filter {
	return nil
}

func (l *InArray) Key() string {
	return l.key
}

func (l *InArray) Value() any {
	return nil
}

func (l *InArray) Values() []any {
	return l.values
}

type NotInArray struct {
	key    string
	values []any
}

func NewNotInArray(key string, values ...any) *NotInArray {
	return &NotInArray{
		key:    key,
		values: values,
	}
}

func (l *NotInArray) Left() Filter {
	return nil
}

func (l *NotInArray) Right() Filter {
	return nil
}

func (l *NotInArray) Key() string {
	return l.key
}

func (l *NotInArray) Value() any {
	return nil
}

func (l *NotInArray) Values() []any {
	return l.values
}

type NotEquals struct {
	key   string
	value any
}

func NewNotEquals(key string, value any) *NotEquals {
	return &NotEquals{
		key:   key,
		value: value,
	}
}

func (l *NotEquals) Left() Filter {
	return nil
}

func (l *NotEquals) Right() Filter {
	return nil
}

func (l *NotEquals) Key() string {
	return l.key
}

func (l *NotEquals) Value() any {
	return l.value
}

func (l *NotEquals) Values() []any {
	return nil
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

func (l *GreaterThan) Values() []any {
	return nil
}

type GreaterThanOrEquals struct {
	key   string
	value any
}

func NewGreaterThanOrEquals(key string, value any) *GreaterThanOrEquals {
	return &GreaterThanOrEquals{
		key:   key,
		value: value,
	}
}

func (l *GreaterThanOrEquals) Left() Filter {
	return nil
}

func (l *GreaterThanOrEquals) Right() Filter {
	return nil
}

func (l *GreaterThanOrEquals) Key() string {
	return l.key
}

func (l *GreaterThanOrEquals) Value() any {
	return l.value
}

func (l *GreaterThanOrEquals) Values() []any {
	return nil
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

func (l *LessThan) Values() []any {
	return nil
}

type LessThanOrEquals struct {
	key   string
	value any
}

func NewLessThanOrEquals(key string, value any) *LessThanOrEquals {
	return &LessThanOrEquals{
		key:   key,
		value: value,
	}
}

func (l *LessThanOrEquals) Left() Filter {
	return nil
}

func (l *LessThanOrEquals) Right() Filter {
	return nil
}

func (l *LessThanOrEquals) Key() string {
	return l.key
}

func (l *LessThanOrEquals) Value() any {
	return l.value
}

func (l *LessThanOrEquals) Values() []any {
	return nil
}

type Exists struct {
	key   string
	value any
}

func NewExists(key string, value any) *Exists {
	return &Exists{
		key:   key,
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

func (l *Exists) Values() []any {
	return nil
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

func (l *Regex) Values() []any {
	return nil
}
