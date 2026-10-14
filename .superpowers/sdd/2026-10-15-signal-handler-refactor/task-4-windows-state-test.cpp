#include <cassert>
#include <csignal>
#include <mutex>

#include "base_signal_handler.hpp"

// Expose private implementation details only in this standalone test helper.
#define private public
#include "signal_handler_windows.hpp"
#undef private

using Handler = searking::SignalHandlerSignalHandler;

namespace {

constexpr int kSignal = SIGTERM;
Handler returned_handler;
Handler expected_install;
const Handler *observed_record;
int raise_count;

void HandlerA(int) {}
void HandlerB(int) {}

}  // namespace

Handler TestSignal(int signum, Handler handler) {
  assert(signum == kSignal);
  assert(handler == expected_install);
  observed_record =
      searking::SignalHandler::GetInstance().go_registered_handlers_.Load(
          signum);
  return returned_handler;
}

int TestRaise(int signum) {
  assert(signum == kSignal);
  ++raise_count;
  return 0;
}

#define signal TestSignal
#define raise TestRaise
#include "signal_handler_windows.cpp"
#undef raise
#undef signal

int main() {
  searking::SignalHandler &instance = searking::SignalHandler::GetInstance();
  instance.go_registered_handlers_.Store(kSignal, nullptr);

  expected_install = &searking::SignalHandler::HandleSignal;
  returned_handler = &HandlerA;
  assert(searking::SignalHandler::SetSig(kSignal) == 0);
  assert(observed_record == nullptr);
  const Handler *saved = instance.go_registered_handlers_.Load(kSignal);
  assert(saved != nullptr && *saved == &HandlerA);

  returned_handler = &searking::SignalHandler::HandleSignal;
  assert(searking::SignalHandler::SetSig(kSignal) == 0);
  assert(observed_record == nullptr);
  assert(instance.go_registered_handlers_.Load(kSignal) == saved);

  returned_handler = SIG_ERR;
  assert(searking::SignalHandler::SetSig(kSignal) == -1);
  assert(observed_record == nullptr);
  assert(instance.go_registered_handlers_.Load(kSignal) == saved);

  returned_handler = &HandlerB;
  assert(searking::SignalHandler::SetSig(kSignal) == 0);
  assert(observed_record == nullptr);
  const Handler *replacement = instance.go_registered_handlers_.Load(kSignal);
  assert(replacement != nullptr && replacement != saved);
  assert(*replacement == &HandlerB);

  instance.go_registered_handlers_.Store(kSignal, nullptr);
  expected_install = SIG_DFL;
  returned_handler = SIG_DFL;
  raise_count = 0;
  instance.InvokeGoSignalHandler(kSignal);
  assert(observed_record == nullptr);
  assert(raise_count == 1);
}
