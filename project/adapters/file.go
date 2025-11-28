package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

type FileAPIClient struct {
	clients *clients.Clients
}

func NewFileAPIClient(clients *clients.Clients) *FileAPIClient {
	if clients == nil {
		panic("NewFilesApiClient: clients is nil")
	}

	return &FileAPIClient{
		clients: clients,
	}
}

func (c *FileAPIClient) UploadFile(ctx context.Context, fileID, fileContent string) error {
	resp, err := c.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, fileContent)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %w", fileID, err)
	}

	if resp.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).With("file", fileID).Info("file already exists")
		return nil
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code while uploading file %s: %d", fileID, resp.StatusCode())
	}

	return nil
}
