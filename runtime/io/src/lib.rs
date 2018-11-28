#[no_mangle]
pub extern "C" fn runtime_main(jsonxx_main: extern "fastcall" fn()) {
    jsonxx_main();
}
