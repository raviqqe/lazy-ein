macro_rules! eval {
    ($thunk:expr) => { ($thunk.entry)(&mut $thunk.payload) };
    ($thunk:expr, $($arg:expr),+) => {
        ($thunk.entry)(&mut $thunk.payload, $($arg),+)
    };
}

macro_rules! closure {
    ($result:ty) => {
        ::core::Closure<extern "fastcall" fn(&mut $result) -> $result, $result>
    };
    ($result:ty, $($arg:ty),+) => {
        ::core::Closure<extern "fastcall" fn(&mut ::core::Environment, $($arg),+) -> $result, ::core::Environment>
    };
}

#[repr(C)]
pub struct Closure<E, P> {
    pub entry: E,
    pub payload: P,
}

impl From<f64> for Number {
    fn from(n: f64) -> Number {
        Closure {
            entry: number_entry,
            payload: n.into(),
        }
    }
}

#[repr(C)]
pub struct Environment(i8); // avoid zero-sized type for compatibility with C

pub type Number = closure!(algebraic::Number);

extern "fastcall" fn number_entry(number: &mut algebraic::Number) -> algebraic::Number {
    *number
}

pub mod algebraic {
    #[derive(Clone, Copy)]
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
