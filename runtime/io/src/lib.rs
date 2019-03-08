extern crate gc;

#[macro_use]
mod core;

use crate::core::algebraic;
use crate::core::{List, Number};

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static mut ein_main: closure!(&'static mut List<Number>, &mut Number);
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe {
        gc::Allocator::initialize();
        gc::Allocator::enable_gc();
    }

    let mut output = eval!(unsafe { eval!(ein_main, &mut 42.0.into()) });

    while let algebraic::List::Cons(elem, list) = *output {
        let n: f64 = (*eval!(unsafe { &mut *elem })).into();

        println!("{}", n);

        output = eval!(unsafe { &mut *list });
    }

    std::process::exit(0)
}
