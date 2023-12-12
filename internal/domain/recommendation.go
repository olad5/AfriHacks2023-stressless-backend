package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecommendationItem struct {
	Index    int
	Heading  string
	Text     string
	ImageUrl string
}
type Recommendation struct {
	ID         primitive.ObjectID
	MetricId   primitive.ObjectID
	MetricType string
	Items      []RecommendationItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
