cmake_minimum_required(VERSION 2.5)
project(judger C)

#set(CMAKE_VERBOSE_MAKEFILE ON)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/output)
set(LIBRARY_OUTPUT_PATH ${CMAKE_CURRENT_SOURCE_DIR}/output)

set(CMAKE_C_FLAGS "-g -Wall -Werror -O3 -std=c99 -pie -fPIC")

# make judger lib
file(GLOB SOURCE "src/*.c" "src/rules/*.c")
add_library(judger_static STATIC ${SOURCE})
SET_TARGET_PROPERTIES (judger_static PROPERTIES OUTPUT_NAME "judger")
target_link_libraries(judger_static pthread seccomp)

add_executable(libjudger.so ${SOURCE})
target_link_libraries(libjudger.so pthread seccomp)
