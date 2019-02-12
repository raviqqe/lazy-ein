.PHONY: all
all: libexecutable libruntime

libexecutable:
	cd executable && cmake . && make

libruntime:
	cd runtime && cargo build --release

.PHONY: clean
clean:
	git clean -dfx
