extern crate atty;
extern crate gc;
extern crate termcolor;

use std::cell::RefCell;
use std::io::Write;
use std::vec::Vec;
use termcolor::{Color, ColorChoice, ColorSpec, StandardStream, WriteColor};

thread_local! {
    // TODO: Vec::with_capacity()
    pub static BLACK_HOLE: RefCell<Vec<*mut u8>> = RefCell::new(Vec::new())
}

#[no_mangle]
pub extern "C" fn core_alloc(size: usize) -> *mut u8 {
    gc::Allocator::alloc(size)
}

#[no_mangle]
pub extern "C" fn core_black_hole(thunk: *mut u8) {
    BLACK_HOLE.with(|black_hole| black_hole.borrow_mut().push(thunk));
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
