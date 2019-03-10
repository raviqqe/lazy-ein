#[repr(C)]
pub struct MainFunction {
    pub entry: extern "fastcall" fn(&mut Environment, &mut Number) -> &'static mut List<Number>,
    pub payload: Environment,
}

impl MainFunction {
    pub fn call(&mut self, num: &mut Number) -> &'static mut List<Number> {
        (self.entry)(&mut self.payload, num)
    }
}

#[repr(C)]
pub struct Environment {
    _private: [u8; 0],
}

extern "fastcall" fn entry<T>(payload: &mut T) -> &mut T {
    payload
}

#[repr(C)]
pub struct Thunk<T> {
    entry: extern "fastcall" fn(&mut T) -> &mut T,
    payload: T,
}

impl<T> Thunk<T> {
    pub fn new(payload: T) -> Self {
        Thunk { entry, payload }
    }

    pub fn force(&mut self) -> &mut T {
        (self.entry)(&mut self.payload)
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
