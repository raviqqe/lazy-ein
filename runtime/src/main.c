void runtime_main(void (*jsonxx_main)());

void jsonxx_main() {}

int main() {
  runtime_main(jsonxx_main);

  return -1;
}
