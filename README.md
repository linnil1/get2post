# GET to POST service in GoLang

Flow: GET (query) -> aws APIGateway -> aws Lambda -> [Lambda web adapter](https://github.com/awslabs/aws-lambda-web-adapter/tree/main) -> [GoLang Binary](./app) -> check secret -> POST (json) -> target web


## Deploy

``` bash
sam build
sam deploy --guide --parameter-overrides SECRET={{password}}
```

Then try it out with [httpbin](https://github.com/postmanlabs/httpbin):

``` bash
curl "https://{{your_apigateway_id}}.execute-api.{{your_region}}.amazonaws.com/get2post?secret={{password}}&url=https://httpbin.org/post&data.content=Test%20from%20aws&data.user.name=Testing&data.user.id=123"
```

The resulting POST body will be structured as follows:
``` json
{
  "content": "Test from aws",
  "user": {
    "id": "123",
    "name": "Testing"
}
```

Parameter Explanations:
* `secret`: The required password for activating the service. Leave it empty if `SECRET` is not provided.
* `url`: The destination URL for posting JSON data.
* `data.*`: Parameters for supplying flattened POST data. Refer to https://github.com/nqd/flat for syntax details.



## Environment
* GoLang 1.20 (alpine)
* Go package: gin
* AWS web adapater https://github.com/awslabs/aws-lambda-web-adapter/tree/main