package service

import (
	"mime/multipart"
	"secondHand/backend"
	"secondHand/model"
	"testing"
)

func TestAddItemGoodButNoImage(t *testing.T) {
	var item = model.Item{SellerId: 10, Price: 10.4, Tag: 1, Description: "Total football"}
	backend.InitPostgreSQLBackend()
	err := AddItem(&item, []multipart.File{})
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
}

func TestAddItemTwoItemsNoImage(t *testing.T) {
	var item1 = model.Item{SellerId: 10, Price: 10.4, Tag: 1, Description: "Total football"}
	var item2 = model.Item{SellerId: 2, Price: 0.5, Tag: 0, Description: "Baseball"}
	backend.InitPostgreSQLBackend()
	err := AddItem(&item1, []multipart.File{})
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	err = AddItem(&item2, []multipart.File{})
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
}
