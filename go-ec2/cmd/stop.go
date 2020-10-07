package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an EC2 instance",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(Region())},
		)

		// Create EC2 service client
		svc := ec2.New(sess)

		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(args[0]),
			},
			DryRun: aws.Bool(false),
		}

		result, err := svc.StopInstances(input)

		if err != nil {
			fmt.Println("Error stopping EC2 instance.", err)
		} else {
			fmt.Println("Stopped EC2 instance:", result.StoppingInstances)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
