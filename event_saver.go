package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/proto"

	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
)

type EventSaver struct {
	logger *shim.ChaincodeLogger
	dispatcher Dispatcher
}

type Event struct {
	Id string `json:"id"`
}

func (p *EventSaver) Init(stub shim.ChaincodeStubInterface) pb.Response {
	p.logger = shim.NewLogger("event_saver")
	p.dispatcher = NewDispatcher()

	p.dispatcher.AddMapping(Functions_GET_EVENT.String(), p.GetEvent)
	p.dispatcher.AddMapping(Functions_SAVE_EVENT.String(), p.SaveEvent)

	p.logger.Infof("< Init call >")

	return shim.Success(nil)
}

func (p *EventSaver) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	functionName, args := stub.GetFunctionAndParameters()

	p.logger.Infof("=====================================================")
	p.logger.Infof("Invoking functionName %v with args %v", functionName, args)

	return p.dispatcher.Dispatch(stub)
}

func (p *EventSaver) GetEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	eventId := args[0]
	p.logger.Infof("eventId:%v", eventId)

	eventKey := "event:" + eventId
	eventJsonBytes, err := stub.GetState(eventKey)
	if err != nil {
		errMsg := fmt.Sprintf("Error while getting patient with id '%v'. Error: %v", eventId, err)
		p.logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}
	p.logger.Infof("Getting event [%v] \n", string(eventJsonBytes))

	return shim.Success(eventJsonBytes)
}

func (p *EventSaver) SaveEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	jsonEvent := args[0]

	event := Event{}
	json.Unmarshal([]byte(jsonEvent), &event)

	eventKey := "event:" + event.Id

	p.logger.Infof("Putting event [%v] with key [%v] \n", jsonEvent, eventKey)
	err := stub.PutState(eventKey, []byte(jsonEvent))
	if err != nil {
		errMsg := fmt.Sprintf("Error while saving patient '%v'. Error: %v", event, err)
		p.logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}

	return shim.Success([]byte(jsonEvent))
}

//func (p *EventSaver) DecodeProtoByteString(encodedEventByteString string) (*Event, error) {
//	var err error
//
//	event := Event{}
//	err = proto.UnmarshalText(encodedEventByteString, &event)
//
//	return &event, err
//}

func main() {
	err := shim.Start(new(EventSaver))
	if err != nil {
		fmt.Printf("Error starting ScheduleChaincode: %s", err)
	}
}
