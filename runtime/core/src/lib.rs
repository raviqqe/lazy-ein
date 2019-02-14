extern crate atty;
extern crate gc;
extern crate termcolor;

use std::io::Write;
use termcolor::{Color, ColorChoice, ColorSpec, StandardStream, WriteColor};

#[no_mangle]
pub extern "C" fn core_alloc(size: usize) -> *mut u8 {
    gc::Allocator::alloc(size)
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
