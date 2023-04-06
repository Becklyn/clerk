package postgres

import (
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/stretchr/testify/assert"
)

func Test_JsonKeyToSelector(t *testing.T) {
	selector := jsonKeyToSelector("data", "a.b.c", nil, true)
	assert.Equal(t, "data->'a'->'b'->'c'", selector)
}

func Test_JsonKeyToSelector_WithEmptyColumn(t *testing.T) {
	selector := jsonKeyToSelector("", "a.b.c", nil, true)
	assert.Equal(t, "a->'b'->'c'", selector)
}

func Test_JsonKeyToSelector_WithStringValue(t *testing.T) {
	selector := jsonKeyToSelector("data", "a.b.c", "foo", true)
	assert.Equal(t, "data->'a'->'b'->>'c'", selector)
}

func Test_TypeCastSelector_Int(t *testing.T) {
	selector := typeCastSelector("selector", 1)
	assert.Equal(t, "(selector)::int", selector)
}

func Test_TypeCastSelector_SliceOfInt(t *testing.T) {
	selector := typeCastSelector("selector", []int{1, 2, 3})
	assert.Equal(t, "(selector)::int", selector)
}

func Test_TypeCastSelector_SliceOfInt_Empty(t *testing.T) {
	selector := typeCastSelector("selector", []int{})
	assert.Equal(t, "(selector)::int", selector)
}

func Test_TypeCastSelector_Float(t *testing.T) {
	selector := typeCastSelector("selector", 1.0)
	assert.Equal(t, "(selector)::float", selector)
}

func Test_TypeCastSelector_SliceOfFloat(t *testing.T) {
	selector := typeCastSelector("selector", []float64{1.0, 2.0, 3.0})
	assert.Equal(t, "(selector)::float", selector)
}

func Test_TypeCastSelector_Bool(t *testing.T) {
	selector := typeCastSelector("selector", true)
	assert.Equal(t, "(selector)::bool", selector)
}

func Test_TypeCastSelector_SliceOfBool(t *testing.T) {
	selector := typeCastSelector("selector", []bool{true, false, true})
	assert.Equal(t, "(selector)::bool", selector)
}

func Test_TypeCastSelector_String(t *testing.T) {
	selector := typeCastSelector("selector", "foo")
	assert.Equal(t, "selector", selector)
}

func Test_TypeCastSelector_SliceOfString(t *testing.T) {
	selector := typeCastSelector("selector", []string{"foo", "bar", "baz"})
	assert.Equal(t, "selector", selector)
}

func Test_SetOfVariables(t *testing.T) {
	vars := setOfVariables(3)
	assert.Equal(t, "(?, ?, ?)", vars)
}

func Test_FiltersToCondition_WithAndFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewAnd(
			clerk.NewEquals("a", 1),
			clerk.NewEquals("b.c", 2),
		),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "((data->'a')::int = ? AND (data->'b'->'c')::int = ?)", stat)
	assert.Equal(t, []any{1, 2}, vals)
}

func Test_FiltersToCondition_WithOrFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewOr(
			clerk.NewEquals("a", 1),
			clerk.NewEquals("b.c", 2),
		),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "((data->'a')::int = ? OR (data->'b'->'c')::int = ?)", stat)
	assert.Equal(t, []any{1, 2}, vals)
}

func Test_FiltersToCondition_WithEqualsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewEquals("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int = ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithNotEqualsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewNotEquals("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int != ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithGreaterThanFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewGreaterThan("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int > ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithGreaterThanOrEqualsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewGreaterThanOrEquals("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int >= ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithLessThanFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewLessThan("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int < ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithLessThanOrEqualsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewLessThanOrEquals("a", 1),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int <= ?", stat)
	assert.Equal(t, []any{1}, vals)
}

func Test_FiltersToCondition_WithExistsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewExists("a", true),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "data ?? ?", stat)
	assert.Equal(t, []any{"a"}, vals)
}

func Test_FiltersToCondition_WithNotExistsFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewExists("a", false),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "NOT(data ?? ?)", stat)
	assert.Equal(t, []any{"a"}, vals)
}

func Test_FiltersToCondition_WithRegexFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewRegex("a", "b"),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "data->>'a' SIMILAR TO ?", stat)
	assert.Equal(t, []any{"b"}, vals)
}

func Test_FiltersToCondition_WithInFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewIn("a", 1, 2),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int IN (?, ?)", stat)
	assert.Equal(t, []any{1, 2}, vals)
}

func Test_FiltersToCondition_WithNotInFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewNotIn("a", 1, 2),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "(data->'a')::int NOT IN (?, ?)", stat)
	assert.Equal(t, []any{1, 2}, vals)
}

func Test_FiltersToCondition_WithInArrayFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewInArray("a", 1, 2),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "EXISTS( SELECT TRUE FROM jsonb_array_elements(data->'a') AS x(o) WHERE x.o IN (?, ?) )", stat)
	assert.Equal(t, []any{1, 2}, vals)
}

func Test_FiltersToCondition_WithNotInArrayFilter(t *testing.T) {
	condition, err := filtersToCondition(
		"data",
		clerk.NewNotInArray("a", 1, 2),
	)
	assert.NoError(t, err)
	stat, vals, err := condition.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "NOT EXISTS( SELECT TRUE FROM jsonb_array_elements(data->'a') AS x(o) WHERE x.o IN (?, ?) )", stat)
	assert.Equal(t, []any{1, 2}, vals)
}
