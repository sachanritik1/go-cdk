package database

import (
	"context"
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	USERS_TABLE = "users"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}

type DynamoDBClient struct {
	databaseStore *dynamodb.Client
}

func NewDynamoDBClient() (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-south-1"),
	)
	if err != nil {
		return nil, err
	}

	db := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{
		databaseStore: db,
	}, nil
}

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.databaseStore.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE),
		Key: map[string]dynamodbtypes.AttributeValue{
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: username,
			},
		},
		ConsistentRead: aws.Bool(true), // optional but recommended for existence checks
	})

	// Handle potential errors from the GetItem operation
	if err != nil {
		return true, err
	}

	// If Item is nil or empty, user does not exist
	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	_, err := u.databaseStore.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(USERS_TABLE),
		Item: map[string]dynamodbtypes.AttributeValue{
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: user.Username,
			},
			"password": &dynamodbtypes.AttributeValueMemberS{
				Value: user.PasswordHash,
			},
		},
	})

	return err
}

func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	result, err := u.databaseStore.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE),
		Key: map[string]dynamodbtypes.AttributeValue{
			"username": &dynamodbtypes.AttributeValueMemberS{
				Value: username,
			},
		},
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return types.User{}, err
	}

	if result.Item == nil {
		return types.User{}, fmt.Errorf("User not found")
	}

	user := types.User{
		Username:     username,
		PasswordHash: result.Item["password"].(*dynamodbtypes.AttributeValueMemberS).Value,
	}

	return user, nil
}
