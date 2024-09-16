
use std::{ffi::{c_char, CStr}, sync::Arc};

use tao::event_loop::EventLoopBuilder;

fn wry_open(url: String) -> wry::Result<()> {
    use tao::{
        event::{Event, StartCause, WindowEvent},
        event_loop::ControlFlow,
        window::WindowBuilder,
    };
    use wry::WebViewBuilder;

    #[cfg(target_os = "linux")]
    use tao::platform::unix::EventLoopBuilderExtUnix;
    #[cfg(target_os = "windows")]
    use tao::platform::windows::EventLoopBuilderExtWindows;

    let event_loop = EventLoopBuilder::new().with_any_thread(true).build();
    let window = WindowBuilder::new()
        .with_title("Gravtest")
        .build(&event_loop)
        .unwrap();

    #[cfg(target_os = "windows")]
    let builder = WebViewBuilder::new(&window);
    #[cfg(target_os = "linux")]
    let builder = {
        use tao::platform::unix::WindowExtUnix;
        use wry::WebViewBuilderExtUnix;
        let vbox = window.default_vbox().unwrap();
        WebViewBuilder::new_gtk(vbox)
    };

    let _webview = builder.with_url(url).build()?;

    event_loop.run(move |event, _, control_flow| {
        *control_flow = ControlFlow::Wait;

        match event {
            Event::NewEvents(StartCause::Init) => println!("Wry has started!"),
            Event::WindowEvent {
                event: WindowEvent::CloseRequested,
                ..
            } => *control_flow = ControlFlow::Exit,
            _ => (),
        }
    });
}

#[no_mangle]
pub extern "C" fn uritaOpenWv(url: *const c_char) -> bool {
    let url = unsafe { CStr::from_ptr(url)};
    let url = url.to_string_lossy();

    match wry_open(url.into_owned()) {
        Ok(()) => (),
        Err(e) => {
            eprintln!("{:?}", e);
            return false;
        }
    }

    true
}
