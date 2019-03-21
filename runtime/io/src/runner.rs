use crate::core::algebraic;
use crate::core::MainFunction;
use crate::effect_ref::EffectRef;
use chashmap::CHashMap;
use coro;
use crossbeam::deque::{Injector, Steal};
use crossbeam::thread::Scope;
use std::sync::atomic::{AtomicBool, AtomicUsize, Ordering};
use std::sync::Mutex;

lazy_static! {
    pub static ref RUNNER: Runner = Runner::new();
}

pub struct Runner {
    effect_injector: Injector<EffectRef>,
    ready_coroutine_injector: Injector<coro::Handle>,
    suspended_coroutine_injector: Injector<coro::Handle>,
    parked_coroutine_map: CHashMap<usize, Mutex<coro::Handle>>,
    gc_started: AtomicBool,
    num_rest_effects: AtomicUsize,
    num_spawned_workers: AtomicUsize,
    num_workers: usize,
}

impl Runner {
    pub fn new() -> Runner {
        Runner {
            effect_injector: Injector::new(),
            ready_coroutine_injector: Injector::new(),
            suspended_coroutine_injector: Injector::new(),
            parked_coroutine_map: CHashMap::new(),
            gc_started: AtomicBool::new(false),
            num_rest_effects: AtomicUsize::new(0),
            num_spawned_workers: AtomicUsize::new(0),
            num_workers: num_cpus::get(),
        }
    }

    pub fn spawn_threads<'a, 'b>(&'a self, scope: &'b Scope<'a>, main: &'static MainFunction) {
        self.spawn_main_thread(scope, main);
        self.spawn_worker_threads(scope);
    }

    fn spawn_main_thread<'a, 'b>(&'a self, scope: &'b Scope<'a>, main: &'static MainFunction) {
        scope.spawn(move |_| {
            gc::Allocator::register_current_thread().unwrap();
            while self.num_spawned_workers.load(Ordering::SeqCst) < self.num_workers {}
            unsafe { gc::Allocator::start_gc() }
            self.gc_started.store(true, Ordering::SeqCst);

            let mut output = main.call(&mut 42.0.into()).force();

            while let algebraic::List::Cons(elem, list) = *output {
                self.num_rest_effects.fetch_add(1, Ordering::SeqCst);
                self.effect_injector
                    .push(EffectRef::new(unsafe { &mut *elem }));
                output = unsafe { &mut *list }.force();
            }

            while self.num_rest_effects.load(Ordering::SeqCst) != 0 {
                delay()
            }

            std::process::exit(0)
        });
    }

    fn spawn_worker_threads<'a, 'b>(&'a self, scope: &'b Scope<'a>) {
        for _ in 0..self.num_workers {
            scope.spawn(move |_| {
                gc::Allocator::register_current_thread().unwrap();
                self.num_spawned_workers.fetch_add(1, Ordering::SeqCst);
                while !self.gc_started.load(Ordering::SeqCst) {}

                loop {
                    let mut handle = match self.steal() {
                        Some(handle) => handle,
                        None => {
                            delay();
                            continue;
                        }
                    };

                    match handle.resume() {
                        Ok(coro::State::Finished) => {
                            self.num_rest_effects.fetch_sub(1, Ordering::SeqCst);
                        }
                        Ok(coro::State::Suspended) => {
                            self.suspended_coroutine_injector.push(handle)
                        }
                        Ok(coro::State::Parked) => {
                            self.parked_coroutine_map
                                .insert(handle.coroutine_id(), Mutex::new(handle));
                        }
                        Err(error) => {
                            eprintln!("{}", error);
                            std::process::exit(1);
                        }
                    }
                }
            });
        }
    }

    pub fn unpark(&self, coroutine_id: usize) {
        loop {
            match self.parked_coroutine_map.remove(&coroutine_id) {
                Some(mutex) => {
                    self.ready_coroutine_injector
                        .push(mutex.into_inner().unwrap());
                    break;
                }
                None => {}
            }
        }
    }

    fn steal(&self) -> Option<coro::Handle> {
        loop {
            match self.ready_coroutine_injector.steal() {
                Steal::Success(handle) => return Some(handle),
                Steal::Empty => break,
                Steal::Retry => {}
            }
        }

        loop {
            match self.effect_injector.steal() {
                Steal::Success(thunk) => {
                    return Some(coro::spawn(move || {
                        let num: f64 = (*unsafe { &mut *thunk.pointer() }.force()).into();
                        println!("{}", num);
                    }))
                }
                Steal::Empty => break,
                Steal::Retry => {}
            }
        }

        loop {
            match self.suspended_coroutine_injector.steal() {
                Steal::Success(handle) => return Some(handle),
                Steal::Empty => return None,
                Steal::Retry => {}
            }
        }
    }
}

fn delay() {
    std::thread::sleep(std::time::Duration::from_millis(1))
}
