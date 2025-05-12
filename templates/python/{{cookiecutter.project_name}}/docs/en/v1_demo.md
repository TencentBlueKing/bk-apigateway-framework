### Functional description

Example

### Request Parameters

#### Description

| Field      | Type      | Required | Description      |
|--------------------|------------|--------|------------|
| a_list           |  list    | 否     | 示例 |
| a_string            |  string    | 是     | 示例|
| a_bool |  bool      | 否     | 示例 |

#### Example

```json
{
    "a_list": ["1234567890"],
    "a_string": "Welcome to Blueking",
    "a_bool": true
}
```

### Response Result

#### Description

| Field      | Type      | Description      |
|-----------|----------|-----------|
|  data     |    dict    |       Result data, detailed information please refer to the following description     |

data

| Field      | Type      | Description      |
|-----------|----------|-----------|
|  summary   |    dict    |      Summary information     |
|  message    |    string  |      Failed reason     |
|  details    |    dict  |      Detailed information     |

#### Example

status == 200

```json
{
  "data": {
    "summary": {
      "total": 1,
      "succeeded": 1,
      "failed": 0
    },
    "message": "",
    "details": {}
  }
}
```

status == 500

```json
{
  "error": {
    "code": "FAILED",
    "message": "reason",
    "data": {
    },
    "details": []
  }
}
```
