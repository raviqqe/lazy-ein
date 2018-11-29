#[no_mangle]
pub extern "C" fn io_main(jsonxx_main: extern "fastcall" fn()) {
    jsonxx_main();
}
