use crate::core;
use crate::core::algebraic;
use crate::runner::RUNNER;
use coro;
use std::ops::{Deref, DerefMut};
use tokio;
use tokio::io;
use tokio::prelude::*;

extern "fastcall" fn entry(payload: &mut algebraic::Number) -> &mut algebraic::Number {
    // TODO: Lock thunks here.

    let mut payload_ref = Ref::new(payload);
    let id = coro::current_id();

    tokio::spawn_async(async move {
        let mut stdin = io::stdin();
        let mut buf: [u8; 1024] = [0; 1024];

        *payload_ref = match await!(stdin.read_async(&mut buf)) {
            Ok(_) => 123.0,
            Err(_) => 456.0,
        }
        .into();

        RUNNER.unpark(id)
    });

    coro::park();
    return payload;
}

pub fn create_input() -> &'static mut core::Number {
    let ptr = unsafe { std::alloc::alloc(std::alloc::Layout::new::<core::Number>()) }
        as *mut core::Number;
    unsafe { *ptr = core::Number::new_with_entry(entry, 42.0.into()) };
    unsafe { &mut *ptr }
}

#[derive(Debug)]
struct Ref<T>(*mut T);

impl<T> Ref<T> {
    fn new(ptr: *mut T) -> Self {
        Ref(ptr)
    }
}

unsafe impl<T> Send for Ref<T> {}

impl<T> Deref for Ref<T> {
    type Target = T;

    fn deref(&self) -> &T {
        unsafe { &*self.0 }
    }
}

impl<T> DerefMut for Ref<T> {
    fn deref_mut(&mut self) -> &mut T {
        unsafe { &mut *self.0 }
    }
}
