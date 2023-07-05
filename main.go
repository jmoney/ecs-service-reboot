package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// returns a boolean if string is nil or empty
func isEmpty(s *string) bool {
	return s == nil || len(*s) == 0
}

func main() {
	cluster := flag.String("cluster", os.Getenv("CLUSTER"), "ECS cluster name")
	serviceName := flag.String("service", os.Getenv("SERVICE_NAME"), "ECS service name")
	awsRegion := flag.String("region", os.Getenv("AWS_REGION"), "AWS region")
	flag.Parse()

	if isEmpty(cluster) || isEmpty(serviceName) || isEmpty(awsRegion) {
		panic("Missing required parameters")
	}
	var err error
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: awsRegion,
		},
	))

	svc := ecs.New(sess)

	logger.Printf("Rebooting ecs service \"%s\" in cluster \"%s\"", *serviceName, *cluster)
	_, err = svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:            cluster,
		Service:            serviceName,
		ForceNewDeployment: aws.Bool(true),
	})

	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(5 * time.Second)

		service, err := svc.DescribeServices(&ecs.DescribeServicesInput{
			Cluster:  cluster,
			Services: []*string{serviceName},
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
