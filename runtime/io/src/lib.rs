extern crate gc;

#[global_allocator]
static GLOBAL: gc::Allocator = gc::Allocator;

#[repr(C)]
pub struct Closure<E, P> {
    pub entry: E,
    pub payload: P,
}
#[repr(C)]
pub struct Payload;

pub type Input = Closure<extern "fastcall" fn(&Payload) -> f64, Payload>;
pub type Output = Closure<extern "fastcall" fn(&Payload) -> f64, Payload>;
pub type Main = Closure<extern "fastcall" fn(&Payload, &Input) -> &'static Output, Payload>;

extern "C" {
    #[allow(improper_ctypes)]
    static mut ein_main: Main;
}

extern "fastcall" fn input_entry(_: &Payload) -> f64 {
    42.0
}

#[no_mangle]
pub extern "C" fn main() {
    let output = unsafe {
        (ein_main.entry)(
            &ein_main.payload,
            &Input {
                entry: input_entry,
                payload: Payload,
            },
        )
    };

    println!("{}", (output.entry)(&output.payload));

    std::process::exit(0)
}
