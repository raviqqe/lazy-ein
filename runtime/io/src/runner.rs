use crate::core::{algebraic, MainFunction};
use crate::io::create_input;
use chashmap::CHashMap;
use coro;
use crossbeam::deque::{Injector, Steal};
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Mutex;

lazy_static! {
    pub static ref RUNNER: Runner = Runner::new();
}

pub struct Runner {
    effect_injector: Injector<coro::Handle>,
    ready_coroutine_injector: Injector<coro::Handle>,
    suspended_coroutine_injector: Injector<coro::Handle>,
    parked_coroutine_map: CHashMap<usize, Mutex<coro::Handle>>,
    num_rest_effects: AtomicUsize,
}

impl Runner {
    pub fn new() -> Runner {
        Runner {
            effect_injector: Injector::new(),
            ready_coroutine_injector: Injector::new(),
            suspended_coroutine_injector: Injector::new(),
            parked_coroutine_map: CHashMap::new(),
            num_rest_effects: AtomicUsize::new(0),
        }
    }

    pub fn run(&'static self, main: &'static MainFunction) {
        let mut runtime = tokio::runtime::Builder::new()
            .after_start(move || unsafe {
                // This thread registration assumes that no heap allocation is done before itself.
                bdwgc_alloc::Allocator::register_current_thread().unwrap();
            })
            .before_stop(move || unsafe {
                bdwgc_alloc::Allocator::unregister_current_thread();
            })
            .build()
            .unwrap();

        for _ in 0..num_cpus::get() {
            runtime.spawn(tokio_async_await::compat::backward::Compat::new(
                async move {
                    loop {
                        let mut handle = match self.steal() {
                            Some(handle) => handle,
                            None => {
                                await!(delay()).unwrap();
                                continue;
                            }
                        };

                        match handle.resume() {
                            Ok(coro::State::Finished) => {}
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
                },
            ));
        }

        runtime
            .block_on::<_, (), ()>(tokio_async_await::compat::backward::Compat::new(
                async move {
                    {
                        let mut output = main.call(create_input()).force();

                        while let algebraic::List::Cons(elem, list) = *output {
                            self.num_rest_effects.fetch_add(1, Ordering::SeqCst);
                            self.effect_injector.push(coro::spawn(move || {
                                let num: f64 = (*unsafe { &mut *elem }.force()).into();
                                println!("{}", num);
                                self.num_rest_effects.fetch_sub(1, Ordering::SeqCst);
                            }));
                            output = unsafe { &mut *list }.force();
                        }
                    }

                    while self.num_rest_effects.load(Ordering::SeqCst) != 0 {
                        await!(delay()).unwrap();
                    }

                    std::process::exit(0)
                },
            ))
            .unwrap();
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

            // TODO: Consider inserting delay here.
        }
    }

    fn steal(&'static self) -> Option<coro::Handle> {
        loop {
            match self.ready_coroutine_injector.steal() {
                Steal::Success(handle) => return Some(handle),
                Steal::Empty => break,
                Steal::Retry => {}
            }
        }

        loop {
            match self.effect_injector.steal() {
                Steal::Success(handle) => return Some(handle),
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

async fn delay() -> Result<(), tokio::timer::Error> {
    await!(tokio::timer::Delay::new(
        std::time::Instant::now() + std::time::Duration::from_millis(1)
    ))
}
