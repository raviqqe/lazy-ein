extern crate coro;
extern crate crossbeam;
extern crate gc;
extern crate num_cpus;

#[macro_use]
mod core;
mod effect_ref;
mod runner;

use crate::core::MainFunction;
use crossbeam::scope;
use runner::Runner;

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static ein_main: MainFunction;
}

#[no_mangle]
pub extern "C" fn main<'a, 'b>() {
    unsafe { gc::Allocator::initialize() }

    let runner = Runner::new();

    scope(|scope| {
        runner.spawn_main_thread(scope, unsafe { &ein_main });
        runner.spawn_worker_threads(scope);
    })
    .unwrap();
}
