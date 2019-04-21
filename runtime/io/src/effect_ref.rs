use crate::core::Number;

pub struct EffectRef {
    pointer: *mut Number,
}

impl EffectRef {
    pub fn new(pointer: *mut Number) -> Self {
        EffectRef { pointer }
    }

    pub fn pointer(&self) -> *mut Number {
        self.pointer
    }
}

unsafe impl Send for EffectRef {}
