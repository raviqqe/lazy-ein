extern crate libc;

use libc::{c_void, size_t};
use std::alloc::{GlobalAlloc, Layout};

#[link(name = "gc")]
extern "C" {
    pub fn GC_malloc(size: size_t) -> *mut c_void;
    pub fn GC_realloc(ptr: *mut c_void, size: size_t) -> *mut c_void;
    pub fn GC_free(ptr: *mut c_void);
}

pub struct Allocator;

unsafe impl GlobalAlloc for Allocator {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        GC_malloc(layout.size()) as *mut u8
    }

    unsafe fn dealloc(&self, ptr: *mut u8, _layout: Layout) {
        GC_free(ptr as *mut c_void)
    }

    unsafe fn realloc(&self, ptr: *mut u8, _layout: Layout, size: usize) -> *mut u8 {
        GC_realloc(ptr as *mut c_void, size as size_t) as *mut u8
    }
}
