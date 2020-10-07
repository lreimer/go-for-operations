package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an EC2 instance",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(Region())},
		)

		// Create EC2 service client
		svc := ec2.New(sess)

		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(args[0]),
			},
			DryRun: aws.Bool(false),
		}

		result, err := svc.StartInstances(input)

		if err != nil {
			fmt.Println("Error starting EC2 instance.", err)
		} else {
			fmt.Println("Started EC2 instance:", result.StartingInstances)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
