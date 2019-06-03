# Scheduler API

Schedule REST API used to manage schedules and appointments.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go 1.12.5](https://golang.org/dl/)
- [Ginko/Gomega](https://github.com/onsi/ginkgo#set-me-up) used to run tests

### Installing/Running

#### No Container

Clone repository:
```
git clone https://github.com/ckaminer/scheduler-api.git
```

In your terminal, set the port you wish to run the server on (default is 8080):
```
export PORT=8080
```

Retrieve dependencies (from the project root):
```
go build
```

Start server (from the project root):
```
go run main.go
```

Alternatively:
```
go build
./scheduler-api
```

#### Using Docker
Clone repository:
```
git clone https://github.com/ckaminer/scheduler-api.git
```

Build Docker Image Locally (from the project root):
```
docker build -t scheduler-api .
```

Run Container:
```
docker run -p 8080:8080 -it scheduler-api
```

## Running the tests

Unit tests (from the project root):
```
ginkgo -r -skipPackage acceptance_tests
```

To run the acceptance tests, first start your server. Then in a new tab run:
```
cd {projectRoot}/acceptance_tests
ginkgo
```

Entire suite (from the project root):
```
ginkgo -r
```

## Endpoints

#### Create Schedule
`POST /schedules`

Sample Request Body:
```
{
  "owner_name": "Tyrion Lannister"
}
```

Expected Response:
```
{
  "id": 1,
  "owner_name": "Tyrion Lannister",
  "appointments": []
}
```

#### View Schedule
`GET /schedules/{scheduleID}`

Expected Response:
```
{
  "id": 1,
  "owner_name": "Tyrion Lannister",
  "appointments": [
    {
      "id": 8,
      "schedule_id": 1,
      "start_time": 5,
      "end_time": 9
    }
  ]
}
```

#### Delete Schedule
`DELETE /schedules/{scheduleID}`

Expected Response:
```
{
  "id": 1,
  "owner_name": "Tyrion Lannister",
  "appointments": [
    {
      "id": 8,
      "schedule_id": 1,
      "start_time": 5,
      "end_time": 9
    }
  ]
}
```

#### Create Appointment
`POST /schedules/{scheduleID}/appointments`

Sample Request Body:
```
{
  "start_time": 5,
  "end_time": 8
}
```

Expected Response:
```
{
  "id": 9,
  "schedule_id": 4,
  "start_time": 5,
  "end_time": 8
}
```

#### View Appointment
`GET /schedules/{scheduleID}/appointments/{appointmentID}`

Expected Response:
```
{
  "id": 9,
  "schedule_id": 4,
  "start_time": 5,
  "end_time": 8
}
```

#### Delete Appointment
`DELETE /schedules/{scheduleID}/appointments/{appointmentID}`

Expected Response:
```
{
  "id": 9,
  "schedule_id": 4,
  "start_time": 5,
  "end_time": 8
}
```
