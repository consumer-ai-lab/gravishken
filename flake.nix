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
        ];
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
        (pkgs.writeShellScriptBin "oof" ''
          #!/usr/bin/env bash
          cd $PROJECT_ROOT
        '')
      ];

      env-packages = pkgs:
        (with pkgs; [
          pkg-config

          go
          # go-tools
          gopls
          bun

          nodePackages_latest.typescript-language-server
          tailwindcss-language-server

          webkitgtk
          # gtk3
          # glib-networking
        ])
        ++ (custom-commands pkgs);
      # stdenv = pkgs.clangStdenv;
      # stdenv = pkgs.gccStdenv;
    in {
      packages = {
      };

      devShells = {
        default =
          pkgs.mkShell.override {
            # inherit stdenv;
          } {
            nativeBuildInputs = (env-packages pkgs) ++ [fhs];
            inputsFrom = [];
            shellHook = ''
              export PROJECT_ROOT="$(pwd)"

              export BUILD_MODE="DEV"
              # export BUILD_MODE="PROD"

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
