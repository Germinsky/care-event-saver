package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/golang/protobuf/proto"

	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
)

type EventSaver struct {
	logger *shim.ChaincodeLogger
	dispatcher Dispatcher
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

	var event Event
	json.Unmarshal(eventJsonBytes, &event)

	eventBytes, err := proto.Marshal(&event)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to marshall message: '%v'. Error: %v", event, err)
		p.logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}

	return shim.Success(eventBytes)
}

func (p *EventSaver) SaveEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	encodedEventByteString := args[0]
	p.logger.Infof("encodedEventByteString: %v", args[0])
	event, err := p.DecodeProtoByteString(encodedEventByteString)
	p.logger.Infof("event: %v", event)

	if err != nil {
		errMsg := fmt.Sprintf("Error while unmarshalling Event: %v", err.Error())
		p.logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}
	jsonEvent, err := json.Marshal(&event)

	eventKey := "event:" + event.Id

	err = stub.PutState(eventKey, jsonEvent)
	if err != nil {
		errMsg := fmt.Sprintf("Error while saving patient '%v'. Error: %v", event, err)
		p.logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}

	return shim.Success(jsonEvent)
}

func (p *EventSaver) DecodeProtoByteString(encodedEventByteString string) (*Event, error) {
	var err error

	event := Event{}
	err = proto.UnmarshalText(encodedEventByteString, &event)

	return &event, err
}

func main() {
	err := shim.Start(new(EventSaver))
	if err != nil {
		fmt.Printf("Error starting ScheduleChaincode: %s", err)
	}
}
