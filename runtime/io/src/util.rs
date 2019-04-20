use crate::core::{Entry, Thunk};
use crate::runner::RUNNER;
use coro;
use std::sync::atomic::{AtomicUsize, Ordering};

pub fn lock<T>(
    coroutine_id: usize,
    payload: &mut T,
    initial_entry: Entry<T>,
) -> Result<Unlocker<T>, &mut T> {
    if unsafe { &mut *payload_to_entry_ptr(payload) }
        .compare_exchange(
            entry_to_usize(initial_entry),
            entry_to_usize::<T>(black_hole_entry),
            Ordering::SeqCst,
            Ordering::SeqCst,
        )
        .is_err()
    {
        return Err(unsafe { &mut *payload_to_thunk_ptr(payload) }.force());
    }

    Ok(Unlocker::new(coroutine_id, payload))
}

#[derive(Clone, Copy)]
pub struct Unlocker<T> {
    coroutine_id: usize,
    payload: *mut T,
}

impl<T> Unlocker<T> {
    pub fn new(coroutine_id: usize, payload: *mut T) -> Self {
        Unlocker {
            coroutine_id,
            payload,
        }
    }

    pub unsafe fn unlock(&self, result: T) {
        *self.payload = result;
        (&mut *payload_to_entry_ptr(self.payload))
            .store(entry_to_usize::<T>(normal_entry), Ordering::SeqCst);
        RUNNER.unpark(self.coroutine_id);
    }

    pub unsafe fn wait_unlock(&self) {
        coro::park();
    }
}

unsafe impl<T> Send for Unlocker<T> {}
unsafe impl<T> Sync for Unlocker<T> {}

fn payload_to_entry_ptr<T>(payload: *mut T) -> *mut AtomicUsize {
    unsafe { std::mem::transmute::<*mut Thunk<T>, *mut AtomicUsize>(payload_to_thunk_ptr(payload)) }
}

fn payload_to_thunk_ptr<T>(payload: *mut T) -> *mut Thunk<T> {
    unsafe {
        std::mem::transmute::<usize, *mut Thunk<T>>(
            std::mem::transmute::<*mut T, usize>(payload) - std::mem::size_of::<usize>(),
        )
    }
}

fn entry_to_usize<T>(entry: Entry<T>) -> usize {
    unsafe { std::mem::transmute::<Entry<T>, usize>(entry) }
}

extern "fastcall" fn black_hole_entry<T>(payload: &mut T) -> &mut T {
    coro::suspend();
    unsafe { &mut *payload_to_thunk_ptr(payload) }.force()
}

extern "fastcall" fn normal_entry<T>(payload: &mut T) -> &mut T {
    payload
}
