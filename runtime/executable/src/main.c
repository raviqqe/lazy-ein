void runtime_main(void *);
struct {
} ein_main;

int main() {
  runtime_main(&ein_main);

  return -1;
}
