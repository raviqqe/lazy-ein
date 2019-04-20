use crate::core;
use crate::core::algebraic;
use crate::util;
use coro;
use tokio;
use tokio::io;
use tokio::prelude::*;

extern "fastcall" fn entry(payload: &mut algebraic::Number) -> &mut algebraic::Number {
    let unlocker = match util::lock(coro::current_id(), payload, entry) {
        Ok(unlocker) => unlocker,
        Err(_payload) => return payload,
    };

    tokio::spawn_async(async move {
        let mut stdin = io::stdin();
        let mut buf: [u8; 1024] = [0; 1024];

        unsafe {
            unlocker.unlock(
                match await!(stdin.read_async(&mut buf)) {
                    Ok(n) => {
                        if n > 0 {
                            buf[0].into()
                        } else {
                            13.0
                        }
                    }
                    Err(_) => 13.0,
                }
                .into(),
            )
        }
    });

    unsafe { unlocker.wait_unlock() }
    payload
}

pub fn create_input() -> &'static mut core::Number {
    let ptr: &'static mut core::Number = unsafe {
        &mut *(std::alloc::alloc(std::alloc::Layout::new::<core::Number>()) as *mut core::Number)
    };
    *ptr = core::Number::new_with_entry(entry, 0.0.into());
    ptr
}
