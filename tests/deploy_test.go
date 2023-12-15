package tests

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/saiaunghlyanhtet/go-mlops/internal/gomlops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SageMakerClientAPI interface {
	CreateEndpoint(ctx context.Context, params *sagemaker.CreateEndpointInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateEndpointOutput, error)
	AddTags(ctx context.Context, params *sagemaker.AddTagsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.AddTagsOutput, error)
}

type MockSageMakerClient struct {
	mock.Mock
}

func (m *MockSageMakerClient) CreateEndpoint(ctx context.Context, params *sagemaker.CreateEndpointInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateEndpointOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.CreateEndpointOutput), args.Error(1)
}

func (m *MockSageMakerClient) AddTags(ctx context.Context, params *sagemaker.AddTagsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.AddTagsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.AddTagsOutput), args.Error(1)
}

func TestDeployModel(t *testing.T) {
	mockClient := new(MockSageMakerClient)
	smc := &gomlops.SageMakerClient{SMClient: mockClient}

	mockClient.On("CreateEndpoint", mock.Anything, mock.Anything).Return(&sagemaker.CreateEndpointOutput{}, nil)

	err := smc.DeployModel("modelPath", "endpointName", "endpointConfigName")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	mockClient.AssertExpectations(t)
}

func TestVersionModel(t *testing.T) {
	mockClient := new(MockSageMakerClient)
	smc := &gomlops.SageMakerClient{SMClient: mockClient}

	modelArn := "arn:aws:sagemaker:us-west-2:123456789012:model/my-model"
	version := "v1"

	mockClient.On("AddTags", mock.Anything, &sagemaker.AddTagsInput{
		ResourceArn: aws.String(modelArn),
		Tags: []types.Tag{
			{
				Key:   aws.String("Version"),
				Value: aws.String(version),
			},
		},
	}).Return(&sagemaker.AddTagsOutput{}, nil)

	err := smc.VersionModel(context.TODO(), "hello", modelArn, version)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestDeployModelIntegration(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Fatalf("unable to load SDK config, %v", err)
	}

	smc := &gomlops.SageMakerClient{SMClient: sagemaker.NewFromConfig(cfg)}

	endpointName := "endpointName"
	err = smc.DeployModel("sampleModelPath", endpointName, "endpointConfigName")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err != nil {
		t.Errorf("Failed to delete endpoint: %v", err)
	}
}
