macro_rules! eval {
    ($thunk:expr) => { ($thunk.entry)(&mut $thunk.payload) };
    ($thunk:expr, $($arg:expr),+) => {
        ($thunk.entry)(&mut $thunk.payload, $($arg),+)
    };
}

macro_rules! closure {
    ($result:ty) => {
        Closure<extern "fastcall" fn(&mut UnsizedPayload) -> $result, UnsizedPayload>
    };
    ($result:ty, $($arg:ty),+) => {
        Closure<extern "fastcall" fn(&mut UnsizedPayload, $($arg),+) -> $result, UnsizedPayload>
    };
}

#[repr(C)]
pub struct Closure<E, P> {
    pub entry: E,
    pub payload: P,
}

impl<E, P> Closure<E, P> {
    pub fn new(entry: E, payload: P) -> Closure<E, P> {
        Closure { entry, payload }
    }
}

#[repr(C)]
pub struct UnsizedPayload;

pub type Number = closure!(algebraic::Number);

pub mod algebraic {
    #[repr(C)]
    pub struct Number(f64);

    impl From<f64> for Number {
        fn from(n: f64) -> Number {
            Number(n)
        }
    }

    impl Into<f64> for Number {
        fn into(self) -> f64 {
            self.0
        }
    }
}
