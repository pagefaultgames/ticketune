package utils

import (
	"fmt"

	"github.com/amatsagu/tempest"
)

var ErrComponentNotFound = fmt.Errorf("component not found")
var ErrAttachmentNotFound = fmt.Errorf("attachment not found")

// Helper function to extract the component from a modal response's label
// Returns the component if found, (otherwise the zero value)
func GetLabelComponent[T tempest.StringSelectComponent | tempest.TextInputComponent](data *tempest.ModalInteractionData, expectedIndex int) T {
	var zero T
	if len(data.Components) <= expectedIndex {
		return zero
	}
	label, ok := data.Components[expectedIndex].(tempest.LabelComponent)
	if !ok {
		return zero
	}
	// Get the child component
	component, ok := label.Component.(T)
	if !ok {
		return zero
	}
	return component
}

// Extract the file attachment given the label of the file upload component corresponding to componentId
// NOTE: This will not work due to a bug in tempest's shape of interaction response data. Leaving commented out until a fix is
// pushed.
// func GetFileUploadAttachment(itx tempest.ModalInteraction, componentId string) ([]tempest.Attachment, error) {
// 	if itx.Data.Resolved == nil || itx.Data.Resolved.Attachments == nil {
// 		return tempest.Attachment{}, ErrAttachmentNotFound
// 	}
// 	var attachmentId tempest.Snowflake = 0
// 	// Iterate through...
// 	for _, component := range itx.Data.Components {
// 		label, ok := component.(tempest.LabelComponent)
// 		if !ok {
// 			continue
// 		}
// 		fileUpload, ok := label.Component.(tempest.FileUploadComponent)
// 		if !ok || fileUpload.CustomID != componentId {
// 			continue
// 		}
// 		// fileUpload.values
// 		attachmentId = tempest.Snowflake(fileUpload.Value)
// 		break
// 	}

// 	if attachmentId == 0 {
// 		return tempest.Attachment{}, ErrComponentNotFound
// 	}
// 	value, ok := itx.Data.Resolved.Attachments[attachmentId]
// 	if !ok {
// 		return tempest.Attachment{}, ErrAttachmentNotFound
// 	}
// 	return value, nil
// }
