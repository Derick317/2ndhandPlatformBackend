package service

import (
	"fmt"
	"mime/multipart"
	"secondHand/backend"
	"secondHand/model"
	"strconv"
)

// ITEM.ImageUrls should be empty.
func AddItem(item *model.Item, imageFiles []multipart.File) error {
	// Clear item's imageurls
	item.ImageUrls = make(model.ImageUrlsType)
	item.ImageUrls["nextKey"] = "0"

	// Save to database
	backend.CreateRecord(item)

	// Save images in GCS
	for _, imageFile := range imageFiles {
		fileName := strconv.Itoa(int(item.ID)) + ":" + item.ImageUrls["nextKey"]
		medialink, err := backend.SaveToGCS(imageFile, fileName)
		if err != nil {
			return err
		}
		item.ImageUrls[fileName] = medialink
		currentKey, err := strconv.Atoi(item.ImageUrls["nextKey"])
		if err != nil {
			return err
		}
		item.ImageUrls["nextKey"] = strconv.Itoa(currentKey + 1)
	}

	// Update record in the database
	numRowsAffected, err := backend.UpdateColumnWithConditions(item, "", nil, "", nil)
	if err != nil {
		return err
	}
	if numRowsAffected != 1 {
		return fmt.Errorf("changed %d records at AddItem", numRowsAffected)
	}
	return nil
}
