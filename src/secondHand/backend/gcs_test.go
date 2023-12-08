package backend

import "testing"

func TestDeleteFromGCS(t *testing.T) {
	InitGCSBackend()
	if err := DeleteFromGCS("1-0"); err != nil {
		t.Errorf("Unexpect error in CreateRecord: %v", err)
	}
}
