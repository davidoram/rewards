rewards
=======

Rewards service.

Goals
-----

- Build a pure go service with as few services as possible
- Learn go packages - http & sql



Testing
-------

$ psql postgres -c "create database reward_test;"
$ export DB_CONNECT=postgres://davidoram:@localhost/reward_test?sslmode=disable; go test src/github.com/davidoram/rewards/context_test.go
$ psql postgres -c "drop database reward_test;"
