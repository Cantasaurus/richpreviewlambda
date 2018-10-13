package main

import(
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-lambda-go/events"
  "github.com/cantasaurus/richpreview"
  "encoding/json"
  "log"
  "net/http"
)

var cors = map[string]string{"Access-Control-Allow-Origin": "*"}

type Request struct{
  Urls []string `json:"urls"`
}

type Response struct{
  Responses []*richpreview.Preview  `json:"responses"`
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
        Headers:    cors,
    }, nil
}

func PostHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
  var request Request
  response := &Response{}
  urlsReceived := make([]string, 0)
  err := json.Unmarshal([]byte(req.Body), &request)

  log.Print(req.Body)
  if err != nil{
    return clientError(http.StatusBadRequest)
  }

  for _, url := range request.Urls{
    urlsReceived = append(urlsReceived, url)
  }

  userAgent := richpreview.DefaultUserAgent()
  results := richpreview.RichPreview(urlsReceived, userAgent)

  for _, result := range results{
    response.Responses = append(response.Responses, result)
  }

  b, _ := json.Marshal(response)
  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body: string(b),
    Headers: cors,
  }, nil
}

func main(){
  lambda.Start(PostHandler)
}
