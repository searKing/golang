#ifndef SEARKING_GOLANG_GO_OS_SIGNAL_CGO_BACKTRACE_H_
#define SEARKING_GOLANG_GO_OS_SIGNAL_CGO_BACKTRACE_H_

#include <string>
namespace searking {
std::string Backtrace(int skip = 1);
} // namespace searking
#endif // SEARKING_GOLANG_GO_OS_SIGNAL_CGO_BACKTRACE_H_
