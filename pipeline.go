package gst

/*
#include <stdlib.h>
#include <gst/gst.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/greenkeytech/glib"
)

type Pipeline struct {
	Bin
}

func (p *Pipeline) g() *C.GstPipeline {
	return (*C.GstPipeline)(p.GetPtr())
}

func (p *Pipeline) AsPipeline() *Pipeline {
	return p
}

func (p *Pipeline) Free() {
	C.g_object_unref(C.gpointer(p.g()))
}

func (p *Pipeline) FlushAndResetTime() {
	flushStart := C.gst_event_new_flush_start()
	flushStop := C.gst_event_new_flush_stop(0)
	startRet := C.gst_element_send_event((*C.GstElement)(p.GetPtr()), flushStart)
	fmt.Printf("startRet: %v", startRet)
	stopRet := C.gst_element_send_event((*C.GstElement)(p.GetPtr()), flushStop)
	fmt.Printf("stopRet: %v", stopRet)
}

func NewPipeline(name string) *Pipeline {
	s := (*C.gchar)(C.CString(name))
	defer C.free(unsafe.Pointer(s))
	p := new(Pipeline)
	p.SetPtr(glib.Pointer(C.gst_pipeline_new(s)))
	return p
}

func ParseLaunch(pipeline_description string) (*Pipeline, error) {
	pd := (*C.gchar)(C.CString(pipeline_description))
	defer C.free(unsafe.Pointer(pd))
	p := new(Pipeline)
	var Cerr *C.GError
	p.SetPtr(glib.Pointer(C.gst_parse_launch(pd, &Cerr)))
	if Cerr != nil {
		err := *(*glib.Error)(unsafe.Pointer(Cerr))
		C.g_error_free(Cerr)
		return nil, &err
	}
	return p, nil
}
