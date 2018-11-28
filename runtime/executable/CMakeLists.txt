cmake_minimum_required(VERSION 3.10)

add_library(excutable STATIC src/main.c)
target_compile_options(excutable PRIVATE -Wall -Wconversion -Wextra)
