.PHONY: all
all: libexecutable libio

libexecutable:
	cd executable && cmake . && make

libio:
	cd io && cargo build --release

.PHONY: clean
clean:
	git clean -dfx
