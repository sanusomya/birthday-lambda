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

func DeleteBirthday(coll *dynamodb.DynamoDB, table string, birthday birthday.Birthday) events.LambdaFunctionURLResponse {

	err := database.Delete(coll, table, birthday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: fmt.Sprintf("unable to delete %s", birthday)}
	}
	jsonData, err := json.Marshal(birthday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: ""}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Body: fmt.Sprintf("sucessfully deleted %s",string(jsonData))}
}
