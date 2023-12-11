package service

import (
	"errors"
	"fmt"
	"log"
	"secondHand/model"
	"secondHand/util"
	"time"
)

var orderCancelChannel = make(chan model.Order)

func InitOrderCanceler() {
	go func() {
		for {
			order := <-orderCancelChannel
			now := time.Now().Unix()
			if order.ExpTime > now {
				time.Sleep((time.Duration)(order.ExpTime-now) * time.Second)
			}
			if err := CancelOrder(order); err != nil && !errors.Is(err, util.ErrOrderNoExists) {
				log.Printf("Auto-canceler failed to cancel order %d: %s\n", order.ID, err.Error())
			} else {
				fmt.Printf("Auto-canceler canceled order %d.\n", order.ID)
			}
		}
	}()
}

func AddToOrderCanceler(order model.Order) {
	orderCancelChannel <- order
}
