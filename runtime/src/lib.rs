#[no_mangle]
pub extern "C" fn runtime_main(jsonxx_main: extern "C" fn()) {
    jsonxx_main();
}
