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
* Go package: gin, flat
* AWS web adapater https://github.com/awslabs/aws-lambda-web-adapter


## License
MIT License

Copyright (c) 2023 linnil1

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
