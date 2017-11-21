package storage

import (
	"encoding/json"
)

// Mock elasticsearch driver with static data
type Mock struct{}

func (m Mock) GetEvents(filter *Filter, tenantId string) ([]*EventDetail, int, error) {
	var detailedEvents eventListWithTotal
	json.Unmarshal(mockEvents, &detailedEvents)

	var events []*EventDetail

	for i := range detailedEvents.Events {
		events = append(events, &detailedEvents.Events[i])
	}

	return events, detailedEvents.Total, nil
}

func (m Mock) GetEvent(eventId string, tenantId string) (*EventDetail, error) {
	var parsedEvent EventDetail
	err := json.Unmarshal(mockEvent, &parsedEvent)
	return &parsedEvent, err
}

func (m Mock) MaxLimit() uint {
	return 100
}

func (m Mock) GetAttributes(queryName string, tenantId string) ([]string, error) {
	var parsedAttribute []string
	err := json.Unmarshal(mockEvent, parsedAttribute)
	return parsedAttribute, err
}

var mockEvent = []byte(`
{
  "source": "service/security",
  "event_id": "7be6c4ff-b761-5f1f-b234-f5d41616c2cd",
  "event_type": "create/role_assignment",
  "event_time": "2017-11-17T08:53:32.667973+0000",
  "resource_id": "f1a7118aee7698ab43deb080df40e01845127240e11bae64293837145a4a7dac",
  "resource_type": "service/security/account/user",
  "id": "7be6c4ff-b761-5f1f-b234-f5d41616c2cd",
  "eventTime": "2017-11-17T08:53:32.667973+0000",
  "action": "create/role_assignment",
  "outcome": "success",
  "initiator": {
	"typeURI": "service/security/account/user",
	"id": "5d847cb1e75047a29aa9dee2cabcce9b"
  },
  "target": {
	"typeURI": "service/security/account/user",
	"id": "f1a7118aee7698ab43deb080df40e01845127240e11bae64293837145a4a7dac"
  },
  "observer": {
	"typeURI": "service/security",
	"id": "a02d5699-4967-522f-8092-c286aea2deab"
  }
}
`)

var mockEvents = []byte(`
{
  "events": [
    {
      "source": "service/security",
      "event_id": "7be6c4ff-b761-5f1f-b234-f5d41616c2cd",
      "event_type": "create/role_assignment",
      "event_time": "2017-11-17T08:53:32.667973+0000",
      "resource_id": "f1a7118aee7698ab43deb080df40e01845127240e11bae64293837145a4a7dac",
      "resource_type": "service/security/account/user",
      "id": "7be6c4ff-b761-5f1f-b234-f5d41616c2cd",
      "eventTime": "2017-11-17T08:53:32.667973+0000",
      "action": "create/role_assignment",
      "outcome": "success",
      "initiator": {
        "typeURI": "service/security/account/user",
        "id": "5d847cb1e75047a29aa9dee2cabcce9b"
      },
      "target": {
        "typeURI": "service/security/account/user",
        "id": "f1a7118aee7698ab43deb080df40e01845127240e11bae64293837145a4a7dac"
      },
      "observer": {
        "typeURI": "service/security",
        "id": "a02d5699-4967-522f-8092-c286aea2deab"
      }
    },
    {
      "source": "service/security",
      "event_id": "f6f0ebf3-bf59-553a-9e38-788f714ccc46",
      "event_type": "create/role_assignment",
      "event_time": "2017-11-07T11:46:19.448565+0000",
      "resource_id": "ba2cc58797d91dc126cc5849e5d802880bb6b01dfd3013a35392ce00ae3b0f43",
      "resource_type": "service/security/account/user",
      "id": "f6f0ebf3-bf59-553a-9e38-788f714ccc46",
      "eventTime": "2017-11-07T11:46:19.448565+0000",
      "action": "create/role_assignment",
      "outcome": "success",
      "initiator": {
        "typeURI": "service/security/account/user",
        "id": "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812"
      },
      "target": {
        "typeURI": "service/security/account/user",
        "id": "ba2cc58797d91dc126cc5849e5d802880bb6b01dfd3013a35392ce00ae3b0f43"
      },
      "observer": {
        "typeURI": "service/security",
        "id": "b54da470-046c-539d-a921-dfa91b32f525"
      }
    },
    {
      "source": "service/security",
      "event_id": "eae03aad-86ab-574e-b428-f9dd58e5a715",
      "event_type": "create/role_assignment",
      "event_time": "2017-11-06T10:15:56.984390+0000",
      "resource_id": "c4d3626f405b99f395a1c581ed630b2d40be8b9701f95f7b8f5b1e2cf2d72c1b",
      "resource_type": "service/security/account/user",
      "id": "eae03aad-86ab-574e-b428-f9dd58e5a715",
      "eventTime": "2017-11-06T10:15:56.984390+0000",
      "action": "create/role_assignment",
      "outcome": "success",
      "initiator": {
        "typeURI": "service/security/account/user",
        "id": "21ff350bc75824262c60adfc58b7fd4a7349120b43a990c2888e6b0b88af6398"
      },
      "target": {
        "typeURI": "service/security/account/user",
        "id": "c4d3626f405b99f395a1c581ed630b2d40be8b9701f95f7b8f5b1e2cf2d72c1b"
      },
      "observer": {
        "typeURI": "service/security",
        "id": "9a3e952c-90a3-544d-9d56-c721e7284e1c"
      }
    },
    {
      "source": "service/security",
      "event_id": "49e2084a-b81c-51f1-9822-78cdd31d0944",
      "event_type": "create/role_assignment",
      "event_time": "2017-11-06T10:11:21.605421+0000",
      "resource_id": "c4d3626f405b99f395a1c581ed630b2d40be8b9701f95f7b8f5b1e2cf2d72c1b",
      "resource_type": "service/security/account/user",
      "id": "49e2084a-b81c-51f1-9822-78cdd31d0944",
      "eventTime": "2017-11-06T10:11:21.605421+0000",
      "action": "create/role_assignment",
      "outcome": "success",
      "initiator": {
        "typeURI": "service/security/account/user",
        "id": "21ff350bc75824262c60adfc58b7fd4a7349120b43a990c2888e6b0b88af6398"
      },
      "target": {
        "typeURI": "service/security/account/user",
        "id": "c4d3626f405b99f395a1c581ed630b2d40be8b9701f95f7b8f5b1e2cf2d72c1b"
      },
      "observer": {
        "typeURI": "service/security",
        "id": "6d4828eb-e497-5649-be10-f29d1ddb0977"
      }
    }
  ],
  "total": 4
}
`)
