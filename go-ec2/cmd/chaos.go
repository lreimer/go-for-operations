package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var count int

// chaosCmd represents the chaos command
var chaosCmd = &cobra.Command{
	Use:   "chaos",
	Short: "Terminate random tagged EC2 instances",
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
			return
		}

		// collect instances from reservations
		var instances []*ec2.Instance
		for _, r := range result.Reservations {
			instances = append(instances, r.Instances...)
		}

		for i, vm := range instances {
			if i >= count {
				break
			}

			input := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{
					aws.String(*vm.InstanceId),
				},
				DryRun: aws.Bool(false),
			}

			result, err := svc.TerminateInstances(input)

			if err != nil {
				fmt.Println("Error terminating EC2 instance.", err)
			} else {
				fmt.Println("Terminated EC2 instance:", result.TerminatingInstances[0].InstanceId)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chaosCmd)

	chaosCmd.Flags().IntVarP(&count, "count", "c", 1, "the EC2 chaos count")
}
