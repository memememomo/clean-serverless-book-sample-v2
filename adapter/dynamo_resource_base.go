package adapter

import "time"

type DynamoCreatedUpdated struct {
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

type DynamoResourceBase struct {
	Version int `dynamo:"Version"`
	DynamoCreatedUpdated
}
