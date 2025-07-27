package routes

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/sanusomya/birthday-lambda/database"
)

func GetAllBirthdaysMonth(coll *dynamodb.DynamoDB ,month string, table string) events.LambdaFunctionURLResponse {
	
	if len(month) == 0{
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "input the query strings of month"}
	}
	birthdays, err := database.FindForThisMonth(coll, table, month)
	jsonData, err := json.Marshal(birthdays)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: ""}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Body: string(jsonData)}
}
