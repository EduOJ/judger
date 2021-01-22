cmake_minimum_required(VERSION 2.5)
project(judger C)

#set(CMAKE_VERBOSE_MAKEFILE ON)
set(LIBRARY_OUTPUT_PATH ${CMAKE_CURRENT_SOURCE_DIR}/output)

set(CMAKE_C_FLAGS "-g -Wall -Werror -O3 -std=c99 -pie -fPIC")

# make judger lib
file(GLOB SOURCE "../sandbox/src/*.c" "../sandbox/src/rules/*.c")

add_library(judger ${SOURCE})
