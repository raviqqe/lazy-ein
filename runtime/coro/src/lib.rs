extern crate coroutine;

mod error;

use coroutine::asymmetric;
use std::cell::RefCell;
use std::error::Error;

thread_local!(
    static CURRENT_COROUTINE: RefCell<*mut asymmetric::Coroutine> =
        RefCell::new(unsafe { std::mem::uninitialized() });
);

fn set_current(coroutine: &mut asymmetric::Coroutine) {
    CURRENT_COROUTINE.with(|current| current.replace(coroutine as *mut asymmetric::Coroutine));
}

fn get_current<'a>() -> &'a mut asymmetric::Coroutine {
    CURRENT_COROUTINE.with(|current| unsafe { &mut **current.borrow() })
}

pub fn suspend() {
    let current = get_current();
    current.yield_with(0);
    set_current(current)
}

pub fn park() {
    let current = get_current();
    current.park_with(0);
    set_current(current)
}

pub fn spawn<F: FnOnce() + 'static>(func: F) -> Handle {
    Handle(asymmetric::Coroutine::spawn_opts(
        move |coroutine, _| {
            set_current(coroutine);
            func();
            0
        },
        coroutine::Options {
            stack_size: 2 * 1024,
            name: None,
        },
    ))
}

pub struct Handle(asymmetric::Handle);

impl Handle {
    pub fn resume(&mut self) -> Result<State, error::Error> {
        self.0
            .resume(0)
            .map_err(|err| error::Error::new(err.description().into()))?;

        match self.0.state() {
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
