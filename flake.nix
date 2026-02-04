{

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;

          overlays = [
            (final: prev: {
              # Hack required to correctly register librsvg as a gdk-pixbuf loader on darwin
              librsvg = prev.librsvg.overrideAttrs (old: {
                postInstall = (prev.lib.optionalString prev.stdenv.hostPlatform.isDarwin ''
                  echo "Renaming SVG pixbuf loader for Darwin..."

                  loaderDir="$out/lib/gdk-pixbuf-2.0/2.10.0/loaders"

                  if [ -d "$loaderDir" ]; then
                    mv "$loaderDir"/libpixbufloader_svg.dylib "$loaderDir"/libpixbufloader_svg.so
                    install_name_tool -add_rpath "$out/lib" "$loaderDir"/libpixbufloader_svg.so
                  fi

                  GDK_PIXBUF_MODULEDIR="$loaderDir" ${pkgs.lib.getDev pkgs.gdk-pixbuf}/bin/gdk-pixbuf-query-loaders > "$out/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache"
                '') + (old.postInstall or "");
              });
            })
          ];
        };
        libraryPath = pkgs.symlinkJoin {
          name = "tonearm-puregotk-lib-folder";
          paths = (
            with pkgs;
            [
              cairo
              gdk-pixbuf
              glib.out
              graphene
              pango.out
              gtk4
              libadwaita
              gobject-introspection
              librsvg
              libsecret
            ]
          );
        };
      in
      {
        devShell = pkgs.mkShell {
          PUREGOTK_LIB_FOLDER = "${libraryPath}/lib";
          GSETTINGS_SCHEMA_DIR = "./internal/settings";
          TONEARM_DEBUG = "1";
          GIO_EXTRA_MODULES = pkgs.lib.optionalString pkgs.stdenv.hostPlatform.isDarwin "${pkgs.glib-networking}/lib/gio/modules";

          hardeningDisable = [ "fortify" ]; # Required for Delve
          # For delve to work, you need to add the following line to your `programs.zed-editor`:
          # package = pkgs.zed-editor.fhs;
          buildInputs = with pkgs; [
            appstream
            delve
            go
            gopls
            gtk4
            librsvg
            libsecret
            graphviz
            glib-networking
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
            pkg-config # Needed for the first compile with CGO
            cacert
          ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
            flatpak-builder
          ];
        };

        packages.tonearm = pkgs.buildGoModule (finalAttrs: {
          pname = "tonearm";
          version = "1.0.1";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-j+7cobxVGNuZFYeRn5ad7XT4um8WNWE1byFo7qo9zK0=";

          ldflags = [
            "-X \"codeberg.org/dergs/tonearm/internal/ui.Commit=${(if (self ? rev) then self.rev else "")}\""
            "-X \"codeberg.org/dergs/tonearm/internal/ui.Version=${finalAttrs.version}\""
          ];

          buildInputs = with pkgs; [
            glib-networking # TLS support for libsoup (HTTPS streaming)
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
            libsecret
          ];
          doCheck = false;
          nativeBuildInputs = with pkgs; [
            pkg-config
            gtk4
            copyDesktopItems
            makeWrapper
            wrapGAppsHook4
          ];

          subPackages = [
            "cmd/tonearm"
          ];

          desktopItems = [
            (pkgs.makeDesktopItem {
              name = "dev.dergs.Tonearm";
              exec = "tonearm %u";
              icon = "dev.dergs.Tonearm";
              comment = "Tonearm is a GTK client for TIDAL written in GoLang.";
              desktopName = "Tonearm";
              mimeTypes = [
                "x-scheme-handler/tidal"
              ];
              categories = [
                "Audio"
                "AudioVideo"
                "Music"
                "GNOME"
                "GTK"
              ];
            })
          ];

          postInstall = ''
            wrapProgram $out/bin/tonearm \
              --prefix GST_PLUGIN_PATH : "$GST_PLUGIN_SYSTEM_PATH_1_0" \
              --set-default PUREGOTK_LIB_FOLDER ${libraryPath}/lib \
              ''${gappsWrapperArgs[@]}
            install -Dm644 internal/icons/hicolor/scalable/apps/dev.dergs.Tonearm.svg -t $out/share/icons/hicolor/scalable/apps
            install -Dm644 internal/icons/hicolor/128x128/apps/dev.dergs.Tonearm.png -t $out/share/icons/hicolor/128x128/apps
            install -Dm644 internal/icons/hicolor/symbolic/apps/dev.dergs.Tonearm-symbolic.svg -t $out/share/icons/hicolor/symbolic/apps
            install -Dm644 internal/settings/dev.dergs.Tonearm.gschema.xml -t $out/share/glib-2.0/schemas
            glib-compile-schemas $out/share/glib-2.0/schemas
          '';

          meta = {
            description = "Tonearm is a GTK client for TIDAL written in GoLang.";
            homepage = "https://codeberg.org/Dergs/Tonearm";
            license = pkgs.lib.licenses.gpl3Plus;
            maintainers = with pkgs.lib.maintainers; [
              drafolin
              nilathedragon
            ];
            mainProgram = "tonearm";
          };
        });

        packages.default = self.packages.${system}.tonearm;
      }
    );
}
