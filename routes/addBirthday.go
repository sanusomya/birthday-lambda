package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	database "github.com/sanusomya/birthday-lambda/database"
	birthday "github.com/sanusomya/birthday-lambda/models"
	utils "github.com/sanusomya/birthday-lambda/utils"
)

func AddBirthday(coll *dynamodb.DynamoDB, table string, birthday birthday.Birthday) events.LambdaFunctionURLResponse {

	name := birthday.Person
	mobile := birthday.Cell
	date := birthday.Birthdate
	month := birthday.Birthmonth
	errMsg := utils.CustomError{}

	month = utils.CorrectMonth(month)
	birthday.Birthmonth = month

	if !utils.CheckDates(date,month){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "date"
		errMsg.Message = "Date and Month combination not right. Date should be with allowed calander values for the month."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	if !utils.ValidMobile(mobile){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "mobile"
		errMsg.Message = "Not a valide mobile number. Mobile number should be 10 digits, not starting with a 0."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	if !utils.ValidName(name){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "name"
		errMsg.Message = "Not a valide Name. Length of name shouldnt exceed 9 and should only contain Alphabets."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	err := database.Add(coll, table, birthday)
	jsonData, err := json.Marshal(birthday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "error occured while unmarshaling"}
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusCreated, Body: fmt.Sprintf("sucessfully added %s", string(jsonData))}
}
