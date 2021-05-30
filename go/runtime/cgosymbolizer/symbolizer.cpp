// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
#include "symbolizer.h"

#include <stdint.h>
#include <string.h>
#include <sys/types.h>

#include <boost/stacktrace/frame.hpp>
#include <boost/stacktrace/stacktrace.hpp>

#include "traceback.h"

static int append_pc_info_to_symbolizer_list(cgoSymbolizerArg* arg);
static int append_entry_to_symbolizer_list(cgoSymbolizerArg* arg);

// For the details of how this is called see runtime.SetCgoTraceback.
void cgoSymbolizer(cgoSymbolizerArg* arg) {
  cgoSymbolizerMore* more = arg->data;
  if (more != NULL) {
    arg->file = more->file;
    arg->lineno = more->lineno;
    arg->func = more->func;
    // set non-zero if more info for this PC
    arg->more = more->more != NULL;
    arg->data = more->more;

    // If returning the last file/line, we can set the
    // entry point field.
    if (!arg->more) {  // no more info
      append_entry_to_symbolizer_list(arg);
    }

    return;
  }
  arg->file = NULL;
  arg->lineno = 0;
  arg->func = NULL;
  arg->more = 0;
  if (arg->pc == 0) {
    return;
  }
  append_pc_info_to_symbolizer_list(arg);

  // If returning only one file/line, we can set the entry point field.
  if (!arg->more) {
    append_entry_to_symbolizer_list(arg);
  }
}

void prepare_syminfo(const boost::stacktrace::detail::native_frame_ptr_t addr,
                     std::string& file, std::size_t& line, std::string& func) {
  auto frame = boost::stacktrace::frame(addr);
  file = frame.source_file();
  line = frame.source_line();
  func = frame.name();

  if (!func.empty()) {
    func = boost::core::demangle(func.c_str());
  } else {
    func = boost::stacktrace::detail::to_hex_array(addr).data();
  }

  if (file.empty() || file.find_first_of("?") == 0) {
    boost::stacktrace::detail::location_from_symbol loc(addr);
    if (!loc.empty()) {
      file = loc.name();
    }
  }
}

static int append_pc_info_to_symbolizer_list(cgoSymbolizerArg* arg) {
  std::string file;
  std::size_t line = 0;
  std::string func;
  prepare_syminfo(boost::stacktrace::frame::native_frame_ptr_t(arg->pc), file,
                  line, func);
  // init head with current stack
  if (arg->file == NULL) {
    arg->file = strdup(file.c_str());
    arg->lineno = line;
    arg->func = strdup(func.c_str());
    return 0;
  }

  cgoSymbolizerMore* more = (cgoSymbolizerMore*)malloc(sizeof(*more));
  if (more == NULL) {
    return 1;
  }
  // append current stack to the tail
  more->more = NULL;
  more->file = strdup(file.c_str());
  more->lineno = line;
  more->func = strdup(func.c_str());
  cgoSymbolizerMore** pp = NULL;
  for (pp = &arg->data; *pp != NULL; pp = &(*pp)->more) {
  }
  *pp = more;
  arg->more = 1;
  return 0;
}

static int append_entry_to_symbolizer_list(cgoSymbolizerArg* arg) {
  auto frame = boost::stacktrace::frame(
      boost::stacktrace::frame::native_frame_ptr_t(arg->pc));
  arg->entry = (uintptr_t)strdup(frame.name().c_str());
  return 0;
}