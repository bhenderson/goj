builds := $(basename $(shell git grep -l "[p]ackage main"))

install: test $(builds)

test:
	go test -timeout=1s -bench=.

$(builds):
	cd $(dir $@) && go install
