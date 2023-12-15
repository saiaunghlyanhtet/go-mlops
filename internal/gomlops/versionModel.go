package gomlops

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/aws-sdk-go/aws"
)

func (smc *SageMakerClient) VersionModel(ctx context.Context, modelName string, modelArn string, version string) error {
	// Tag the model using SageMaker AddTags API
	_, err := smc.SMClient.AddTags(ctx, &sagemaker.AddTagsInput{
		ResourceArn: aws.String(modelArn),
		Tags: []types.Tag{
			{
				Key:   aws.String("Version"),
				Value: aws.String(version),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error tagging model: %v", err)
	}

	return nil
}
