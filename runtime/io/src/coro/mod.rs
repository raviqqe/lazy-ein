extern crate coroutine;

mod error;

use coroutine::asymmetric;
use std::cell::RefCell;
use std::error::Error;
use std::sync::atomic::{AtomicUsize, Ordering};

static CUNRRENT_COROUTINE_ID: AtomicUsize = AtomicUsize::new(0);

thread_local!(
    static CURRENT_COROUTINE: RefCell<Coroutine> =
        RefCell::new(unsafe { std::mem::uninitialized() });
);

fn set_current(coroutine: Coroutine) {
    CURRENT_COROUTINE.with(|current| current.replace(coroutine));
}

fn get_current() -> Coroutine {
    CURRENT_COROUTINE.with(|current| *current.borrow())
}

fn new_id() -> usize {
    CUNRRENT_COROUTINE_ID.fetch_add(1, Ordering::SeqCst)
}

pub fn current_id() -> usize {
    get_current().id()
}

pub fn suspend() {
    let current = get_current();
    current.ptr().yield_with(0);
    set_current(current)
}

pub fn park() {
    let current = get_current();
    current.ptr().park_with(0);
    set_current(current)
}

pub fn spawn<F: FnOnce() + 'static>(func: F) -> Handle {
    let id = new_id();

    Handle::new(
        id,
        asymmetric::Coroutine::spawn_opts(
            move |ptr, _| {
                set_current(Coroutine::new(id, ptr));
                func();
                0
            },
            coroutine::Options {
                stack_size: 2 * 1024 * 1024,
                name: None,
            },
        ),
    )
}

#[derive(Clone, Copy)]
struct Coroutine {
    id: usize,
    ptr: *mut asymmetric::Coroutine,
}

impl Coroutine {
    fn new(id: usize, ptr: &mut asymmetric::Coroutine) -> Self {
        Coroutine { id, ptr }
    }

    fn id(&self) -> usize {
        self.id
    }

    fn ptr<'a>(&self) -> &'a mut asymmetric::Coroutine {
        unsafe { &mut *self.ptr }
    }
}

pub struct Handle {
    coroutine_id: usize,
    handle: asymmetric::Handle,
}

impl Handle {
    pub fn new(coroutine_id: usize, handle: asymmetric::Handle) -> Self {
        Handle {
            coroutine_id,
            handle,
        }
    }

    pub fn coroutine_id(&self) -> usize {
        self.coroutine_id
    }

    pub fn resume(&mut self) -> Result<State, error::Error> {
        self.handle
            .resume(0)
            .map_err(|err| error::Error::new(err.description().into()))?;

        match self.handle.state() {
            asymmetric::State::Finished => Ok(State::Finished),
            asymmetric::State::Parked => Ok(State::Parked),
            asymmetric::State::Suspended => Ok(State::Suspended),
            asymmetric::State::Running | asymmetric::State::Panicked => {
                Err(error::Error::new("invalid coroutine state".into()))
            }
        }
    }
}

unsafe impl Send for Handle {}

#[derive(Debug, Eq, PartialEq)]
pub enum State {
    Finished,
    Parked,
    Suspended,
}

mod test {
    #[test]
    fn spawn() {
        let mut handle = super::spawn(|| {});

        assert_eq!(handle.resume().unwrap(), super::State::Finished);
    }

    #[test]
    fn suspend() {
        let mut handle = super::spawn(|| super::suspend());

        assert_eq!(handle.resume().unwrap(), super::State::Suspended);
        assert_eq!(handle.resume().unwrap(), super::State::Finished);
    }

    #[test]
    fn park() {
        let mut handle = super::spawn(|| super::park());

        assert_eq!(handle.resume().unwrap(), super::State::Parked);
        assert_eq!(handle.resume().unwrap(), super::State::Finished);
    }
}
