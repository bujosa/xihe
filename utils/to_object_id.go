package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectId(id interface{}) (primitive.ObjectID, error) {
	switch id := id.(type) {
	case primitive.ObjectID:
		return id, nil
	case string:
		if objectId, err := primitive.ObjectIDFromHex(id); err != nil {
			return primitive.NilObjectID, err
		} else {
			return objectId, nil
		}
	default:
		return primitive.NilObjectID, fmt.Errorf("invalid type for id: %T", id)
	}
}
