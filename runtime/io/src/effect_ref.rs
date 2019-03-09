use crate::core::Number;

pub struct EffectRef {
    pointer: *mut Number,
}

impl EffectRef {
    pub fn new(thunk: &mut Number) -> Self {
        EffectRef {
            pointer: thunk as *mut Number,
        }
    }

    pub fn pointer(&self) -> *mut Number {
        self.pointer
    }
}

unsafe impl Sync for EffectRef {}
unsafe impl Send for EffectRef {}
