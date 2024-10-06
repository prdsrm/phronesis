## Lambda function

So, I'm rewriting it as a lambda function.
Following: <https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html>.

## Usage as a cloud function

NOTE: I implemented it, checkout 733806a852803b865d6a0810b9c2be962d59c497, but it doesn't work on DigitalOcean because... `gotd/td` build requires too much memory,
since its generated code, so the compilation process is killed. Can't really do anything about it. Cloud Functions support on DO is way too early.

Deploy this cloud function:
- Follow <https://docs.digitalocean.com/products/functions/how-to/create-functions/> carefully
- `cd ..`
- `doctl serverless deploy phronesis`

## Quick testing, deploy locally

Run `cd local && go run`
Connect:
`curl http://localhost/connect -d '{"user_id": 0}`
