import ctypes as C
import platform

try:
  if platform.system() == 'Darwin':
    libuvc = C.cdll.LoadLibrary("libuvc.dylib")
  elif platform.system() == 'Linux':
    libuvc = C.cdll.LoadLibrary("libuvc.so")
  else:
    libuvc = C.cdll.LoadLibrary("libuvc")
except OSError:
  print("Error: could not find libuvc!")
  exit(1)

class uvc_context(C.Structure):
  _fields_ = [("usb_ctx", C.c_void_p),
              ("own_usb_ctx", C.c_uint8),
              ("open_devices", C.c_void_p),
              ("handler_thread", C.c_ulong),
              ("kill_handler_thread", C.c_int)]

class uvc_device(C.Structure):
  _fields_ = [("ctx", C.POINTER(uvc_context)),
              ("ref", C.c_int),
              ("usb_dev", C.c_void_p)]

class uvc_stream_ctrl(C.Structure):
  _fields_ = [("bmHint", C.c_uint16),
              ("bFormatIndex", C.c_uint8),
              ("bFrameIndex", C.c_uint8),
              ("dwFrameInterval", C.c_uint32),
              ("wKeyFrameRate", C.c_uint16),
              ("wPFrameRate", C.c_uint16),
              ("wCompQuality", C.c_uint16),
              ("wCompWindowSize", C.c_uint16),
              ("wDelay", C.c_uint16),
              ("dwMaxVideoFrameSize", C.c_uint32),
              ("dwMaxPayloadTransferSize", C.c_uint32),
              ("dwClockFrequency", C.c_uint32),
              ("bmFramingInfo", C.c_uint8),
              ("bPreferredVersion", C.c_uint8),
              ("bMinVersion", C.c_uint8),
              ("bMaxVersion", C.c_uint8),
              ("bInterfaceNumber", C.c_uint8)]

class uvc_format_desc(C.Structure):
  pass

class uvc_frame_desc(C.Structure):
  pass

uvc_frame_desc._fields_ = [
              ("parent", C.POINTER(uvc_format_desc)),
              ("prev", C.POINTER(uvc_frame_desc)),
              ("next", C.POINTER(uvc_frame_desc)),
              # /** Type of frame, such as JPEG frame or uncompressed frme */
              ("bDescriptorSubtype", C.c_uint), # enum uvc_vs_desc_subtype bDescriptorSubtype;
              # /** Index of the frame within the list of specs available for this format */
              ("bFrameIndex", C.c_uint8),
              ("bmCapabilities", C.c_uint8),
              # /** Image width */
              ("wWidth", C.c_uint16),
              # /** Image height */
              ("wHeight", C.c_uint16),
              # /** Bitrate of corresponding stream at minimal frame rate */
              ("dwMinBitRate", C.c_uint32),
              # /** Bitrate of corresponding stream at maximal frame rate */
              ("dwMaxBitRate", C.c_uint32),
              # /** Maximum number of bytes for a video frame */
              ("dwMaxVideoFrameBufferSize", C.c_uint32),
              # /** Default frame interval (in 100ns units) */
              ("dwDefaultFrameInterval", C.c_uint32),
              # /** Minimum frame interval for continuous mode (100ns units) */
              ("dwMinFrameInterval", C.c_uint32),
              # /** Maximum frame interval for continuous mode (100ns units) */
              ("dwMaxFrameInterval", C.c_uint32),
              # /** Granularity of frame interval range for continuous mode (100ns) */
              ("dwFrameIntervalStep", C.c_uint32),
              # /** Frame intervals */
              ("bFrameIntervalType", C.c_uint8),
              # /** number of bytes per line */
              ("dwBytesPerLine", C.c_uint32),
              # /** Available frame rates, zero-terminated (in 100ns units) */
              ("intervals", C.POINTER(C.c_uint32))]

uvc_format_desc._fields_ = [
              ("parent", C.c_void_p),
              ("prev", C.POINTER(uvc_format_desc)),
              ("next", C.POINTER(uvc_format_desc)),
              # /** Type of image stream, such as JPEG or uncompressed. */
              ("bDescriptorSubtype", C.c_uint), # enum uvc_vs_desc_subtype bDescriptorSubtype;
              # /** Identifier of this format within the VS interface's format list */
              ("bFormatIndex", C.c_uint8),
              ("bNumFrameDescriptors", C.c_uint8),
              # /** Format specifier */
              ("guidFormat", C.c_char * 16), # union { uint8_t guidFormat[16]; uint8_t fourccFormat[4]; }
              # /** Format-specific data */
              ("bBitsPerPixel", C.c_uint8),
              # /** Default {uvc_frame_desc} to choose given this format */
              ("bDefaultFrameIndex", C.c_uint8),
              ("bAspectRatioX", C.c_uint8),
              ("bAspectRatioY", C.c_uint8),
              ("bmInterlaceFlags", C.c_uint8),
              ("bCopyProtect", C.c_uint8),
              ("bVariableSize", C.c_uint8),
              # /** Available frame specifications for this format */
              ("frame_descs", C.POINTER(uvc_frame_desc))]

