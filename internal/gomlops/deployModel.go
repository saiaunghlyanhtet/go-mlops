package gomlops

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go/aws"
)

type SageMakerClientAPI interface {
	CreateEndpoint(ctx context.Context, params *sagemaker.CreateEndpointInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateEndpointOutput, error)
	AddTags(ctx context.Context, params *sagemaker.AddTagsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.AddTagsOutput, error)
}

type SageMakerClient struct {
	SMClient SageMakerClientAPI
}

func (smc *SageMakerClient) DeployModel(modelPath string, endpointName string, endpointConfigName string) error {
	// Deploy the model using SageMaker CreateEndpoint API
	_, err := smc.SMClient.CreateEndpoint(context.TODO(), &sagemaker.CreateEndpointInput{
		EndpointConfigName: aws.String(endpointConfigName),
		EndpointName:       aws.String(endpointName),
	})
	if err != nil {
		return fmt.Errorf("error deploying model: %v", err)
	}

	return nil
}
