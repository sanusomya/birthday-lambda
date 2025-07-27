package routes

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/sanusomya/birthday-lambda/database"
)

func GetAllBirthdays(coll *dynamodb.DynamoDB, table string) events.LambdaFunctionURLResponse {

	birthdays, err := database.GetAll(coll, table)
	jsonData, err := json.Marshal(birthdays)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: ""}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Body: string(jsonData)}
}
