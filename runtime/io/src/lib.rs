extern crate chashmap;
extern crate coro;
extern crate crossbeam;
extern crate gc;
#[macro_use]
extern crate lazy_static;
extern crate num_cpus;

#[macro_use]
mod core;
mod effect_ref;
mod runner;

use crate::core::MainFunction;
use crossbeam::scope;
use runner::RUNNER;

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static ein_main: MainFunction;
}

#[no_mangle]
pub extern "C" fn main<'a, 'b>() {
    unsafe { gc::Allocator::initialize() }

    scope(|scope| RUNNER.spawn_threads(scope, unsafe { &ein_main })).unwrap();
}
