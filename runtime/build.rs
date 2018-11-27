extern crate cc;

fn main() {
    cc::Build::new()
        .object("src/jsonxx.o")
        .compile("libjsonxx.a");
}
