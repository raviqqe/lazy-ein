#[link(name = "jsonxx", kind = "static")]
extern "C" {
    fn jsonxx_main();
}

fn main() {
    unsafe {
        jsonxx_main();
    };
}
