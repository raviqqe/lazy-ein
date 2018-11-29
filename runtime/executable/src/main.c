void io_main(void (*jsonxx_main)());
void jsonxx_main();

int main() {
  io_main(jsonxx_main);

  return -1;
}