class timeval(C.Structure):
  _fields_ = [("tv_sec", C.c_long), ("tv_usec", C.c_long)]

class uvc_frame(C.Structure):
  _fields_ = [# /** Image data for this frame */
              ("data", C.POINTER(C.c_uint8)),
              # /** Size of image data buffer */
              ("data_bytes", C.c_size_t),
              # /** Width of image in pixels */
              ("width", C.c_uint32),
              # /** Height of image in pixels */
              ("height", C.c_uint32),
              # /** Pixel data format */
              ("frame_format", C.c_uint), # enum uvc_frame_format frame_format
              # /** Number of bytes per horizontal line (undefined for compressed format) */
              ("step", C.c_size_t),
              # /** Frame number (may skip, but is strictly monotonically increasing) */
              ("sequence", C.c_uint32),
              # /** Estimate of system time when the device started capturing the image */
              ("capture_time", timeval),
              # /** Handle on the device that produced the image.
              #  * @warning You must not call any uvc_* functions during a callback. */
              ("source", C.POINTER(uvc_device)),
              # /** Is the data buffer owned by the library?
              #  * If 1, the data buffer can be arbitrarily reallocated by frame conversion
              #  * functions.
              #  * If 0, the data buffer will not be reallocated or freed by the library.
              #  * Set this field to zero if you are supplying the buffer.
              #  */
              ("library_owns_data", C.c_uint8)]

class uvc_device_handle(C.Structure):
  _fields_ = [("dev", C.POINTER(uvc_device)),
              ("prev", C.c_void_p),
              ("next", C.c_void_p),
              ("usb_devh", C.c_void_p),
              ("info", C.c_void_p),
              ("status_xfer", C.c_void_p),
              ("status_buf", C.c_ubyte * 32),
              ("status_cb", C.c_void_p),
              ("status_user_ptr", C.c_void_p),
              ("button_cb", C.c_void_p),
              ("button_user_ptr", C.c_void_p),
              ("streams", C.c_void_p),
              ("is_isight", C.c_ubyte)]

class lep_oem_sw_version(C.Structure):
  _fields_ = [("gpp_major", C.c_ubyte),
              ("gpp_minor", C.c_ubyte),
              ("gpp_build", C.c_ubyte),
              ("dsp_major", C.c_ubyte),
              ("dsp_minor", C.c_ubyte),
              ("dsp_build", C.c_ubyte),
              ("reserved", C.c_ushort)]

def call_extension_unit(devh, unit, control, data, size):
  return libuvc.uvc_get_ctrl(devh, unit, control, data, size, 0x81)

def set_extension_unit(devh, unit, control, data, size):
  return libuvc.uvc_set_ctrl(devh, unit, control, data, size, 0x81)

PT_USB_VID = 0x1e4e
PT_USB_PID = 0x0100

AGC_UNIT_ID = 3
OEM_UNIT_ID = 4
RAD_UNIT_ID = 5
SYS_UNIT_ID = 6
VID_UNIT_ID = 7

UVC_FRAME_FORMAT_UYVY = 4
UVC_FRAME_FORMAT_I420 = 5
UVC_FRAME_FORMAT_RGB = 7
UVC_FRAME_FORMAT_BGR = 8
UVC_FRAME_FORMAT_Y16 = 13

