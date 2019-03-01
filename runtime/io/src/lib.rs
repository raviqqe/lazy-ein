extern crate gc;

#[macro_use]
mod core;

use core::algebraic;
use core::{Closure, Number, UnsizedPayload};

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    #[allow(improper_ctypes)]
    static mut ein_main: closure!(&'static mut Number, &mut Number);
}

extern "fastcall" fn input_entry(_: &mut UnsizedPayload) -> algebraic::Number {
    42.0.into()
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe {
        gc::Allocator::initialize();
        gc::Allocator::enable_gc();
    }

    let output: f64 =
        eval!(unsafe { eval!(ein_main, &mut Number::new(input_entry, UnsizedPayload)) }).into();

    println!("{}", output);

    std::process::exit(0)
}
