void io_main(void *);
struct {
} ein_main;

int main() {
  io_main(&ein_main);

  return -1;
}
