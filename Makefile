fmt:
	python ./scripts/make.py fmt

tests:
	python ./scripts/make.py tests

tests-cover:
	python ./scripts/make.py tests:cover

lint:
	python ./scripts/make.py lint

security:
	python ./scripts/make.py security

generate:
	python ./scripts/make.py generate

all: fmt tests tests-cover lint security generate

.NOTPARALLEL:
.PHONY: all fmt tests tests-cover lint security generate
