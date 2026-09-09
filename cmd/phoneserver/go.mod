module github.com/mudphilo/phonenumber/cmd/phoneserver

go 1.26.0

replace github.com/mudphilo/phonenumber => ../..

require (
	github.com/aws/aws-lambda-go v1.13.1
	github.com/mudphilo/phonenumber v0.0.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/nyaruka/phonenumbers/v2 v2.0.11 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
