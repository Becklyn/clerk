package mongodb

import (
	"fmt"

	"github.com/Becklyn/clerk/v4"

	"go.mongodb.org/mongo-driver/bson"
)

func resolveFilters(filters ...clerk.Filter) (bson.D, error) {
	bsonFilters := bson.D{}
	for _, filter := range filters {
		parsedFilter, err := resolveFilter(filter)
		if err != nil {
			return nil, err
		}
		bsonFilters = append(bsonFilters, parsedFilter)
	}
	return bsonFilters, nil
}

func resolveFilter(filter clerk.Filter) (bson.E, error) {
	switch filter := filter.(type) {
	case *clerk.And:
		left, err := resolveFilter(filter.Left())
		if err != nil {
			return bson.E{}, err
		}
		right, err := resolveFilter(filter.Right())
		if err != nil {
			return bson.E{}, err
		}
		return bson.E{
			Key: "$and",
			Value: bson.A{
				bson.D{left},
				bson.D{right},
			},
		}, nil
	case *clerk.Or:
		left, err := resolveFilter(filter.Left())
		if err != nil {
			return bson.E{}, err
		}
		right, err := resolveFilter(filter.Right())
		if err != nil {
			return bson.E{}, err
		}
		return bson.E{
			Key: "$or",
			Value: bson.A{
				bson.D{left},
				bson.D{right},
			},
		}, nil
	case *clerk.Equals:
		return bson.E{
			Key:   filter.Key(),
			Value: filter.Value(),
		}, nil
	case *clerk.NotEquals:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$ne",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.GreaterThan:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$gt",
					Value: filter.Value(),
				},
			},
		}, nil

	case *clerk.GreaterThanOrEquals:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$gte",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.LessThan:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$lt",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.LessThanOrEquals:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$lte",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.Exists:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$exists",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.Regex:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$regex",
					Value: filter.Value(),
				},
			},
		}, nil
	case *clerk.In:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$in",
					Value: filter.Values(),
				},
			},
		}, nil
	case *clerk.InArray:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$in",
					Value: filter.Values(),
				},
			},
		}, nil
	case *clerk.NotIn:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$nin",
					Value: filter.Values(),
				},
			},
		}, nil
	case *clerk.NotInArray:
		return bson.E{
			Key: filter.Key(),
			Value: bson.D{
				{
					Key:   "$nin",
					Value: filter.Values(),
				},
			},
		}, nil
	default:
		return bson.E{}, fmt.Errorf("%w: %T", clerk.ErrorInvalidFilter, filter)
	}
}
