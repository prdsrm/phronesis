## Usage as a cloud function

Deploy this cloud function:
- Follow <https://docs.digitalocean.com/products/functions/how-to/create-functions/> carefully
- `cd ..`
- `doctl serverless deploy phronesis`

## Quick testing, deploy locally

Run `cd local && go run`
Connect:
`curl http://localhost/connect -d '{"user_id": 0}`
