cmake_minimum_required(VERSION 3.10)

add_library(executable STATIC src/main.c)
target_compile_options(executable PRIVATE -Wall -Wconversion -Wextra)
