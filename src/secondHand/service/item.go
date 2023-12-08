package service

import (
	"fmt"
	"mime/multipart"
	"secondHand/backend"
	"secondHand/model"
)

// ITEM.ImageUrls should be empty.
func AddItem(item *model.Item, imageFiles []multipart.File) error {
	var err error = nil

	// Clear item's imageurls
	item.ImageNextKey = 0

	// Save to database
	if err = backend.CreateRecord(item); err != nil {
		return err
	}

	// Save images in GCS
	for _, imageFile := range imageFiles {
		fileName := fmt.Sprintf("%d-%d", item.ID, item.ImageNextKey)
		medialink, err := backend.SaveToGCS(imageFile, fileName)
		if err != nil {
			return err
		}
		item.ImageUrls.Add(fileName, medialink)
		item.ImageNextKey += 1
	}

	// Update record in the database
	numRowsAffected, err := backend.UpdateColumnsWithConditions(item, "", nil, nil)
	if err != nil {
		return err
	}
	if numRowsAffected != 1 {
		return fmt.Errorf("changed %d records at AddItem", numRowsAffected)
	}
	return nil
}

// DeleteItem deletes item's record from the database and its images from Google cloud storage
// Item's status should be "Deleted"
func DeleteItem(item *model.Item) error {
	if err := backend.DeleteFromDBByPrimaryKey(&model.Item{}, item.ID); err != nil {
		return err
	}
	for imageKey := range item.ImageUrls {
		if err := backend.DeleteFromGCS(imageKey); err != nil {
			return err
		}
	}
	return nil
}

// TestAndSetItemStatus changes the ID-th item's status to NEWSTATUS if its
// original status is TARGET; otherwise, does nothing. It returns its current
// status and whether update the status successfully and possible error.
//
// It reports unsuccess if
//   - some error happens;
//   - original status is not TARGET;
//   - other goroutine has changed the status.
//
// This function is atomic.
func TestAndSetItemStatus(id uint64, target model.ItemStatusType,
	newStatus model.ItemStatusType) (model.ItemStatusType, bool, error) {
	var item model.Item
	// Read initial status and version
	if err := backend.ReadFromDBByPrimaryKey(&item, id); err != nil {
		return model.StatusCounter, false, err
	}
	if item.Status != target {
		return item.Status, false, nil
	}

	// Set new status
	num, err := backend.UpdateColumnsWithConditions(&item, "version", item.Version,
		map[string]interface{}{"status": newStatus, "version": item.Version + 1})
	if err != nil {
		return model.StatusCounter, false, err
	}
	// Check success
	if num == 1 {
		return newStatus, true, nil
	}
	if err := backend.ReadFromDBByPrimaryKey(&item, id); err != nil {
		return model.StatusCounter, false, err
	}
	return item.Status, false, nil
}

// The record will be save in ITEM, so ITEM should be a pointer.
func QueryItem(item *model.Item, itemId uint64) error {
	return backend.ReadFromDBByPrimaryKey(item, itemId)
}
