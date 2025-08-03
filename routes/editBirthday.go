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

func EditBirthday(coll *dynamodb.DynamoDB, table string, name string, mobile int64, bday birthday.Birthday) events.LambdaFunctionURLResponse {

	if len(name) == 0 || mobile == 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "input the query strings of name and mobile."}
	}

	bname := bday.Person
	bmobile := bday.Cell
	date := bday.Birthdate
	month := bday.Birthmonth  
	errMsg := utils.CustomError{}

	month = utils.CorrectMonth(month)
	bday.Birthmonth = month
	
	if !utils.CheckDates(date,month){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "date"
		errMsg.Message = "Date and Month combination not right. Date should be with allowed calander values for the month."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}  

	if !utils.ValidMobile(mobile) || !utils.ValidMobile(bmobile){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "mobile"
		errMsg.Message = "Not a valide mobile number. Mobile number should be 10 digits, not starting with a 0."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	if !utils.ValidName(name) || !utils.ValidName(bname) {
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "name"
		errMsg.Message = "Not a valide Name. Length of name shouldnt exceed 9 and should only contain Alphabets."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
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
