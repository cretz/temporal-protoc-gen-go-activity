package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DataDog/temporalite/temporaltest"
	"github.com/cretz/temporal-protoc-gen-go-activity/example/greetingspb"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	// Start test server
	log.Print("Starting server")
	ts := temporaltest.NewServer()
	defer ts.Stop()

	// Start worker on a task queue that registers our workflow and activities
	const taskQueue = "my-task-queue"
	log.Print("Starting worker")
	ts.Worker(taskQueue, func(registry worker.Registry) {
		registry.RegisterWorkflow(GreetingSample)
		registry.RegisterActivity(&activities{Name: "Temporal", Greeting: "Hello"})
	})

	// Execute the workflow and get result
	startOpts := client.StartWorkflowOptions{TaskQueue: taskQueue}
	run, err := ts.Client().ExecuteWorkflow(context.TODO(), startOpts, GreetingSample)
	if err != nil {
		log.Fatalf("Failed starting workflow: %v", err)
	}
	var sayResult string
	if err := run.Get(context.TODO(), &sayResult); err != nil {
		log.Fatalf("Failed running workflow: %v", err)
	}

	log.Printf("Workflow result: %v", sayResult)
}

// GreetingSample workflow definition.
// This greetings sample workflow executes 3 activities in sequential.
// It gets greeting and name from 2 different activities,
// and then pass greeting and name as input to a 3rd activity to generate final greetings.
func GreetingSample(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Get greeting
	greetingResult, err := greetingspb.Greetings.GetGreeting(ctx, nil)
	if err != nil {
		logger.Error("Get greeting failed.", "Error", err)
		return "", err
	}

	// Get name
	nameResult, err := greetingspb.Greetings.GetName(ctx, nil)
	if err != nil {
		logger.Error("Get name failed.", "Error", err)
		return "", err
	}

	// Say Greeting
	sayResult, err := greetingspb.Greetings.SayGreeting(ctx, &greetingspb.SayGreetingRequest{
		Greeting: greetingResult.Greeting,
		Name:     nameResult.Name,
	})
	if err != nil {
		logger.Error("Say greeting failed.", "Error", err)
		return "", err
	}

	logger.Info("GreetingSample completed.", "Result", sayResult.Greeting)
	return sayResult.Greeting, nil
}

type activities struct {
	Name     string
	Greeting string
}

var _ greetingspb.GreetingsImpl = &activities{}

func (a *activities) GetGreeting(
	context.Context,
	*greetingspb.GetGreetingRequest,
) (*greetingspb.GetGreetingResponse, error) {
	return &greetingspb.GetGreetingResponse{Greeting: a.Greeting}, nil
}

func (a *activities) GetName(context.Context, *greetingspb.GetNameRequest) (*greetingspb.GetNameResponse, error) {
	return &greetingspb.GetNameResponse{Name: a.Name}, nil
}

func (a *activities) SayGreeting(
	ctx context.Context,
	in *greetingspb.SayGreetingRequest,
) (*greetingspb.SayGreetingResponse, error) {
	return &greetingspb.SayGreetingResponse{
		Greeting: fmt.Sprintf("Greeting: %s %s!\n", in.Greeting, in.Name),
	}, nil
}
