package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// terminateCmd represents the terminate command
var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "Terminate an EC2 instance",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(Region())},
		)

		// Create EC2 service client
		svc := ec2.New(sess)

		input := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String(args[0]),
			},
			DryRun: aws.Bool(false),
		}

		result, err := svc.TerminateInstances(input)

		if err != nil {
			fmt.Println("Error terminating EC2 instance.", err)
		} else {
			fmt.Println("Terminated EC2 instance:", result.TerminatingInstances)
		}
	},
}

func init() {
	rootCmd.AddCommand(terminateCmd)
}
