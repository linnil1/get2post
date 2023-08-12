# GET to POST service in GoLang

Flow: GET -> aws APIGateway -> aws Lambda -> [Lambda web adapter](https://github.com/awslabs/aws-lambda-web-adapter/tree/main) -> [GoLang Binary](./app) -> check secret -> POST -> target web


## Deploy

``` bash
sam build
sam deploy --guide --parameter-overrides APP_SECRET={{password}}
```

Then try it out:
(This example is to trigger Discord webhook by GET)

```
curl "https://{{your_apigateway_id}}.execute-api.us-east-2.amazonaws.com/get2post?secret={{password}}&url=https://discord.com/api/webhooks/{{your_webhook_id_and_token}}&content=Test%20from%20aws"
```


## Environment
* GoLang 1.20 (alpine)
* Go package: gin
* AWS web adapater https://github.com/awslabs/aws-lambda-web-adapter/tree/main