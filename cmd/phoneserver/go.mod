module github.com/mudphilo/phonenumber/cmd/phoneserver

go 1.24.5

replace github.com/mudphilo/phonenumber => ../..

require (
	github.com/aws/aws-lambda-go v1.13.1
	github.com/mudphilo/phonenumber v0.0.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
