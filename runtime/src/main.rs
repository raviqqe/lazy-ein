#[link(name = "jsonxx", kind = "static")]
extern "fastcall" {
    fn jsonxx_main();
}

fn main() {
    unsafe {
        jsonxx_main();
    };
}
