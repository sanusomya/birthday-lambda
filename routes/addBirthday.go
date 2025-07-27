package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/sanusomya/birthday-lambda/database"
	birthday "github.com/sanusomya/birthday-lambda/models"
)

func AddBirthday(coll *dynamodb.DynamoDB, table string, birthday birthday.Birthday) events.LambdaFunctionURLResponse {

	err := database.Add(coll, table, birthday)
	jsonData, err := json.Marshal(birthday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "error occured while unmarshaling"}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusCreated, Body: fmt.Sprintf("sucessfully added %s",string(jsonData))}
}
