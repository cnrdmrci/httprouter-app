# httprouter-app

### Prepare Url Mathing in "main.go" file

```golang
var projectNameToProjectUrl map[string]string = map[string]string {
	"receiving-api": "https://receivingApiUrl",
}
```

### Do Request

```bash
curl --location 'http://127.0.0.1:8080/plans/1' --header 'api-name: receiving-api'
```

### Response

```text
Request  -> https://receivingApiUrl/plans/1
--------------------------------------------
Response ->
{
    "id" : 1 
}
```