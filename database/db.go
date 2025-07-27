package database

import (
	"strconv"

	birthday "github.com/sanusomya/birthday-lambda/models"
	// "github.com/sanusomya/birthday-backend/utils"
	// "context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func ConnectDB() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return dynamodb.New(sess)
}

func GetAll(svc *dynamodb.DynamoDB, tableName string) ([]birthday.Birthday, error) {
	var birth []birthday.Birthday
	ans, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return []birthday.Birthday{}, err
	}
	if *ans.Count == 0 {
		return []birthday.Birthday{}, nil
	}
	for _, item := range ans.Items {
		var temp birthday.Birthday
		dynamodbattribute.UnmarshalMap(item, &temp)
		birth = append(birth, temp)
	}
	return birth, nil
}

func Add(svc *dynamodb.DynamoDB, tableName string, b birthday.Birthday) error {

	av, err := dynamodbattribute.MarshalMap(b)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	return err
}

func Delete(svc *dynamodb.DynamoDB, tableName string, b birthday.Birthday) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Person": {
				S: aws.String(b.Person),
			},
			"Cell": {
				N: aws.String(strconv.Itoa(int(b.Cell))),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	return err
}

func Edit(svc *dynamodb.DynamoDB, tableName string, name string, mobile int64, b birthday.Birthday) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":m": {
				S: aws.String(b.Birthmonth),
			},
			":d": {
				N: aws.String(strconv.Itoa(int(b.Birthdate))),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Person": {
				S: aws.String(name),
			},
			"Cell": {
				N: aws.String(strconv.Itoa(int(mobile))),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Birthmonth = :m, Birthdate = :d"),
	}

	_, err := svc.UpdateItem(input)
	return err
}


func FindForThisMonth(svc *dynamodb.DynamoDB, tableName string, mon string) ([]birthday.Birthday, error) {
	var bdays []birthday.Birthday
	filt := expression.Name("Birthmonth").Equal(expression.Value(mon))
	proj := expression.NamesList(expression.Name("Person"), expression.Name("Birthdate"), expression.Name("Cell"), expression.Name("Birthmonth"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return bdays, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		return bdays, err
	}
	for _, i := range result.Items {
		bday := birthday.Birthday{}

		err = dynamodbattribute.UnmarshalMap(i, &bday)
		if err != nil {
			return bdays, err
		}
		bdays = append(bdays, bday)
	}
	return bdays, err
}

func FindForToday(svc *dynamodb.DynamoDB, tableName string, mon string, date int8) ([]birthday.Birthday, error) {
	var bdays []birthday.Birthday
	filt := expression.And(expression.Name("Birthmonth").Equal(expression.Value(mon)), expression.Name("Birthdate").Equal(expression.Value(date)))
	proj := expression.NamesList(expression.Name("Person"), expression.Name("Birthdate"), expression.Name("Cell"), expression.Name("Birthmonth"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return bdays, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		return bdays, err
	}
	for _, i := range result.Items {
		bday := birthday.Birthday{}

		err = dynamodbattribute.UnmarshalMap(i, &bday)
		if err != nil {
			return bdays, err
		}
		bdays = append(bdays, bday)
	}
	return bdays, err
}

func Get(svc *dynamodb.DynamoDB, tableName string, name string, mobile int64) (birthday.Birthday, error) {
	bday := birthday.Birthday{} 
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Person": {
				S: aws.String(name),
			},
			"Cell": {
				N: aws.String(strconv.Itoa(int(mobile))),
			},
		},
	})
	if err != nil {
		return birthday.Birthday{},err
	}
	
	err = dynamodbattribute.UnmarshalMap(result.Item, &bday)
	if err != nil {
		return birthday.Birthday{},err
	}
	return bday,err
}
