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

func EditBirthdayNumber(coll *dynamodb.DynamoDB, table string, name string, mobile int64, body int64) events.LambdaFunctionURLResponse {

	if len(name) == 0 || mobile == 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "input the query strings of name and mobile."}
	}

	errMsg := utils.CustomError{}
	
	if !utils.ValidMobile(mobile) || !utils.ValidMobile(body){
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "mobile"
		errMsg.Message = "Not a valide mobile number. Mobile number should be 10 digits, not starting with a 0."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	if !utils.ValidName(name) {
		errMsg.StatusCode = http.StatusBadRequest
		errMsg.Attribute = "name"
		errMsg.Message = "Not a valide Name. Length of name shouldnt exceed 9 and should only contain Alphabets."
		body,_ := json.Marshal(errMsg)
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: string(body)}
	}

	bday, err := database.Get(coll, table, name, int64(mobile))
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: fmt.Sprintf("cannot find birthday %s",body)} 
	}
	jsonData, err := json.Marshal(bday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: ""}
	}
	err = database.Delete(coll, table, bday)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: fmt.Sprintf("Some error occured while deleting %s",string(jsonData))} 
	}
	temp := birthday.Birthday{
		Person:     bday.Person,
		Cell:       body,
		Birthdate:  bday.Birthdate,
		Birthmonth: bday.Birthmonth,
	}
	err = database.Add(coll, table, temp)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError, Body: fmt.Sprintf("Some error occured while adding %s",string(jsonData))} 
	}

	jsonData, _ = json.Marshal(temp)
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusCreated, Body: fmt.Sprintf("sucessfully edited %s",string(jsonData))}
}
