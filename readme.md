## Golang Exponea SDK

### Installation
```bash
go get github.com/flowup/exponea-go
```

### Usage
There are two types of requests available on the Exponea API: **Events**
and **Customers**. While **Events** are for tracking what happened in the real time, **Customers**
customizes attributes for the given set of customer ids.

At first, initialization of the API client is needed.
```go
client := exponea.NewClient("your-project-id")
```

> Please note that any call to the API will return the response with occured
errors (string array) and error which may be caused by the network or serializer.

#### Tracking Events
Tracking events is done by calling the `SendEvent` method on the `exponea.Client`.

```go
resp, err := client.SendEvent(&exponea.Event{
  Customers: map[string]string{
    "registered": "peter.malina@flowup.eu",
  },
  Type: "Something happened",
  Properties: map[string]string{
      "property": "and it's value",
  },
})
```

#### Customizing Customer values
Customization customer's attributes can be done by simply calling `SendCustomer` method.

```go
resp, err := client.SendCustomer(&exponea.Customer{
  Customers: map[string]string{
    "registered": "peter.malina@flowup.eu",
  },
  Properties: map[string]string{
    "property": "and it's valueee",
  },
})
```