use std::ffi::{c_char, CStr};

fn wry_open(url: String) -> wry::Result<()> {
    use tao::{
        event::{Event, StartCause, WindowEvent},
        event_loop::{ControlFlow, EventLoopBuilder},
        window::{Fullscreen, WindowBuilder},
    };
    use wry::WebViewBuilder;

    #[cfg(target_os = "linux")]
    use tao::platform::unix::EventLoopBuilderExtUnix;
    #[cfg(target_os = "windows")]
    use tao::platform::windows::EventLoopBuilderExtWindows;

    let event_loop = EventLoopBuilder::new().with_any_thread(true).build();
    let window = WindowBuilder::new()
        .with_title("gravishken")
        .with_fullscreen(Some(Fullscreen::Borderless(None)))
        .with_decorations(false)
        .with_closable(false)
        .with_minimizable(false)
        .with_focused(true)
        .build(&event_loop)
        .unwrap();

    #[cfg(target_os = "windows")]
    let (datadir, builder) = {
        // NOTE: APPDATA is set on windows automatically
        let datadir = std::path::PathBuf::from(std::env::var("APPDATA").unwrap());
        let datadir = Some(datadir.join("Gravishken"));
        let wv = WebViewBuilder::new(&window);
        (datadir, wv)
    };
    #[cfg(target_os = "linux")]
    let (datadir, builder) = {
        use tao::platform::unix::WindowExtUnix;
        use wry::WebViewBuilderExtUnix;
        let vbox = window.default_vbox().unwrap();
        let wv = WebViewBuilder::new_gtk(vbox);
        let datadir = None;
        (datadir, wv)
    };

    // - [TODO: proxy yo!](https://docs.rs/wry/latest/wry/struct.WebViewAttributes.html#structfield.proxy_config)
    let mut ctx = wry::WebContext::new(datadir);
    let _webview = builder.with_url(url).with_web_context(&mut ctx).build()?;

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
    let url = unsafe { CStr::from_ptr(url) };
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
