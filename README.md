## Lambda function

So, I'm rewriting it as a lambda function.
Following: <https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html>.

## Usage as a cloud function

NOTE: it doesn't work on DigitalOcean because... `gotd/td` build requires too much memory,
since its generated code.

Deploy this cloud function:
- Follow <https://docs.digitalocean.com/products/functions/how-to/create-functions/> carefully
- `cd ..`
- `doctl serverless deploy phronesis`

## Quick testing, deploy locally

Run `cd local && go run`
Connect:
`curl http://localhost/connect -d '{"user_id": 0}`
