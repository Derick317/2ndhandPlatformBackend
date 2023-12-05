package service

import (
	"fmt"
	"mime/multipart"
	"secondHand/backend"
	"secondHand/constants"
	"secondHand/model"
	"strconv"
)

// ITEM.ImageUrls should be empty.
func AddItem(item *model.Item, imageFiles []multipart.File) error {
	// Clear item's imageurls
	item.ImageUrls.Add(constants.NEXTKEY_KEY, "0")

	// Save to database
	backend.CreateRecord(item)

	// Save images in GCS
	for _, imageFile := range imageFiles {
		fileName := strconv.FormatUint(item.ID, 10) + "-" + item.ImageUrls[constants.NEXTKEY_KEY]
		medialink, err := backend.SaveToGCS(imageFile, fileName)
		if err != nil {
			return err
		}
		item.ImageUrls.Add(fileName, medialink)
		currentKey, err := strconv.Atoi(item.ImageUrls[constants.NEXTKEY_KEY])
		if err != nil {
			return err
		}
		item.ImageUrls[constants.NEXTKEY_KEY] = strconv.Itoa(currentKey + 1)
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
		return model.TagCounter, false, err
	}
	if item.Status != target {
		return item.Status, false, nil
	}

	// Set new status
	num, err := backend.UpdateColumnsWithConditions(&item, "version", item.Version,
		map[string]interface{}{"status": newStatus, "version": item.Version + 1})
	if err != nil {
		return model.TagCounter, false, err
	}
	// Check success
	if num == 1 {
		return newStatus, true, nil
	}
	if err := backend.ReadFromDBByPrimaryKey(&item, id); err != nil {
		return model.TagCounter, false, err
	}
	return item.Status, false, nil
}

// The record will be save in ITEM, so ITEM should be a pointer.
func QueryItem(item *model.Item, itemId uint64) error {
	return backend.ReadFromDBByPrimaryKey(item, itemId)
}
