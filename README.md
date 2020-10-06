# go-event-sourcing

## Features
- Use HTTP protocol
- Persitent event to embbeded database (leveldb)
- Subscribe api using HTTP ServerSentEvent
  - Subscribe projection event: merge new_event with old_event using _domain_id and _domain_name
  - Subscribe aggregation event: Configurable group_by_key, aggregated function: min, max, count, last
  - Subscribe changed data capture (CDC): response domain event only change properties
- Support client group consumers
  - Same client_id will shared events (one event only bound to one consumer) across consumers
  - Unique client_id will get all events


## Architecture Diagram
![alt text](https://raw.githubusercontent.com/samuraiiway/go-event-sourcing/develop/go-event-store.png)

## Event Anatomy
### Internal Properties
- _event_id: auto generated
- _domain_name: user defined
- _domain_id: user defined or auto generated, used for projection
- _timestamp: auto generated
### Domain Properties
- field: value
- ...
- Only support value for string, number, boolean
- Not support nested json nor array
### Example Request/Response
POST http://localhost:8080/api/event/transaction
```json
{
    "_domain_id": "20200927-0000012345",
    "service_id": "transfer",
    "payer_id": "123450",
    "payee_id": "98765",
    "amount": 100,
    "status": "success"
}
```
```json
{
    "_domain_id": "20200927-0000012345",
    "_domain_name": "transaction",
    "_event_id": "163b42844ec238f8-ee06f19f93527f6",
    "_timestamp": 1601947228484818000,
    "amount": 100,
    "payee_id": "98765",
    "payer_id": "123450",
    "service_id": "transfer",
    "status": "success"
}
```

## Projection Anatomy
### Internal Properties
- _event_id: last event id
- _domain_name: event domain name
- _domain_id:  event domain id as primary key
- _timestamp: last event timestamp
- client_group: path variable at the end
### Example Response
GET http://localhost:8080/stream/projection/transaction/client_group_1
```json
{"_domain_id":"20200927-0000012345","_domain_name":"transaction","_event_id":"163b432f943d6230-972778c45d952007","_timestamp":1601947964089921000,"amount":100,"payee_id":"98765","payer_id":"123450","service_id":"transfer","status":"pending"}

{"_domain_id":"20200927-0000012345","_domain_name":"transaction","_event_id":"163b432fa4a067f8-c14ae2f7c8f48313","_timestamp":1601947964364847000,"amount":100,"payee_id":"98765","payer_id":"123450","service_id":"transfer","status":"confirm"}

{"_domain_id":"20200927-0000012345","_domain_name":"transaction","_event_id":"163b432fc1b97118-31d0c889289c0813","_timestamp":1601947964853026000,"amount":100,"payee_id":"98765","payer_id":"123450","service_id":"transfer","status":"success"}
```

## Aggregation Anatomy
### Internal Properties
- _event_id: last event id
- _domain_name: event domain name
- _aggregation_name: user defined
- _aggregation_key: group_by_key
- client_group: path variable at the end
### Example Configuration
```json
{
    "domain_name": "transaction",
    "aggregated_id": "transaction_service_summary",
    "group_by_key_id": ["service_id", "status"],
    "aggregated_function": [
        {
            "property_name": "amount",
            "field_name": "sum_amount",
            "function": "sum"
        },
        {
            "property_name": "amount",
            "field_name": "min_amount",
            "function": "min"
        },
        {
            "property_name": "amount",
            "field_name": "max_amount",
            "function": "max"
        },
        {
            "property_name": "amount",
            "field_name": "count_amount",
            "function": "count"
        },
        {
            "property_name": "payer_id",
            "field_name": "last_payer",
            "function": "last"
        }
    ]
}
```
### Example Response
GET http://localhost:8080/stream/aggregation/transaction/client_group_1
```json
{"_aggregation_key":":transfer:success","_aggregation_name":"transaction_service_summary","_domain_name":"transaction","_event_id":"163b436b54015f28-d24b1e9a5570c704","_timestamp":1601948220710284000,"count_amount":600005,"last_payer":"123450","max_amount":499999,"min_amount":0,"sum_amount":129999700500}

{"_aggregation_key":":transfer:success","_aggregation_name":"transaction_service_summary","_domain_name":"transaction","_event_id":"163b436c3f1359a8-93e13da04fecc833","_timestamp":1601948224654116000,"count_amount":600006,"last_payer":"123450","max_amount":499999,"min_amount":0,"sum_amount":129999700510}

{"_aggregation_key":":transfer:success","_aggregation_name":"transaction_service_summary","_domain_name":"transaction","_event_id":"163b436cba3f5cf8-76b07a4f849d3a6c","_timestamp":1601948226720592000,"count_amount":600007,"last_payer":"123450","max_amount":499999,"min_amount":0,"sum_amount":129999700511}

{"_aggregation_key":":transfer:success","_aggregation_name":"transaction_service_summary","_domain_name":"transaction","_event_id":"163b436da0e52f70-86c74811083b4eeb","_timestamp":1601948230590217000,"count_amount":600008,"last_payer":"123450","max_amount":499999,"min_amount":0,"sum_amount":129999701511}
```

## Changed Data Capture Anatomy
### Internal Properties
- _event_id: changed event id
- _timestamp: changed event timestamp
### Changed Properties
- field_new: new value of changed event
- field_old: old value of changed event

### Example Response
GET http://localhost:8080/stream/changed/transaction/client_group_1
```json
{"_event_id":"163b43c5170f4008-c5a00358c668398e","_timestamp":1601948606234840000,"status_new":"failed","status_old":"success"}

{"_event_id":"163b43c7d1d93428-ecfce96c00159c2a","_timestamp":1601948617958572000,"status_new":"canceled","status_old":"failed"}

{"_event_id":"163b43c9f19e9a30-8b0a67832d2a8035","_timestamp":1601948627081537000,"status_new":"refunded","status_old":"canceled"}
```
