use std::sync::atomic::{AtomicUsize, Ordering};

#[repr(C)]
pub struct MainFunction {
    entry: extern "fastcall" fn(&Environment, &mut Number) -> &'static mut List<Number>,
    payload: Environment,
}

impl MainFunction {
    pub fn call(&self, num: &mut Number) -> &'static mut List<Number> {
        (self.entry)(&self.payload, num)
    }
}

#[repr(C)]
pub struct Environment {
    _private: [u8; 0],
}

extern "fastcall" fn entry<T>(payload: &mut T) -> &mut T {
    payload
}

pub type Entry<T> = extern "fastcall" fn(payload: &mut T) -> &mut T;

#[repr(C)]
pub struct Thunk<T> {
    entry: extern "fastcall" fn(&mut T) -> &mut T,
    payload: T,
}

impl<T> Thunk<T> {
    pub fn new(payload: T) -> Self {
        Thunk { entry, payload }
    }

    pub fn new_with_entry(entry: extern "fastcall" fn(&mut T) -> &mut T, payload: T) -> Self {
        Thunk { entry, payload }
    }

    pub fn force(&mut self) -> &mut T {
        (unsafe {
            std::mem::transmute::<usize, Entry<T>>(
                std::mem::transmute::<&Entry<T>, &AtomicUsize>(&self.entry).load(Ordering::SeqCst),
            )
        })(&mut self.payload)
    }
}

impl From<f64> for Number {
    fn from(n: f64) -> Number {
        Thunk::new(n.into())
    }
}

pub type List<T> = Thunk<algebraic::List<T>>;
pub type Number = Thunk<algebraic::Number>;

pub mod algebraic {
    #[repr(C)]
    pub enum List<T> {
        #[allow(dead_code)]
        Cons(*mut T, *mut super::List<T>),
        #[allow(dead_code)]
        Nil,
    }

    #[derive(Clone, Copy)]
    #[repr(C)]
    pub struct Number(f64);

    impl From<f64> for Number {
        fn from(n: f64) -> Number {
            Number(n)
        }
    }

    impl From<Number> for f64 {
        fn from(n: Number) -> f64 {
            n.0
        }
    }
}
