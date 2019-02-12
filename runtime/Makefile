.PHONY: all
all: libio libcore

libio:
	cd io && cargo build --release

libcore:
	cd core && cargo build --release

.PHONY: clean
clean:
	git clean -dfx
