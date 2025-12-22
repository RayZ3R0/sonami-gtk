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
        cgoDependencies = with pkgs; [
          gtk4
          gtk3
          gobject-introspection
          libadwaita
        ];
        runtimeDependencies = with pkgs; [ librsvg ];
      in
      {
        devShell = pkgs.mkShell {
          buildInputs =
            with pkgs;
            [
              go
              pkg-config
            ]
            ++ cgoDependencies
            ++ runtimeDependencies;
        };

        packages.tidal-wave = pkgs.buildGoModule {
          pname = "tidalwave";
          version = "0.0.1";
          src = "./.";
          vendorHash = "sha256-QNXf32swHYKjfl8Zf5ufRi1rGtHm/N0RspqP5VJ9H34=";

          buildInputs = cgoDependencies;
          doCheck = false;
          nativeBuildInputs = with pkgs; [
            pkg-config
            copyDesktopItems
          ];

          subPackages = [
            "cmd/tidalwave"
          ];

          desktopItems = [
            (pkgs.makeDesktopItem {
              name = "org.codeberg.dergs.tidalwave";
              exec = "tidalwave";
              icon = "tidalwave";
              comment = "Tidal Wave is a GTK client for TIDAL written in GoLang.";
              desktopName = "Tidal Wave";
              categories = [
                "Audio"
                "AudioVideo"
                "Music"
                "GNOME"
                "GTK"
              ];
            })
          ];

          meta = {
            description = "Tidal Wave is a GTK client for TIDAL written in GoLang.";
            homepage = "https://codeberg.org/Dergs/TidalWave";
            maintainers = with pkgs.lib.maintainers; [ nilathedragon ];
            mainProgram = "tidalwave";
          };
        };

        packages.default = self.packages.${system}.tidal-wave;
      }
    );
}
