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

func EditBirthday(coll *dynamodb.DynamoDB, table string, name string, mobile int64, bday birthday.Birthday) events.LambdaFunctionURLResponse {

	if len(name) == 0 || mobile == 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "input the query strings of name and mobile."}
	}
	jsonData, err := json.Marshal(bday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: "unable to unmarshall body."}
	}
	err = database.Edit(coll, table, name, mobile, bday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("Some issue while editing %s.",string(jsonData))}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusPartialContent, Body: fmt.Sprintf("added changes to %s.",string(jsonData))}
}
