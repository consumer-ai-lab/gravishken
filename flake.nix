{
  description = "yaaaaaaaaaaaaaaaaaaaaa";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
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

      urita = pkgs.rustPlatform.buildRustPackage {
        name = "urita";
        src = ./urita;
        cargoLock = {
          lockFile = ./urita/Cargo.lock;
        };

        nativeBuildInputs = with pkgs; [
          pkg-config
        ];
        buildInputs = with pkgs; [
          webkitgtk_4_1
          # libsoup
        ];
      };
      gravishken = pkgs.buildGoModule {
        name = "gravishken";
        src = ./.;
        vendorHash = "";

        nativeBuildInputs = with pkgs; [
          pkg-config
          wrapGAppsHook3
          bun
        ];
        buildInputs =
          (with pkgs; [
            # webkitgtk
            glib
            glib-networking
            # gtk3

            libpng
            xclip
            libxkbcommon
            xorg.libXtst
            xorg.libX11
            xorg.libxcb
            xorg.xkbutils
            xorg.xcbutil
          ])
          ++ [
            urita
          ];

        # subPackages = [
        # ];
      };

      windows-pkgs = inputs.nixpkgs.legacyPackages.x86_64-linux.pkgsCross.mingwW64;
      rust-bin = inputs.rust-overlay.lib.mkRustBin {} windows-pkgs.buildPackages;

      windows-shell = windows-pkgs.mkShell {
        nativeBuildInputs = [
          windows-pkgs.buildPackages.pkg-config
          rust-bin.stable.latest.minimal
        ];

        depsBuildBuild = [];
        buildInputs = [
          windows-pkgs.windows.pthreads
        ];

        env = {
          CARGO_BUILD_TARGET = "x86_64-pc-windows-gnu";
          CARGO_TARGET_X86_64_PC_WINDOWS_GNU_LINKER = "${windows-pkgs.stdenv.cc.targetPrefix}cc";

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
        (pkgs.writeShellScriptBin "build-windows-installer" ''
          #!/usr/bin/env bash
          nix develop .#windows -c run build-windows-installer
        '')
        (pkgs.writeShellScriptBin "build-windows-app" ''
          #!/usr/bin/env bash
          nix develop .#windows -c run build-windows-app
        '')
        (pkgs.writeShellScriptBin "build-windows-server" ''
          #!/usr/bin/env bash
          nix develop .#windows -c run build-windows-server
        '')
      ];

      env-packages = pkgs:
        (with pkgs; [
          (python311.withPackages (ps:
            with ps; [
              pandas
              numpy
              seaborn
              matplotlib
            ]))
          python311Packages.pip
          python311Packages.virtualenv

          # go-tools
          unstable.gopls
          unstable.rust-analyzer
          unstable.rustfmt

          # easy opensource installer creator
          nsis

          nodejs

          nodePackages_latest.typescript-language-server
          tailwindcss-language-server
        ])
        ++ (custom-commands pkgs);
      # stdenv = pkgs.clangStdenv;
      # stdenv = pkgs.gccStdenv;
    in {
      packages = {
        default = gravishken;
        inherit gravishken urita;
      };

      devShells = {
        windows = windows-shell;

        default =
          pkgs.mkShell.override {
            # inherit stdenv;
          } {
            nativeBuildInputs = (env-packages pkgs) ++ [fhs];
            inputsFrom = [
              gravishken
              urita
            ];
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
