rewards
=======

Rewards service.

Goals
-----

- Build a pure go service with as few services as possible
- Learn go packages - http & sql



Testing
-------

$ psql postgres -c "create db reward_test;"
$ DB_CONNECT=postgres://<username>:@localhost/reward_test?sslmode=disable go test github.com/davidoram/rewards/context
$ psql postgres -c "drop db reward_test;"
