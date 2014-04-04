builds := $(basename $(shell git grep -l "[p]ackage main"))

install: test $(builds)

test:
	go test

$(builds):
	cd $(dir $@) && go install
