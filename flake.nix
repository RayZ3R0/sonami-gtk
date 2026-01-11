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
          name = "tidalwave-puregotk-lib-folder";
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
            ]
          );
        };
      in
      {
        devShell = pkgs.mkShell {
          PUREGOTK_LIB_FOLDER = "${libraryPath}/lib";
          GSETTINGS_SCHEMA_DIR = "./internal/settings";
          TIDAL_WAVE_DEBUG = "1";

          hardeningDisable = [ "fortify" ]; # Required for Delve
          # For delve to work, you need to add the following line to your `programs.zed-editor`:
          # package = pkgs.zed-editor.fhs;
          buildInputs = with pkgs; [
            delve
            go
            gopls
            gtk4
            librsvg
            graphviz
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
            pkg-config # Needed for the first compile with CGO
          ];
        };

        packages.tidalwave = pkgs.buildGoModule (finalAttrs: {
          pname = "tidalwave";
          version = "0.0.1";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-GXc+L59nQgHem6KdK7u0XouBxu0Ta0y57TjXPT3fMmk=";

          ldflags = [
            "-X \"codeberg.org/dergs/tidalwave/internal/ui.Commit=${(if (self ? rev) then self.rev else "")}\""
            "-X \"codeberg.org/dergs/tidalwave/internal/ui.Version=${finalAttrs.version}\""
          ];

          buildInputs = with pkgs; [
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
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
            "cmd/tidalwave"
          ];

          desktopItems = [
            (pkgs.makeDesktopItem {
              name = "org.codeberg.dergs.tidalwave";
              exec = "tidalwave %u";
              icon = "tidalwave";
              comment = "Tidal Wave is a GTK client for TIDAL written in GoLang.";
              desktopName = "Tidal Wave";
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
            wrapProgram $out/bin/tidalwave \
              --prefix GST_PLUGIN_PATH : "$GST_PLUGIN_SYSTEM_PATH_1_0" \
              --set-default PUREGOTK_LIB_FOLDER ${libraryPath}/lib \
              ''${gappsWrapperArgs[@]}
            install -Dm644 internal/icons/hicolor/scalable/apps/logo.png $out/share/icons/hicolor/scalable/apps/tidalwave.png
            install -Dm644 internal/settings/org.codeberg.dergs.tidalwave.gschema.xml -t $out/share/glib-2.0/schemas
            glib-compile-schemas $out/share/glib-2.0/schemas
          '';

          meta = {
            description = "Tidal Wave is a GTK client for TIDAL written in GoLang.";
            homepage = "https://codeberg.org/Dergs/TidalWave";
            license = pkgs.lib.licenses.gpl3;
            maintainers = with pkgs.lib.maintainers; [ nilathedragon ];
            mainProgram = "tidalwave";
          };
        });

        packages.default = self.packages.${system}.tidalwave;
      }
    );
}
