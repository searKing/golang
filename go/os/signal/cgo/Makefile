######################################
#
#
######################################
  
#target
#共享库文件名，lib*.a
TARGET  := libsignal_cgo.a
  
#compile and lib parameter
#编译参数
CC      := gcc
AR      = ar
RANLIB  = ranlib
LIBS    := -ldl
LDFLAGS :=
DEFINES := 
INCLUDE := -I.
CFLAGS  := -g #-Wall -O3
CFLAGS  := $(CFLAGS) $(DEFINES) $(INCLUDE)
CXXFLAGS:= $(CFLAGS)  -std=c++11
  
#source file
#源文件，自动找所有.c和.cpp文件，并将目标定义为同名.o文件
SOURCE  := $(wildcard *.c) $(wildcard *.cpp)
OBJS    := $(patsubst %.c,%.o,$(patsubst %.cpp,%.o,$(SOURCE)))
  
.PHONY : everything objs clean veryclean rebuild
  
everything : $(TARGET)
  
all : $(TARGET)
  
objs : $(OBJS)
  
rebuild: veryclean everything
                
clean :
	rm -fr *.o
    
veryclean : clean
	rm -fr $(TARGET)
  
$(TARGET) : $(OBJS)
	$(AR) cru $(TARGET) $(OBJS)
	$(RANLIB) $(TARGET)