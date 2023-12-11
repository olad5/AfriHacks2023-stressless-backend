package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jaswdr/faker"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	seedUser()
	seedMetrics()
}

var layout = "2006-01-02T15:04:05.999Z"

func seedUser() {
	f, e := os.Create("./seed-users.csv")
	if e != nil {
		fmt.Println(e)
	}

	writer := csv.NewWriter(f)
	heading := []string{"_id", "email", "password", "first_name", "last_name", "is_onboarding_complete", "last_metric_log", "updated_at", "created_at"}
	err := writer.Write(heading)
	if err != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	users := getUsers()
	var row []string
	for _, user := range users {
		row = []string{user.ID.Hex(), user.Email, user.Password, user.FirstName, user.LastName, strconv.FormatBool(user.IsOnBoardingComplete), user.LastMetricLog.Format(layout), user.UpdatedAt.Format(layout), user.CreatedAt.Format(layout)}
		e := writer.Write(row)
		if e != nil {
			fmt.Println(e)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
	}
}

func getUsers() []domain.User {
	hashedPassword, err := hashAndSalt([]byte("some-password"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return []domain.User{
		{
			ID:                   MustObjectIDFromHex("64a72bf98fc6411d15248485"),
			Email:                "jasontodd@gmail.com",
			FirstName:            "jason",
			LastName:             "todd",
			Password:             hashedPassword,
			IsOnBoardingComplete: true,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   MustObjectIDFromHex("649dccc44c16a88a075cb19d"),
			Email:                "joecole@yahoo.co.uk",
			FirstName:            "joe",
			LastName:             "cole",
			Password:             hashedPassword,
			IsOnBoardingComplete: true,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   MustObjectIDFromHex("649dccc44c16a88a075cb19e"),
			Email:                "mileyjohnson@gmail.com",
			FirstName:            "miley",
			LastName:             "johnson",
			Password:             hashedPassword,
			IsOnBoardingComplete: true,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
	}
}

func seedMetrics() {
	f, e := os.Create("./seed-metrics.csv")
	if e != nil {
		fmt.Println(e)
	}

	writer := csv.NewWriter(f)

	heading := []string{"_id", "owner_id", "mood", "sleep_quality", "stress_less_score", "stress_level", "feeling", "created_at", "updated_at"}
	err := writer.Write(heading)
	if err != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	fakerInstance := faker.New()
	users := getUsers()
	var row []string
	for _, user := range users {
		for i := 0; i < 7; i++ {
			metric := domain.Metric{
				ID:              primitive.NewObjectID(),
				OwnerId:         user.ID,
				StressLevel:     randomIntWithMaxValueInclusive(1, 5),
				StressLessScore: randomIntWithMaxValueInclusive(1, 100),
				SleepQuality:    randomSleepQuality(),
				Mood:            randomMood(),
				Feeling:         fakerInstance.App().Faker.Lorem().Paragraph(2),
				CreatedAt:       createPreviousIsoDateFromString(i),
				UpdatedAt:       createPreviousIsoDateFromString(i),
			}
			row = []string{metric.ID.Hex(), metric.OwnerId.Hex(), string(metric.Mood), string(metric.SleepQuality), fmt.Sprint(metric.StressLessScore), fmt.Sprint(metric.StressLevel), metric.Feeling, metric.CreatedAt.Format(layout), metric.UpdatedAt.Format(layout)}
			e := writer.Write(row)
			if e != nil {
				fmt.Println(e)
			}
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
	}
}

func MustObjectIDFromHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return id
}

func hashAndSalt(plainPassword []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(plainPassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	return string(hash), nil
}

func createPreviousIsoDateFromString(numOfDaysFromCurrentTime int) time.Time {
	return time.Now().AddDate(0, 0, -numOfDaysFromCurrentTime)
}

func randomIntWithMaxValueInclusive(min, max int) int {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time
	return rand.Intn(max-min+1) + min
}

func randomSleepQuality() domain.SleepQuality {
	qualities := []string{"excellent", "good", "fair", "poor", "worst"}
	randomIndex := randomIntWithMaxValueInclusive(0, len(qualities)-1)
	return domain.SleepQuality(qualities[randomIndex])
}

func randomMood() domain.Mood {
	moods := []string{"overjoyed", "happy", "neutral", "sad", "depressed"}
	randomIndex := randomIntWithMaxValueInclusive(0, len(moods)-1)
	return domain.Mood(moods[randomIndex])
}
