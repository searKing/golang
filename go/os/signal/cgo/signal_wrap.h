#ifndef SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_WRAP_H_
#define SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_WRAP_H_

#include <functional>
#include <string>
namespace searking {
int SignalAction(bool enable, int signum);
} // namespace searking
#endif // SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_WRAP_H_
