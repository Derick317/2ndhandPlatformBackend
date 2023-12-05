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

func TestStatusTestAndSetOne(t *testing.T) {
	var item1 = model.Item{SellerId: 10, Price: 10.4, Tag: 1, Description: "Total football"}
	backend.InitPostgreSQLBackend()
	err := AddItem(&item1, []multipart.File{})
	if err != nil {
		t.Errorf("Unexpect error when adding item: %v", err)
	}
	newStatus, ok, err := TestAndSetItemStatus(1, model.ItemStatusType(model.Available),
		model.ItemStatusType(model.OnOrder))
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if !ok {
		t.Errorf("Failed to change status, current status is: %v", newStatus)
	}
	if newStatus != model.ItemStatusType(model.OnOrder) {
		t.Errorf("Expect new status is %v, but it is %v", model.OnOrder, newStatus)
	}
}

func TestStatusTestAndSetMultiple(t *testing.T) {
	const NUM_ROUTINE = 100
	backend.InitPostgreSQLBackend()
	var item = model.Item{SellerId: 10, Price: 10.4, Tag: 1, Description: "Total football"}
	err := AddItem(&item, []multipart.File{})
	if err != nil {
		t.Errorf("Unexpect error when adding item: %v", err)
	}
	ch := make(chan bool)
	for idx := 0; idx < NUM_ROUTINE; idx++ {
		go func() {
			newStatus, ok, err := TestAndSetItemStatus(1, model.ItemStatusType(model.Available),
				model.ItemStatusType(model.OnOrder))
			if err != nil {
				t.Errorf("Unexpect error: %v", err)
			}
			if ok && newStatus != model.ItemStatusType(model.OnOrder) {
				t.Errorf("Expect new status is %v, but it is %v", model.OnOrder, newStatus)
			}
			ch <- ok
		}()
	}

	num_ok := 0
	for idx := 0; idx < NUM_ROUTINE; idx++ {
		if <-ch {
			num_ok++
		}
	}
	if num_ok != 1 {
		t.Errorf("Unexpected num_ok: %d", num_ok)
	}
	if err := backend.ReadFromDBByPrimaryKey(&item, 1); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if item.Version != 1 || item.Status != model.ItemStatusType(model.OnOrder) {
		t.Errorf("Unexpect item: %v", item)
	}
}

func TestQueryItem(t *testing.T) {
	var item1 = model.Item{SellerId: 10, Price: 10.4, Tag: 1, Description: "Total football"}
	var item2 = model.Item{SellerId: 2, Price: 0.5, Tag: 0, Description: "Baseball"}
	backend.InitPostgreSQLBackend()
	if err := AddItem(&item1, []multipart.File{}); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := AddItem(&item2, []multipart.File{}); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	var item0 model.Item
	if err := QueryItem(&item0, 1); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if item0.SellerId != item1.SellerId {
		t.Errorf("Unexpect seller ID: %d", item0.SellerId)
	}
	if item0.Tag != item1.Tag {
		t.Errorf("Unexpect tag: %d", item0.SellerId)
	}
	if item0.Description != item1.Description {
		t.Errorf("Unexpect description: %s", item0.Description)
	}
}
