#![feature(async_await, await_macro, futures_api)]

extern crate bdwgc_alloc;
extern crate chashmap;
extern crate coro;
extern crate crossbeam;
#[macro_use]
extern crate lazy_static;
extern crate num_cpus;
#[macro_use]
extern crate tokio;
extern crate tokio_async_await;

mod core;
mod io;
mod runner;
mod util;

use crate::core::MainFunction;
use runner::RUNNER;

#[global_allocator]
static GLOBAL_ALLOCATOR: bdwgc_alloc::Allocator = bdwgc_alloc::Allocator;

extern "C" {
    static ein_main: MainFunction;
}

#[no_mangle]
pub extern "C" fn main<'a, 'b>() {
    unsafe { bdwgc_alloc::Allocator::initialize() }

    RUNNER.run(unsafe { &ein_main });
}
