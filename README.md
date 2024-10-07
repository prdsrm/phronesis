## Usage

Run `cd local && go run`
Get some bot in the database:`psql $DATABASE_URL`, and `SELECT * FROM bots ORDER BY RANDOM() LIMIT 1`
Connect: `curl http://localhost/connect -d '{"user_id": 0}`
