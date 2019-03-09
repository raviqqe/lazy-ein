extern crate crossbeam;
extern crate gc;
extern crate num_cpus;

#[macro_use]
mod core;
mod effect_ref;

use crate::core::algebraic;
use crate::core::{List, Number};
use crossbeam::deque::{Injector, Steal};
use crossbeam::scope;
use effect_ref::EffectRef;
use std::sync::atomic::{AtomicBool, AtomicUsize, Ordering};
use std::sync::Arc;

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static mut ein_main: closure!(&'static mut List<Number>, &mut Number);
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe { gc::Allocator::initialize() }

    let thunk_injector = Arc::new(Injector::new());
    let num_rest_effects = AtomicUsize::new(0);
    let gc_enabled = AtomicBool::new(false);
    let gc_ready = AtomicUsize::new(0);

    scope(|scope| {
        // Main thread

        {
            let thunk_injector = thunk_injector.clone();
            let num_rest_effects = &num_rest_effects;
            let gc_enabled = &gc_enabled;
            let gc_ready = &gc_ready;

            scope.spawn(move |_| {
                gc::Allocator::register_current_thread().unwrap();
                gc_ready.fetch_add(1, Ordering::SeqCst);
                while !gc_enabled.load(Ordering::SeqCst) {}

                let mut output = eval!(unsafe { eval!(ein_main, &mut 42.0.into()) });

                while let algebraic::List::Cons(elem, list) = *output {
                    num_rest_effects.fetch_add(1, Ordering::SeqCst);
                    thunk_injector.push(EffectRef::new(unsafe { &mut *elem }));
                    output = eval!(unsafe { &mut *list });
                }

                while num_rest_effects.load(Ordering::SeqCst) != 0 {
                    delay()
                }

                std::process::exit(0)
            });
        }

        // Worker threads

        let num_workers = num_cpus::get();

        for _ in 0..num_workers {
            let thunk_injector = thunk_injector.clone();
            let num_rest_effects = &num_rest_effects;
            let gc_enabled = &gc_enabled;
            let gc_ready = &gc_ready;

            scope.spawn(move |_| {
                gc::Allocator::register_current_thread().unwrap();
                gc_ready.fetch_add(1, Ordering::SeqCst);
                while !gc_enabled.load(Ordering::SeqCst) {}

                loop {
                    match thunk_injector.steal() {
                        Steal::Success(thunk) => {
                            let num: f64 = (*eval!(unsafe { &mut *thunk.pointer() })).into();
                            println!("{}", num);
                            num_rest_effects.fetch_sub(1, Ordering::SeqCst);
                        }
                        Steal::Empty => delay(),
                        Steal::Retry => {}
                    }
                }
            });
        }

        // GC starter thread

        let gc_enabled = &gc_enabled;
        let gc_ready = &gc_ready;

        scope.spawn(move |_| {
            while gc_ready.load(Ordering::SeqCst) < num_workers + 1 {}
            unsafe { gc::Allocator::enable_gc() }
            gc_enabled.store(true, Ordering::SeqCst);
        });
    })
    .unwrap();
}

fn delay() {
    std::thread::sleep(std::time::Duration::from_millis(1))
}
