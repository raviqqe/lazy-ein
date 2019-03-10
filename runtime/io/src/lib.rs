extern crate crossbeam;
extern crate gc;
extern crate num_cpus;

#[macro_use]
mod core;
mod effect_ref;

use crate::core::algebraic;
use crate::core::MainFunction;
use crossbeam::deque::{Injector, Steal};
use crossbeam::scope;
use effect_ref::EffectRef;
use std::sync::atomic::{AtomicBool, AtomicUsize, Ordering};

#[global_allocator]
static mut GLOBAL_ALLOCATOR: gc::Allocator = gc::Allocator;

extern "C" {
    static mut ein_main: MainFunction;
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe { gc::Allocator::initialize() }

    let thunk_injector = Injector::new();
    let num_rest_effects = AtomicUsize::new(0);
    let gc_started = AtomicBool::new(false);
    let gc_ready = AtomicUsize::new(0);

    scope(|scope| {
        // Main thread

        {
            let thunk_injector = &thunk_injector;
            let num_rest_effects = &num_rest_effects;
            let gc_started = &gc_started;
            let gc_ready = &gc_ready;

            scope.spawn(move |_| {
                gc::Allocator::register_current_thread().unwrap();
                gc_ready.fetch_add(1, Ordering::SeqCst);
                while !gc_started.load(Ordering::SeqCst) {}

                let mut output = unsafe { ein_main.call(&mut 42.0.into()) }.force();

                while let algebraic::List::Cons(elem, list) = *output {
                    num_rest_effects.fetch_add(1, Ordering::SeqCst);
                    thunk_injector.push(EffectRef::new(unsafe { &mut *elem }));
                    output = unsafe { &mut *list }.force();
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
            let thunk_injector = &thunk_injector;
            let num_rest_effects = &num_rest_effects;
            let gc_started = &gc_started;
            let gc_ready = &gc_ready;

            scope.spawn(move |_| {
                gc::Allocator::register_current_thread().unwrap();
                gc_ready.fetch_add(1, Ordering::SeqCst);
                while !gc_started.load(Ordering::SeqCst) {}

                loop {
                    match thunk_injector.steal() {
                        Steal::Success(thunk) => {
                            let num: f64 = (*unsafe { &mut *thunk.pointer() }.force()).into();
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

        let gc_started = &gc_started;
        let gc_ready = &gc_ready;

        scope.spawn(move |_| {
            while gc_ready.load(Ordering::SeqCst) < num_workers + 1 {}
            unsafe { gc::Allocator::start_gc() }
            gc_started.store(true, Ordering::SeqCst);
        });
    })
    .unwrap();
}

fn delay() {
    std::thread::sleep(std::time::Duration::from_millis(1))
}
