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
        pkgs = import nixpkgs { inherit system; };
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

          hardeningDisable = [ "fortify" ]; # Required for Delve
          # For delve to work, you need to add the following line to your `programs.zed-editor`:
          # package = pkgs.zed-editor.fhs;
          buildInputs = with pkgs; [
            appstream
            delve
            flatpak-builder
            go_1_26
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
            sass
          ];
        };

        packages.tonearm = pkgs.buildGoModule.override { go = pkgs.go_1_26; } (finalAttrs: {
          pname = "sonami";
          version = "1.5.7";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-jwv80SkHVPqsWdIsVyFEw1J+8kOpg38gObYgtxnlv6o=";

          ldflags = [
            "-X \"codeberg.org/dergs/sonami/internal/ui.Commit=${(if (self ? rev) then self.rev else "")}\""
            "-X \"codeberg.org/dergs/sonami/internal/ui.Version=${finalAttrs.version}\""
          ];

          buildInputs = with pkgs; [
            glib-networking # TLS support for libsoup (HTTPS streaming)
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
            gtk4
            libsecret
          ];
          doCheck = false;
          nativeBuildInputs = with pkgs; [
            pkg-config
            copyDesktopItems
            makeWrapper
            wrapGAppsHook4
          ];

          subPackages = [
            "cmd/sonami"
          ];

          desktopItems = [
            (pkgs.makeDesktopItem {
              name = "io.github.rayz3r0.SonamiGtk";
              exec = "tonearm %u";
              icon = "io.github.rayz3r0.SonamiGtk";
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
            wrapProgram $out/bin/sonami \
              --prefix GST_PLUGIN_PATH : "$GST_PLUGIN_SYSTEM_PATH_1_0" \
              --set-default PUREGOTK_LIB_FOLDER ${libraryPath}/lib \
              ''${gappsWrapperArgs[@]}
            install -Dm644 internal/icons/hicolor/scalable/apps/io.github.rayz3r0.SonamiGtk.svg -t $out/share/icons/hicolor/scalable/apps
            install -Dm644 internal/icons/hicolor/128x128/apps/io.github.rayz3r0.SonamiGtk.png -t $out/share/icons/hicolor/128x128/apps
            install -Dm644 internal/icons/hicolor/symbolic/apps/io.github.rayz3r0.SonamiGtk-symbolic.svg -t $out/share/icons/hicolor/symbolic/apps
            install -Dm644 internal/settings/io.github.rayz3r0.SonamiGtk.gschema.xml -t $out/share/glib-2.0/schemas
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
            mainProgram = "sonami";
          };
        });

        packages.default = self.packages.${system}.tonearm;
      }
    );
}
