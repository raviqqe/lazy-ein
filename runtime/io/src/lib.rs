extern crate gc;

#[global_allocator]
static GLOBAL: gc::Allocator = gc::Allocator;

#[repr(C)]
pub struct Closure<E, P> {
    pub entry: E,
    pub payload: P,
}

pub type Input = Closure<extern "fastcall" fn(&()) -> f64, ()>;
pub type Output = Closure<extern "fastcall" fn(&()) -> f64, ()>;
pub type Main = Closure<extern "fastcall" fn(&(), &Input) -> Box<Output>, ()>;

extern "fastcall" fn input_entry(_: &()) -> f64 {
    42.0
}

#[no_mangle]
pub extern "C" fn io_main(main: &'static Main) {
    let output = (main.entry)(
        &main.payload,
        &Input {
            entry: input_entry,
            payload: (),
        },
    );

    println!("{}", (output.entry)(&output.payload));

    std::process::exit(0)
}
