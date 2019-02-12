.PHONY: all
all: libexecutable libio libcore

libexecutable:
	cd executable && cmake . && make

libio:
	cd io && cargo build --release

libcore:
	cd core && cargo build --release

.PHONY: clean
clean:
	git clean -dfx
