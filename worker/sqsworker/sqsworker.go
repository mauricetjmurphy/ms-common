package sqsworker

import (
	"context"
	"sync"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/mauricetjmurphy/ms-common/clients/aws/sqs"
	"github.com/pkg/errors"
)

// TaskFn defines the function to be performed on the task data.
type TaskFn func(msg types.Message) error

// WorkersParams presents the workers properties to be created.
type WorkersParams struct {
	// Presents the SQS configuration properties.
	*sqs.ClientParams

	// NumWorkers is the number of concurrent workers to spawn.
	NumWorkers int

	// Tasks define the functions that are run to handle an SQS message.
	Tasks TaskFn

	// Logger define the logging functions.
	Logger Logger
}

var recoverFunc func(error) error
var recoverFuncOnce sync.Once

// SetRecoverFunc sets the function to be called if a panic is recovered when calling a task function.
func SetRecoverFunc(f func(err error) error) {
	recoverFuncOnce.Do(func() {
		recoverFunc = f
	})
}

// New creates the worker instance to read messages from an SQS queue and handle them via "task" functions.
// A "done" channel allows the caller to gracefully kill workers.
func New(ctx context.Context, params WorkersParams, stop <-chan struct{}) (*sync.WaitGroup, <-chan error, error) {
	if err := validateParams(params); err != nil {
		return nil, nil, err
	}

	if params.Logger == nil {
		params.Logger = &StandardLogger{}
	}

	client, err := sqs.New(ctx, params.ClientParams)
	if err != nil {
		return nil, nil, errors.Wrap(err, "sqsworker : failed to load sqs client")
	}

	decoratedTasks := decorateTaskFn(params.Tasks)

	errCh := make(chan error)
	var wg sync.WaitGroup

	for i := 0; i < params.NumWorkers; i++ {
		wg.Add(1)
		go func(_ int) {
			defer wg.Done()
			w := &workerImpl{
				client: client,
				tasks:  decoratedTasks,
				errCh:  errCh,
				stop:   stop,
				logger: params.Logger,
			}
			w.loop()
		}(i)
	}

	return &wg, errCh, nil
}

type workerImpl struct {
	client sqs.Client
	tasks  TaskFn
	errCh  chan<- error
	stop   <-chan struct{}
	logger Logger
}

// loop repeatedly fetches and handles jobs from SQS
func (w *workerImpl) loop() {
	for {
		select {
		case <-w.stop:
			return
		default:
			if err := w.execute(); err != nil && err != ErrNoMessageInQueue {
				w.handleError(err)
			}
		}
	}
}

func (w *workerImpl) handleError(err error) {
	select {
	case w.errCh <- err:
	default:
	}
}

func (w *workerImpl) execute() error {
	resp, err := w.callReceiveMessages(1)
	if err != nil {
		return errors.Wrap(err, "sqsworker : failed to call callReceiveMessages")
	}
	if len(resp.Messages) == 0 {
		return ErrNoMessageInQueue
	}
	message := resp.Messages[0]

	w.logger.Debugf("sqsworker : in process on mgs %v", *message.MessageId)

	err = w.tasks(message)
	if err != nil {
		return errors.Wrap(err, "sqsworker : failed to execute task on mgs "+*message.MessageId)
	}

	if err := w.deleteTask(message); err != nil {
		w.logger.Errorf("sqsworker : failed to delete on mgs %v", *message.MessageId)
	}

	w.logger.Debugf("sqsworker : completed process on mgs %v", *message.MessageId)

	return nil
}

// deleteTask removes a task from the queue, typically after successfully processing it.
func (w *workerImpl) deleteTask(message types.Message) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return w.client.DeleteMessage(ctx, message)
}

type receiveMessageResult struct {
	err  error
	resp *awssqs.ReceiveMessageOutput
}

func (w *workerImpl) callReceiveMessages(maxReceivedMessages int) (*awssqs.ReceiveMessageOutput, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultCh := make(chan receiveMessageResult, 1)
	go func() {
		defer close(resultCh)
		resp, err := w.client.ReceiveMessages(ctx, maxReceivedMessages)
		resultCh <- receiveMessageResult{
			err:  err,
			resp: resp,
		}
	}()

	select {
	case <-w.stop:
		return &awssqs.ReceiveMessageOutput{}, nil
	case result := <-resultCh:
		return result.resp, result.err
	}
}

func validateParams(params WorkersParams) error {
	switch {
	case params.NumWorkers < 1:
		return errors.New("invalid NumWorkers")
	case params.QueueURL == "":
		return errors.New("invalid QueueUrl")
	case params.Region == "":
		return errors.New("invalid Region - required because no Client is set")
	}
	return nil
}

// decorateTaskFn adds panic protection to the TaskFn
func decorateTaskFn(fn TaskFn) TaskFn {
	return func(msg types.Message) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
				if recoverFunc != nil {
					err = recoverFunc(err)
				}
			}
		}()
		return fn(msg)
	}
}
