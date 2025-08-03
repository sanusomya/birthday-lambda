package routes

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/sanusomya/birthday-lambda/database"
	utils "github.com/sanusomya/birthday-lambda/utils"
)

func GetAllBirthdaysToday(coll *dynamodb.DynamoDB ,month string, date int8, table string) events.LambdaFunctionURLResponse {
	
	if len(month) == 0 || date == 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "input the query strings of month and date"}
	}

	month = utils.CorrectMonth(month)
	
	errMsg := utils.CustomError{}
	
	if !utils.CheckDates(date,month){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "date"
		errMsg.Message = "Date and Month combination not right. Date should be with allowed calander values for the month."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}  


	birthdays, err := database.FindForToday(coll, table, month, date)
	jsonData, err := json.Marshal(birthdays)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: ""}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Body: string(jsonData)}
}
