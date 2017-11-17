package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"testing"
	"github.com/golang/protobuf/proto"
	"encoding/json"
	"github.com/google/uuid"
)

func TestGettingEvent (t *testing.T) {
	eventSaver := new(EventSaver)
	stub := createMockStub(t, eventSaver)

	txId := "mockTxID"
	testEvent := Event{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		"",
		1510830496,
		EventType_RETRIEVE,
		EventStatus_IN_PROGRESS,
	}

	protoTestEvent, _ := proto.Marshal(&testEvent);

	stub.MockInit("initTx", [][]byte{})

	stub.MockTransactionStart(txId)
	eventSaver.SaveEvent(stub, []string{string(protoTestEvent)})
	stub.MockTransactionStart(txId)

	stub.MockTransactionStart(txId)
	response := eventSaver.GetEvent(stub, []string{string(testEvent.Id)})
	stub.MockTransactionStart(txId)

	retrievedEvent := Event{}
	proto.Unmarshal(response.Payload, &retrievedEvent)

	assertEqual(t, retrievedEvent.SourceId, testEvent.SourceId, "")
	assertEqual(t, retrievedEvent.TargetId, testEvent.TargetId, "")
	assertEqual(t, retrievedEvent.PayloadId, testEvent.PayloadId, "")
	assertEqual(t, retrievedEvent.PayloadHash, testEvent.PayloadHash, "")
	assertEqual(t, retrievedEvent.EventType, testEvent.EventType, "")
	assertEqual(t, retrievedEvent.EventStatus, testEvent.EventStatus, "")
}


func TestSavingEvent (t *testing.T) {
	eventSaver := new(EventSaver)
	stub := createMockStub(t, eventSaver)

	txId := "mockTxID"
	testEvent := Event{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		"",
		1510830496,
		EventType_CREATE,
		EventStatus_FAILED,
	}

	protoTestEvent, _ := proto.Marshal(&testEvent);

	stub.MockInit("initTx", [][]byte{})

	stub.MockTransactionStart(txId)
	response := eventSaver.SaveEvent(stub, []string{string(protoTestEvent)})
	stub.MockTransactionStart(txId)

	if s := response.GetStatus(); s != 200 {
		t.Errorf("the status is %d, instead of 200", s)
		t.Errorf("message: %s", response.Message)
	}

	var savedEvent Event
	eventKey := "event:" + testEvent.Id
	json.Unmarshal(stub.State[eventKey], &savedEvent)

	assertEqual(t, savedEvent.SourceId, testEvent.SourceId, "")
	assertEqual(t, savedEvent.TargetId, testEvent.TargetId, "")
	assertEqual(t, savedEvent.PayloadId, testEvent.PayloadId, "")
	assertEqual(t, savedEvent.PayloadHash, testEvent.PayloadHash, "")
	assertEqual(t, savedEvent.EventType, testEvent.EventType, "")
	assertEqual(t, savedEvent.EventStatus, testEvent.EventStatus, "")
}

func createMockStub(t *testing.T, eventSaver *EventSaver) *shim.MockStub {
	stub := shim.NewMockStub("mockStub", eventSaver)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	return stub
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
