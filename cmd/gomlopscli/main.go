package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/saiaunghlyanhtet/go-mlops/internal/gomlops"
	"github.com/spf13/cobra"
)

func main() {
	var modelPath, modelName, endpointName, endpointConfigName, modelArn, version string

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "gomlops",
		Short: "A CLI for deploying models with SageMaker",
		Long:  `A Command Line Interface for deploying models with SageMaker`,
	}

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a model",
		Run: func(cmd *cobra.Command, args []string) {
			smc := &gomlops.SageMakerClient{SMClient: sagemaker.NewFromConfig(cfg)}
			err := smc.DeployModel(modelPath, endpointName, endpointConfigName)
			if err != nil {
				fmt.Println("Error deploying model:", err)
				os.Exit(1)
			}
			fmt.Println("Model deployed successfully")
		},
	}

	deployCmd.Flags().StringVarP(&modelPath, "model", "m", "", "Path to the model")
	deployCmd.Flags().StringVarP(&endpointName, "endpoint", "e", "", "Name of the endpoint")
	deployCmd.Flags().StringVarP(&endpointConfigName, "config", "c", "", "Name of the endpoint configuration")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version a model",
		Run: func(cmd *cobra.Command, args []string) {
			smc := &gomlops.SageMakerClient{SMClient: sagemaker.NewFromConfig(cfg)}
			err := smc.VersionModel(context.TODO(), modelName, modelArn, version)
			if err != nil {
				fmt.Println("Error versioning model:", err)
				os.Exit(1)
			}
			fmt.Println("Model versioned successfully")
		},
	}

	versionCmd.Flags().StringVarP(&modelName, "name", "n", "", "Name of the model")
	versionCmd.Flags().StringVarP(&modelArn, "arn", "a", "", "ARN of the model")
	versionCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the model")

	var testCmd = &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Test command executed")
		},
	}

	rootCmd.AddCommand(deployCmd, versionCmd, testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
