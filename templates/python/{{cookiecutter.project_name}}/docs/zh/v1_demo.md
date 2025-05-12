### 功能描述

文档示例

### 请求参数

#### 说明

| 字段               |  类型      | 必选   |  描述      |
|--------------------|------------|--------|------------|
| a_list           |  list    | 否     | 示例 |
| a_string            |  string    | 是     | 示例|
| a_bool |  bool      | 否     | 示例 |

#### 示例

```json
{
    "a_list": ["1234567890"],
    "a_string": "Welcome to Blueking",
    "a_bool": true
}
```

### 响应结果

#### 说明

status 为 200 表示成功, 其他状态码表示失败


| 字段      | 类型      | 描述      |
|-----------|----------|-----------|
|  data     |    dict    |      结果数据，详细信息请见下面说明     |

data

| 字段      | 类型      | 描述      |
|-----------|----------|-----------|
|  summary   |    dict    |      汇总信息     |
|  message    |    string  |     失败原因     |
|  details    |    dict  |     详细信息     |

#### 示例

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
