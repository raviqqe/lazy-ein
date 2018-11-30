libexecutable = executable/libexecutable.a
libio = io/target/release/libio.a

.PHONY: all
all: $(libexecutable) $(libio)

$(libexecutable):
	cd executable && cmake . && make

$(libio):
	cd io && cargo build --release

.PHONY: clean
clean:
	git clean -dfx
