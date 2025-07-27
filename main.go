package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sanusomya/birthday-lambda/database"
	birthday "github.com/sanusomya/birthday-lambda/models"
	"github.com/sanusomya/birthday-lambda/routes"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	path := request.RequestContext.HTTP.Path

	var response events.LambdaFunctionURLResponse

	// create dtabase collection
	coll := database.ConnectDB()

	// table name
	table := os.Getenv("table")

	switch path {
	// section for get all
	case fmt.Sprintf("/%s/api", os.Getenv("stage")):
		method := request.RequestContext.HTTP.Method
		if method == "GET" {
			response = routes.GetAllBirthdays(coll, table)
		} else if method == "POST"{
			bday := birthday.Birthday{}
			err := json.Unmarshal([]byte(request.Body), &bday)
			if err != nil {
				response = events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "cannot decode body"}
			}
			response = routes.AddBirthday(coll, table, bday)
		}else{
			bday := birthday.Birthday{}
			err := json.Unmarshal([]byte(request.Body), &bday)
			query := request.QueryStringParameters
			mobile, _ := strconv.Atoi(query["mobile"])
			name := query["name"]
			if err != nil {
				response = events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest, Body: "cannot decode body"}
			}
			response = routes.EditBirthday(coll, table, name, int64(mobile), bday)
		}

	// section for get all today
	case fmt.Sprintf("/%s/api/today", os.Getenv("stage")):
		query := request.QueryStringParameters
		date, _ := strconv.Atoi(query["date"])
		month := query["month"]
		response = routes.GetAllBirthdaysToday(coll, month, int8(date), table)

	// section for get all this month
	case fmt.Sprintf("/%s/api/month", os.Getenv("stage")):
		query := request.QueryStringParameters
		month := query["month"]
		response = routes.GetAllBirthdaysMonth(coll, month, table)

	// section for editing name
	case fmt.Sprintf("/%s/api/name", os.Getenv("stage")):
		query := request.QueryStringParameters
		name := query["name"]
		mobile, _ := strconv.Atoi(query["mobile"])
		response = routes.EditBirthdayName(coll, table, name, int64(mobile), request.Body)

	// section for editing number
	case fmt.Sprintf("/%s/api/number", os.Getenv("stage")):
		query := request.QueryStringParameters
		name := query["name"]
		mobile, _ := strconv.Atoi(query["mobile"])
		body, _ := strconv.Atoi(request.Body)
		response = routes.EditBirthdayNumber(coll, table, name, int64(mobile), int64(body))

	default:
		response = events.LambdaFunctionURLResponse{StatusCode: 404, Body: "Not Found"}
	}

	return response, nil

}
