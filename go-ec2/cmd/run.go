package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run a new EC2 instance",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(Region())},
		)

		// Create EC2 service client
		svc := ec2.New(sess)

		// Specify the details of the instance that you want to create.
		runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
			// An Amazon Linux AMI ID for t2.micro instances in the eu-central-1 region
			ImageId:      aws.String("ami-0233214e13e500f77"),
			InstanceType: aws.String("t2.micro"),
			MinCount:     aws.Int64(1),
			MaxCount:     aws.Int64(1),
		})

		if err != nil {
			fmt.Println("Could not create EC2 instance.", err)
			return
		}

		fmt.Println("Created EC2 instance:", *runResult.Instances[0].InstanceId)

		// Add tags to the created instance
		_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{runResult.Instances[0].InstanceId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(args[0]),
				},
				{
					Key:   aws.String("ChaosEnabled"),
					Value: aws.String("true"),
				},
			},
		})
		if errtag != nil {
			log.Println("Could not tag EC2 instance.", runResult.Instances[0].InstanceId, errtag)
			return
		}

		fmt.Println("Successfully tagged EC2 instance.")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
