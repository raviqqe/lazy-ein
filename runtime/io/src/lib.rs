extern crate gc;

#[macro_use]
mod core;

use crate::core::Number;

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static mut ein_main: closure!(&'static mut Number, &mut Number);
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe {
        gc::Allocator::initialize();
        gc::Allocator::enable_gc();
    }

    let output: f64 = (*eval!(unsafe { eval!(ein_main, &mut 42.0.into()) })).into();

    println!("{}", output);

    std::process::exit(0)
}
