package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v4"
	sq "github.com/Masterminds/squirrel"
	"github.com/samber/lo"
)

var ErrNoSelectorInEmptySlice = errors.New("no selector in empty sclice")

func jsonKeyToSelector(column string, key string, value any) string {
	parts := strings.Split(key, ".")
	initial := func() string {
		if column == "" {
			return parts[0]
		}
		return column
	}()
	if column == "" {
		parts = parts[1:]
	}
	return lo.Reduce(parts, func(acc string, part string, i int) string {
		switch v := value.(type) {
		case string:
			if i == len(parts)-1 {
				return fmt.Sprintf("%s->>'%s'", acc, part)
			}
		case []any:
			if i == len(parts)-1 && len(v) > 0 {
				if _, ok := v[0].(string); ok {
					return fmt.Sprintf("%s->>'%s'", acc, part)
				}
			}
		}
		return fmt.Sprintf("%s->'%s'", acc, part)
	}, initial)
}

func typeCastSelector(selector string, value any) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("(%s)::int", selector)
	case []int:
		return fmt.Sprintf("(%s)::int", selector)
	case int32:
		return fmt.Sprintf("(%s)::int", selector)
	case []int32:
		return fmt.Sprintf("(%s)::int", selector)
	case int64:
		return fmt.Sprintf("(%s)::bigint", selector)
	case []int64:
		return fmt.Sprintf("(%s)::bigint", selector)
	case float32:
		return fmt.Sprintf("(%s)::float", selector)
	case []float32:
		return fmt.Sprintf("(%s)::float", selector)
	case float64:
		return fmt.Sprintf("(%s)::float", selector)
	case []float64:
		return fmt.Sprintf("(%s)::float", selector)
	case bool:
		return fmt.Sprintf("(%s)::bool", selector)
	case []bool:
		return fmt.Sprintf("(%s)::bool", selector)
	case []any:
		if len(v) == 0 {
			return selector
		}
		return typeCastSelector(selector, v[0])
	default:
		return selector
	}
}

func setOfVariables(len int) string {
	variables := "("
	for i := 0; i < len; i++ {
		if i == 0 {
			variables += "?"
		} else {
			variables += ", ?"
		}
	}
	variables += ")"
	return variables
}

func filtersToCondition(column string, filters ...clerk.Filter) (sq.Sqlizer, error) {
	var condition sq.Sqlizer
	for _, filter := range filters {
		resolvedFilter, err := filterToCondition(column, filter)
		if err != nil {
			return nil, err
		}
		if condition == nil {
			condition = resolvedFilter
		} else {
			condition = sq.And{condition, resolvedFilter}
		}
	}
	return condition, nil
}

func filterToCondition(column string, filter clerk.Filter) (sq.Sqlizer, error) {
	switch filter := filter.(type) {
	case *clerk.And:
		left, err := filterToCondition(column, filter.Left())
		if err != nil {
			return nil, err
		}
		right, err := filterToCondition(column, filter.Right())
		if err != nil {
			return nil, err
		}
		return sq.And{left, right}, nil
	case *clerk.Or:
		left, err := filterToCondition(column, filter.Left())
		if err != nil {
			return nil, err
		}
		right, err := filterToCondition(column, filter.Right())
		if err != nil {
			return nil, err
		}
		return sq.Or{left, right}, nil
	case *clerk.Equals:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s = ?", selector),
			filter.Value(),
		), nil
	case *clerk.NotEquals:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s != ?", selector),
			filter.Value(),
		), nil
	case *clerk.GreaterThan:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s > ?", selector),
			filter.Value(),
		), nil
	case *clerk.GreaterThanOrEquals:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s >= ?", selector),
			filter.Value(),
		), nil
	case *clerk.LessThan:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s < ?", selector),
			filter.Value(),
		), nil
	case *clerk.LessThanOrEquals:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s <= ?", selector),
			filter.Value(),
		), nil
	case *clerk.Exists:
		exists, ok := filter.Value().(bool)
		if !ok || !exists {
			return sq.Expr(
				fmt.Sprintf("NOT(%s ?? ?)", column),
				filter.Key(),
			), nil
		}
		return sq.Expr(
			fmt.Sprintf("%s ?? ?", column),
			filter.Key(),
		), nil
	case *clerk.Regex:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Value()),
			filter.Value(),
		)
		return sq.Expr(
			fmt.Sprintf("%s SIMILAR TO ?", selector),
			filter.Value(),
		), nil
	case *clerk.In:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Values()),
			filter.Values(),
		)
		variables := setOfVariables(len(filter.Values()))
		return sq.Expr(
			fmt.Sprintf("%s IN %s", selector, variables),
			filter.Values()...,
		), nil
	case *clerk.NotIn:
		selector := typeCastSelector(
			jsonKeyToSelector(column, filter.Key(), filter.Values()),
			filter.Values(),
		)
		variables := setOfVariables(len(filter.Values()))
		return sq.Expr(
			fmt.Sprintf("%s NOT IN %s", selector, variables),
			filter.Values()...,
		), nil
	case *clerk.InArray:
		selector := jsonKeyToSelector(column, filter.Key(), filter.Values())
		variables := setOfVariables(len(filter.Values()))
		typeCastedValues := typeCastSelector("values", filter.Values())
		return sq.Expr(
			fmt.Sprintf("EXISTS( SELECT TRUE FROM jsonb_array_elements(%s) values WHERE %s IN %s )", selector, typeCastedValues, variables),
			filter.Values()...,
		), nil
	case *clerk.NotInArray:
		selector := jsonKeyToSelector(column, filter.Key(), filter.Values())
		variables := setOfVariables(len(filter.Values()))
		typeCastedValues := typeCastSelector("values", filter.Values())
		return sq.Expr(
			fmt.Sprintf("NOT EXISTS( SELECT TRUE FROM jsonb_array_elements(%s) values WHERE %s IN %s )", selector, typeCastedValues, variables),
			filter.Values()...,
		), nil
	default:
		return nil, fmt.Errorf("%w: %T", clerk.ErrorInvalidFilter, filter)
	}
}
