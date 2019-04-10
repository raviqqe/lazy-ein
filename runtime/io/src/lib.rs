#![feature(async_await, await_macro, futures_api)]

extern crate chashmap;
extern crate coro;
extern crate crossbeam;
extern crate gc;
#[macro_use]
extern crate lazy_static;
extern crate num_cpus;
#[macro_use]
extern crate tokio;
extern crate tokio_async_await;

#[macro_use]
mod core;
mod effect_ref;
mod io;
mod runner;

use crate::core::MainFunction;
use runner::RUNNER;

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static ein_main: MainFunction;
}

#[no_mangle]
pub extern "C" fn main<'a, 'b>() {
    unsafe { gc::Allocator::initialize() }
    unsafe { gc::Allocator::start_gc() }

    RUNNER.run(unsafe { &ein_main });
}
