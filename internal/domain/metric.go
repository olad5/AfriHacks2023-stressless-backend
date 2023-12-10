package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mood string

const (
	OVERJOYED Mood = "overjoyed"
	HAPPY     Mood = "happy"
	NEUTRAL   Mood = "neutral"
	SAD       Mood = "sad"
	DEPRESSED Mood = "depressed"
)

type SleepQuality string

const (
	EXCELLENT SleepQuality = "excellent"
	GOOD      SleepQuality = "good"
	FAIR      SleepQuality = "fair"
	POOR      SleepQuality = "poor"
	WORST     SleepQuality = "worst"
)

type Metric struct {
	ID              primitive.ObjectID
	OwnerId         primitive.ObjectID
	StressLevel     int
	Mood            Mood
	SleepQuality    SleepQuality
	Feeling         string
	StressLessScore int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
