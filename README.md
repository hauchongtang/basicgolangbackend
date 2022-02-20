# Instructions
- To run, simply clone the repo and execute `go run .`
- Head to https://localhost:10000/events or which ever domain preferred.
- This is only an API for consumption (i.e. only GET queries)

# Example query `GET`
```
[
  {
    "stopname": "Hall 12 \u002613",
    "data": [
      {
        "bus": {
          "TYPE": "RED",
          "id": 31516
        },
        "coordinates": [
          "1.354883000000",
          "103.687518000000"
        ],
        "arrive_in": 317.99311714285716
      },
      {
        "bus": {
          "TYPE": "RED",
          "id": 31509
        },
        "coordinates": [
          "1.340124000000",
          "103.682122000000"
        ],
        "arrive_in": 435.8393116401205
      }
    ],
    "coordinates": [
      "1.3516864875",
      "103.6806499958"
    ]
  },
  {
    "stopname": "Saraca Hall",
    "data": [
      {
        "bus": {
          "TYPE": "RED",
          "id": 31516
        },
        "coordinates": [
          "1.354883000000",
          "103.687518000000"
        ],
        "arrive_in": 134.13240000000002
      },
      {
        "bus": {
          "TYPE": "RED",
          "id": 31509
        },
        "coordinates": [
          "1.340124000000",
          "103.682122000000"
        ],
        "arrive_in": 356.13303502367626
      }
    ],
    "coordinates": [
      "1.3548398855",
      "103.6845606565"
    ]
  },
  ]
```
