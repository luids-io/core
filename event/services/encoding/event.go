package encoding

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"

	"github.com/luids-io/core/event"
	pb "github.com/luids-io/core/protogen/eventpb"
)

// Event returns a new Event object from protobuf
func Event(pbevent *pb.Event) (event.Event, error) {
	e := event.Event{}
	e.ID = pbevent.GetId()
	e.Type = event.Type(pbevent.GetType())
	e.Code = event.Code(pbevent.GetCode())
	e.Level = event.Level(pbevent.GetLevel())
	e.Timestamp, _ = ptypes.Timestamp(pbevent.GetTimestamp())

	pbsource := pbevent.GetSource()
	if pbsource == nil {
		return event.Event{}, errors.New("source is empty")
	}
	e.Source.Hostname = pbsource.GetHostname()
	e.Source.Instance = pbsource.GetInstance()
	e.Source.Program = pbsource.GetProgram()
	//decode event data
	pbdata := pbevent.GetData()
	if pbdata == nil {
		return event.Event{}, errors.New("data is empty")
	}
	switch pbdata.GetDataEnc() {
	case pb.Event_Data_JSON:
		rawData := pbdata.GetData()
		err := json.Unmarshal(rawData, &e.Data)
		if err != nil {
			return event.Event{}, fmt.Errorf("unmarshalling data: %v", err)
		}
	case pb.Event_Data_NODATA:
		e.Data = make(map[string]interface{})
	}
	return e, nil
}

// EventPB returns a new protobuf event
func EventPB(e event.Event) (*pb.Event, error) {
	pbevent := &pb.Event{}
	pbevent.Id = e.ID
	pbevent.Type = pb.Event_Type(e.Type)
	pbevent.Code = int32(e.Code)
	pbevent.Level = pb.Event_Level(e.Level)
	pbevent.Timestamp, _ = ptypes.TimestampProto(e.Timestamp)
	pbevent.Source = &pb.Event_Source{
		Hostname: e.Source.Hostname,
		Program:  e.Source.Program,
		Instance: e.Source.Instance,
	}
	//if no data
	if e.Data == nil || len(e.Data) == 0 {
		pbevent.Data = &pb.Event_Data{
			DataEnc: pb.Event_Data_NODATA,
		}
		return pbevent, nil
	}
	// encode data to json
	encoded, err := json.Marshal(e.Data)
	if err != nil {
		return nil, fmt.Errorf("cannot encode data to json: %v", err)
	}
	pbevent.Data = &pb.Event_Data{
		DataEnc: pb.Event_Data_JSON,
		Data:    encoded,
	}
	return pbevent, nil
}
