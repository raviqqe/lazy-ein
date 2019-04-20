extern crate atty;
extern crate bdwgc_alloc;
extern crate coro;
extern crate termcolor;

use std::alloc::GlobalAlloc;
use std::io::Write;
use termcolor::{Color, ColorChoice, ColorSpec, StandardStream, WriteColor};

#[no_mangle]
pub extern "C" fn core_alloc(size: usize) -> *mut u8 {
    unsafe { bdwgc_alloc::Allocator.alloc(std::alloc::Layout::from_size_align(size, 8).unwrap()) }
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
