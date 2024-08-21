{
  description = "yaaaaaaaaaaaaaaaaaaaaa";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs:
    inputs.flake-utils.lib.eachDefaultSystem (system: let
      flakePackage = flake: package: flake.packages."${system}"."${package}";
      flakeDefaultPackage = flake: flakePackage flake "default";

      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [
          (final: prev: {
            unstable = import inputs.nixpkgs-unstable {
              inherit system;
            };
          })
          # (final: prev: {
          #   webkitgtk = prev.webkitgtk.overrideAttrs (oldAttrs: {
          #     buildInputs = oldAttrs.buildInputs ++ [ prev.glib-networking ];
          #   });
          # })
        ];
      };

      gravtest = pkgs.buildGoModule {
        name = "gravtest";
        src = ./.;
        vendorHash = "";

        nativeBuildInputs = with pkgs; [
          pkg-config
          wrapGAppsHook3
        ];
        buildInputs = with pkgs; [
          webkitgtk
          glib
          glib-networking
        ];

        # subPackages = [
        # ];
      };

      windows-pkgs = inputs.nixpkgs.legacyPackages.x86_64-linux.pkgsCross.mingwW64;

      # - [fatal error: EventToken.h: No such file or directory](https://github.com/webview/webview/issues/1036)
      # - [MinGW-w64 requirements](https://github.com/webview/webview?tab=readme-ov-file#mingw-w64-requirements)
      # - [WinLibs - GCC+MinGW-w64 compiler for Windows](https://winlibs.com/#download-release)
      winlibs = windows-pkgs.stdenv.mkDerivation {
        name = "winlibs";
        src = windows-pkgs.fetchzip {
          url = "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.2.0posix-18.1.8-12.0.0-ucrt-r1/winlibs-x86_64-posix-seh-gcc-14.2.0-llvm-18.1.8-mingw-w64ucrt-12.0.0-r1.zip";
          sha256 = "sha256-xBRZ8NJmWXpvraaTpXBkd2QbhF5hR/8g/UBPwCd12hc=";
        };

        phases = ["installPhase"];
        installPhase = ''
          mkdir $out
          cp -r $src/* $out/.
        '';
      };
      mcfgthread = windows-pkgs.stdenv.mkDerivation {
        name = "mcfgthread";
        src = windows-pkgs.fetchurl {
          url = "https://mirror.msys2.org/mingw/mingw64/mingw-w64-x86_64-mcfgthread-1.8.3-1-any.pkg.tar.zst";
          sha256 = "sha256-ogfmo9utCtE2WpWtmPDuf+M6WIvpp1Xvxn+aqRu+nbs=";
        };

        nativeBuildInputs = [
          pkgs.zstd
        ];

        phases = ["installPhase"];
        installPhase = ''
          mkdir $out
          cp $src $out/src
          cd $out

          tar --zstd -xvf src
          rm src
          mv mingw64/* .
          rmdir mingw64
        '';
      };
      windows-shell = windows-pkgs.mkShell {
        nativeBuildInputs = [
          windows-pkgs.buildPackages.pkg-config
          windows-pkgs.openssl
          winlibs
          mcfgthread
        ];

        depsBuildBuild = [];
        buildInputs = [
          windows-pkgs.buildPackages.pkg-config
          windows-pkgs.openssl
          windows-pkgs.windows.mingw_w64_pthreads
          windows-pkgs.windows.pthreads
          winlibs
          mcfgthread
        ];

        env = {
          CARGO_BUILD_TARGET = "x86_64-pc-windows-gnu";
          DEV_SHELL = "WIN";
        };
      };
      fhs = pkgs.buildFHSEnv {
        name = "fhs-shell";
        targetPkgs = p: (env-packages p) ++ (custom-commands p);
        runScript = "${pkgs.zsh}/bin/zsh";
        profile = ''
          export FHS=1
          # source ./.venv/bin/activate
          # source .env
        '';
      };
      custom-commands = pkgs: [
        (pkgs.writeShellScriptBin "run" ''
          #!/usr/bin/env bash
          $PROJECT_ROOT/run.sh $@
        '')
        (pkgs.writeShellScriptBin "build-windows-app" ''
          #!/usr/bin/env bash
          nix develop .#windows -c run build-windows-app
        '')
      ];

      env-packages = pkgs:
        (with pkgs; [
          pkg-config

          go
          # go-tools
          unstable.gopls
          bun

          nodePackages_latest.typescript-language-server
          tailwindcss-language-server

          webkitgtk
          # gtk3
          # glib-networking

          libpng
          xclip
          libxkbcommon
          xorg.libXtst
          xorg.libX11
          xorg.libxcb
          xorg.xkbutils
          xorg.xcbutil
        ])
        ++ (custom-commands pkgs);
      # stdenv = pkgs.clangStdenv;
      # stdenv = pkgs.gccStdenv;
    in {
      packages = {
        default = gravtest;
        inherit gravtest winlibs mcfgthread;
      };

      devShells = {
        windows = windows-shell;

        default =
          pkgs.mkShell.override {
            # inherit stdenv;
          } {
            nativeBuildInputs = (env-packages pkgs) ++ [fhs];
            inputsFrom = [];
            shellHook = ''
              export PROJECT_ROOT="$(git rev-parse --show-toplevel)"

              # - [Workaround for blank window with WebKit/DMA-BUF/NVIDIA/X11 by SteffenL · Pull Request #1060 · webview/webview · GitHub](https://github.com/webview/webview/pull/1060)
              # export WEBKIT_DISABLE_COMPOSITING_MODE=1
              export WEBKIT_DISABLE_DMABUF_RENDERER=1

              # makes the scale "normal"
              export GDK_BACKEND=x11
            '';
          };
      };
    });
}
