package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var version string
var commit string

var versionFlag bool
var cfgFile string

var replicas int
var namespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-chaos [deployment/name]",
	Short: "Perform chaos on your K8s deployments",
	Long:  "Delete specified number of Pods for speficied deployments",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("kubectl-chaos %v %v\n", version, commit)
			return
		}

		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

		config, err := kubeConfig.ClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Perform chaos on Kubernetes namespace", namespace)

		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), args[0], metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		labelMap := deployment.Spec.Selector.MatchLabels

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelMap).String(),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for i, p := range pods.Items {
			if i >= replicas {
				break
			}
			fmt.Println("Deleting pod", p.Name)
			clientset.CoreV1().Pods(namespace).Delete(context.TODO(), p.Name, metav1.DeleteOptions{})
		}
	},
}

// SetVersion sets the application version for the root command
func SetVersion(v string) {
	version = v
}

// SetCommit sets the commit hash for the root command
func SetCommit(c string) {
	commit = c
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-chaos.yaml)")

	rootCmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace to use")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Display version info")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kubectl-chaos" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubectl-chaos")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
