package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	logger      *log.Logger
	cluster     string
	serviceName string
)

func init() {
	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	cluster = os.Getenv("CLUSTER")
	serviceName = os.Getenv("SERVICE_NAME")
}

func main() {
	var err error
	sess := session.Must(session.NewSession())

	svc := ecs.New(sess)

	logger.Printf("Rebooting ecs service")
	_, err = svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:            aws.String(cluster),
		Service:            aws.String(serviceName),
		ForceNewDeployment: aws.Bool(true),
	})

	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(5 * time.Second)

		service, err := svc.DescribeServices(&ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: []*string{aws.String(serviceName)},
		})

		if err == nil {
			for _, deployment := range service.Services[0].Deployments {
				if *deployment.Status == "PRIMARY" {
					logger.Printf("%s|%d tasks running|%d tasks pending|%d tasks desired|%d tasks failed", *deployment.RolloutState, *deployment.RunningCount, *deployment.PendingCount, *deployment.DesiredCount, *deployment.FailedTasks)
					if deployment.RolloutState != nil && (*deployment.RolloutState == "COMPLETED" || *deployment.RolloutState == "FAILED") {
						logger.Printf("Rollout complete: %s", *deployment.RolloutStateReason)
						return
					}
				}
			}
		} else {
			logger.Printf("Error: %s", err.Error())
		}
	}
}
