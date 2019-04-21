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
use std::io::Write;
use termcolor::{Color, ColorChoice, ColorSpec, StandardStream, WriteColor};

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

#[no_mangle]
pub extern "C" fn core_alloc(size: usize) -> *mut u8 {
    unsafe { std::alloc::alloc(std::alloc::Layout::from_size_align(size, 8).unwrap()) }
}

#[no_mangle]
pub extern "C" fn core_black_hole() {
    coro::suspend();
}

#[no_mangle]
pub extern "C" fn core_panic() {
    let mut stderr = StandardStream::stderr(ColorChoice::Auto);

    if atty::is(atty::Stream::Stderr) {
        stderr
            .set_color(ColorSpec::new().set_fg(Some(Color::Red)))
            .unwrap()
    }

    writeln!(&mut stderr, "Match error!").unwrap();

    std::process::exit(1)
}
