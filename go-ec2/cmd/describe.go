package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe EC2 instances",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(Region())},
		)

		// Create EC2 service client
		svc := ec2.New(sess)

		// Call to get detailed information on each instance
		result, err := svc.DescribeInstances(nil)
		if err != nil {
			fmt.Println("Error describing EC2 instances.", err)
		} else {
			fmt.Println("Instances:", result)
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
