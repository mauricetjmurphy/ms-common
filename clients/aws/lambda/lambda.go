package lambda

import (
	"context"

	"github.com/NBCUniversal/gvs-ms-common/clients/aws/lambda/utils"
	"github.com/NBCUniversal/gvs-ms-common/logx"
	"github.com/aws/aws-lambda-go/lambda"
)

type Func func(ctx context.Context, payload interface{}) error

// StartFunc invokes Lambda when run inside a Lambda function with generic events, otherwise just execute it locally.
func StartFunc(f Func, params ...interface{}) {
	if utils.IsLambdaEnv() {
		lambda.Start(f)
	} else {
		runLocal(f, params)
	}
}

func runLocal(f Func, params ...interface{}) {
	logx.Infof("lambda : execute the lambda on local with params %v", params)
	err := f(context.Background(), params)
	if err != nil {
		logx.Panicf("lambda : failed to execute lambda local %v on params %v", err, params)
	}
	logx.Infof("lambda : completed run the lambda on local")
}
