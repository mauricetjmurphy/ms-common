package eventbridge

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/NBCUniversal/gvs-ms-common/clients/aws/lambda/utils"
	"github.com/NBCUniversal/gvs-ms-common/logx"
)

type Handler func(ctx context.Context, event events.CloudWatchEvent) error

// StartFunc invokes Lambda when run inside a Lambda function triggered from event bridge event.
func StartFunc(f Handler) {
	if utils.IsLambdaEnv() {
		lambda.Start(f)
	} else {
		runLocal(f)
	}
}

func runLocal(f Handler) {
	params, err := utils.ParseArgs()
	if err != nil {
		logx.Panicf("lambda :  parse arguments on err %v", err)
	}
	logx.Infof("lambda : locally run the lambda with params %v", params)
	body, err := json.Marshal(params)
	if err != nil {
		logx.Panicf("lambda : unable to parse json body %v", err)
	}
	req := events.CloudWatchEvent{
		Detail: body,
	}
	err = f(context.Background(), req)
	if err != nil {
		logx.Panicf("lambda : execute lambda local on err %v", err)
	}
	logx.Infof("lambda : completed run the lambda on local")
}