VS_FMT_GUID_GREY = C.create_string_buffer(
    b"Y8  \x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

VS_FMT_GUID_Y16 = C.create_string_buffer(
    b"Y16 \x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

VS_FMT_GUID_YUYV = C.create_string_buffer(
    b"UYVY\x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

VS_FMT_GUID_NV12 = C.create_string_buffer(
    b"NV12\x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

VS_FMT_GUID_YU12 = C.create_string_buffer(
    b"I420\x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

VS_FMT_GUID_BGR3 = C.create_string_buffer(
    b"\x7d\xeb\x36\xe4\x4f\x52\xce\x11\x9f\x53\x00\x20\xaf\x0b\xa7\x70", 16
)

VS_FMT_GUID_RGB565 = C.create_string_buffer(
    b"RGBP\x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71", 16
)

libuvc.uvc_get_format_descs.restype = C.POINTER(uvc_format_desc)

def print_device_info(devh):
  vers = lep_oem_sw_version()
  call_extension_unit(devh, OEM_UNIT_ID, 9, C.byref(vers), 8)
  print("Version gpp: {0}.{1}.{2} dsp: {3}.{4}.{5}".format(
    vers.gpp_major, vers.gpp_minor, vers.gpp_build,
    vers.dsp_major, vers.dsp_minor, vers.dsp_build,
  ))

  flir_pn = C.create_string_buffer(32)
  call_extension_unit(devh, OEM_UNIT_ID, 8, flir_pn, 32)
  print("FLIR part #: {0}".format(flir_pn.raw))

  flir_sn = C.create_string_buffer(8)
  call_extension_unit(devh, SYS_UNIT_ID, 3, flir_sn, 8)
  print("FLIR serial #: {0}".format(repr(flir_sn.raw)))

def uvc_iter_formats(devh):
  p_format_desc = libuvc.uvc_get_format_descs(devh)
  while p_format_desc:
    yield p_format_desc.contents
    p_format_desc = p_format_desc.contents.next

def uvc_iter_frames_for_format(devh, format_desc):
  p_frame_desc = format_desc.frame_descs
  while p_frame_desc:
    yield p_frame_desc.contents
    p_frame_desc = p_frame_desc.contents.next

def print_device_formats(devh):
  for format_desc in uvc_iter_formats(devh):
    print("format: {0}".format(format_desc.guidFormat[0:4]))
    for frame_desc in uvc_iter_frames_for_format(devh, format_desc):
      print("  frame {0}x{1} @ {2}fps".format(frame_desc.wWidth, frame_desc.wHeight, int(1e7 / frame_desc.dwDefaultFrameInterval)))

def uvc_get_frame_formats_by_guid(devh, vs_fmt_guid):
  for format_desc in uvc_iter_formats(devh):
    if vs_fmt_guid[0:4] == format_desc.guidFormat[0:4]:
      return [fmt for fmt in uvc_iter_frames_for_format(devh, format_desc)]
  return []
#TODO 実際のlibuvcの型定義と照らし合わせて整合性を確認
# int uvc_init(uvc_context **ctx, libusb_context *usb_ctx);
libuvc.uvc_init.restype = C.c_int
libuvc.uvc_init.argtypes = [C.POINTER(C.POINTER(uvc_context)), C.c_void_p]

# void uvc_exit(uvc_context *ctx);
libuvc.uvc_exit.restype = None
libuvc.uvc_exit.argtypes = [C.POINTER(uvc_context)]

# int uvc_find_device(uvc_context *ctx, uvc_device **dev, int vid, int pid, const char *serial)
libuvc.uvc_find_device.restype = C.c_int
libuvc.uvc_find_device.argtypes = [C.POINTER(uvc_context), C.POINTER(C.POINTER(uvc_device)), C.c_int, C.c_int, C.c_char_p]

# void uvc_unref_device(uvc_device *dev);
libuvc.uvc_unref_device.restype = None
libuvc.uvc_unref_device.argtypes = [C.POINTER(uvc_device)]

# int uvc_open(uvc_device *dev, uvc_device_handle **devh);
libuvc.uvc_open.restype = C.c_int
libuvc.uvc_open.argtypes = [C.POINTER(uvc_device), C.POINTER(C.POINTER(uvc_device_handle))]

# void uvc_close(uvc_device_handle *devh);
libuvc.uvc_close.restype = None
libuvc.uvc_close.argtypes = [C.POINTER(uvc_device_handle)]

# int uvc_get_stream_ctrl_format_size(uvc_device_handle *devh, uvc_stream_ctrl *ctrl,
#     enum uvc_frame_format format, int width, int height, int fps);
libuvc.uvc_get_stream_ctrl_format_size.restype = C.c_int
libuvc.uvc_get_stream_ctrl_format_size.argtypes = [
    C.POINTER(uvc_device_handle), C.POINTER(uvc_stream_ctrl),
    C.c_uint, C.c_int, C.c_int, C.c_int
]

# typedef void (*uvc_frame_callback_t)(struct uvc_frame *, void *user_ptr);
FRAME_CB = C.CFUNCTYPE(None, C.POINTER(uvc_frame), C.c_void_p)

# int uvc_start_streaming(uvc_device_handle *devh, uvc_stream_ctrl *ctrl,
#     uvc_frame_callback_t cb, void *user_ptr, uint8_t flags);
libuvc.uvc_start_streaming.restype = C.c_int
libuvc.uvc_start_streaming.argtypes = [
    C.POINTER(uvc_device_handle), C.POINTER(uvc_stream_ctrl),
    FRAME_CB, C.c_void_p, C.c_uint8
]

# void uvc_stop_streaming(uvc_device_handle *devh);
libuvc.uvc_stop_streaming.restype = None
libuvc.uvc_stop_streaming.argtypes = [C.POINTER(uvc_device_handle)]
