package cloudvision

import (
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/workshopapps/pictureminer.api/internal/model"
)

// ProcessImages processes images in the background using google cloud vision API.
// It has not been integrated with the code due to an issue with setting up a valid
// billing account on google cloud.
func ProcessImages(details *model.ProcessBatchAPIResponse) (*visionpb.BatchAnnotateImagesResponse, error) {
	ctx := context.TODO()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	request := make([]*visionpb.AnnotateImageRequest, 0, len(details.Images))
	for _, imageURL := range details.Images {
		req := &visionpb.AnnotateImageRequest{
			Image: vision.NewImageFromURI(imageURL),
			Features: []*visionpb.Feature{
				{Type: visionpb.Feature_LABEL_DETECTION, MaxResults: 5},
				{Type: visionpb.Feature_LANDMARK_DETECTION, MaxResults: 5},
			},
		}

		request = append(request, req)
	}

	response, err := client.BatchAnnotateImages(ctx, &visionpb.BatchAnnotateImagesRequest{
		Requests: request,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
