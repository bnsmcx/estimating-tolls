package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var apiKey string

func main() {
  url := "https://apis.tollguru.com/v2/origin-destination-waypoints"

	payload := strings.NewReader("{\"from\":{\"address\":\"Walt Whitman Brg Philadelphia, PA 19148, USA\"},\"to\":{\"address\":\"Ocean City, NJ 08226, USA\"},\"vehicle\":{\"type\":\"5AxlesTruck\",\"height\":{\"value\":10,\"unit\":\"feet\"},\"weight\":{\"value\":45000,\"unit\":\"pounds\"},\"length\":{\"value\":30,\"unit\":\"feet\"},\"axles\":5},\"fuelOptions\":{\"fuelCost\":{\"value\":4.38,\"currency\":\"USD\",\"units\":\"$/gallon\"}},\"fuelEfficiency\":{\"city\":6.4,\"hwy\":8.2,\"units\":\"mpg\"},\"driver\":{\"wage\":30,\"rounding\":15,\"valueOfTime\":0},\"state_mileage\":true,\"hos\":{\"rule\":60,\"dutyHoursBeforeEndOfWorkDay\":11,\"dutyHoursBeforeRestBreak\":7,\"drivingHoursBeforeEndOfWorkDay\":11,\"timeRemaining\":60}}")

	req, err := http.NewRequest("post", url, payload)
  if err != nil {
    log.Fatal(err)
  }

	req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", "x-api-key: " + apiKey)

	res, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Fatal(err)
  }
	defer res.Body.Close()

  fmt.Println(res.Status)
  fmt.Printf("%+v\n", res.Request)
  fmt.Println(res.TransferEncoding)

	body, err := io.ReadAll(res.Body)
  if err != nil {
    log.Fatal(err)
  }

	fmt.Println(res)
	fmt.Println(string(body))
}

func init() {
  key, exists := os.LookupEnv("TOLL_GURU_KEY")
  if !exists {
    log.Fatal("Toll guru api key missing, set the TOLL_GURU_KEY environment variable.")
  }
  apiKey = key
}
