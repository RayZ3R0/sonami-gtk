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

        packages.tonearm = pkgs.buildGoModule (finalAttrs: {
          pname = "tonearm";
          version = "0.0.1";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-GXc+L59nQgHem6KdK7u0XouBxu0Ta0y57TjXPT3fMmk=";

          ldflags = [
            "-X \"codeberg.org/dergs/tonearm/internal/ui.Commit=${(if (self ? rev) then self.rev else "")}\""
            "-X \"codeberg.org/dergs/tonearm/internal/ui.Version=${finalAttrs.version}\""
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
            "cmd/tonearm"
          ];

          desktopItems = [
            (pkgs.makeDesktopItem {
              name = "dev.dergs.tonearm";
              exec = "tonearm %u";
              icon = "dev.dergs.tonearm";
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
            install -Dm644 internal/icons/hicolor/256x256/apps/dev.dergs.tonearm.png $out/share/icons/hicolor/256x256/apps/dev.dergs.tonearm.png
            install -Dm644 internal/settings/dev.dergs.tonearm.gschema.xml -t $out/share/glib-2.0/schemas
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
