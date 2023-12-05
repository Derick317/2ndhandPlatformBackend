package backend

import (
	"secondHand/model"
	"testing"
)

func TestReadFromDBByKeyFirstMatch(t *testing.T) {
	InitPostgreSQLBackend()
	apple := model.Item{SellerId: 1, Price: 1.5, Description: "apple"}
	if err := CreateRecord(&apple); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}
	var destItem model.Item
	if err := ReadFromDBByKey(&destItem, "seller_id", "1", true); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}

	if destItem.SellerId != apple.SellerId || destItem.Price != apple.Price || destItem.Description != apple.Description {
		t.Errorf("DestItem should be the same as apple, but it is %v", destItem)
	}
}

func TestReadFromDBByKeyAllMatch(t *testing.T) {
	InitPostgreSQLBackend()
	apple := model.Item{SellerId: 1, Price: 1.5, Description: "apple"}
	pineapple := model.Item{SellerId: 1, Price: 13.2, Description: "apple"}
	peach := model.Item{SellerId: 2, Description: "peach"}
	banana := model.Item{SellerId: 1, Description: "banana"}
	if err := CreateRecord(&apple); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}
	if err := CreateRecord(&pineapple); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}
	if err := CreateRecord(&banana); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}
	if err := CreateRecord(&peach); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}

	var destItems []model.Item
	if err := ReadFromDBByKeys(&destItems, []string{"seller_id", "description"},
		[]string{"1", "apple"}, false); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if destItems[0].SellerId != 1 && destItems[1].SellerId != 1 {
		t.Errorf("Incorrect destItem: %v", destItems)
	}
}

func TestDeleteFromDBByKeysOneRecord(t *testing.T) {
	InitPostgreSQLBackend()
	order := model.Order{ItemId: 1, BuyerId: 1, ExpTime: 10}
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	order.ItemId = 2
	order.ID = 0
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}

	num, err := DeleteFromDBByKeys(&model.Order{}, []string{"item_id", "buyer_id"}, []string{"1", "1"})
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if num != 1 {
		t.Errorf("Unexpect num: %d", num)
	}
}

func TestDeleteFromDBByKeysTwoRecords(t *testing.T) {
	InitPostgreSQLBackend()
	order := model.Order{ItemId: 1, BuyerId: 1, ExpTime: 10}
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	order.ItemId = 2 // {ItemId: 2, BuyerId: 1, ExpTime: 10}
	order.ID = 0
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	order.BuyerId = 2 // {ItemId: 2, BuyerId: 2, ExpTime: 10}
	order.ID = 0
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}

	num, err := DeleteFromDBByKeys(&model.Order{}, []string{"buyer_id"}, []string{"1"})
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if num != 2 {
		t.Errorf("Unexpect num: %d", num)
	}
}

func TestDeleteFromDBByKeyTwoRecords(t *testing.T) {
	InitPostgreSQLBackend()
	order := model.Order{ItemId: 1, BuyerId: 1, ExpTime: 10}
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	order.ItemId = 2 // {ItemId: 2, BuyerId: 1, ExpTime: 10}
	order.ID = 0
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	order.BuyerId = 2 // {ItemId: 2, BuyerId: 2, ExpTime: 10}
	order.ID = 0
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}

	num, err := DeleteFromDBByKey(&model.Order{}, "buyer_id", "1")
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if num != 2 {
		t.Errorf("Unexpect num: %d", num)
	}
}

func TestDeleteFromDBByKeyNoRecord(t *testing.T) {
	InitPostgreSQLBackend()
	order := model.Order{ItemId: 1, BuyerId: 1, ExpTime: 10}
	if err := CreateRecord(&order); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}

	num, err := DeleteFromDBByKey(&model.Order{}, "buyer_id", "2")
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if num != 0 {
		t.Errorf("Unexpect num: %d", num)
	}
}
