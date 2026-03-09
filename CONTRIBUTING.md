## Issues

We use Codeberg issues to track public bugs. Please ensure your description is clear and has sufficient instructions to be able to reproduce the issue.

## License

By contributing to Sonami, you agree that your contributions will be licensed under the LICENSE file in the root directory of this source tree.

## Code of Conduct

We expect all contributors to adhere to the [GNOME Code of Conduct](https://conduct.gnome.org/).

## Development Environment

Sonami provides a standardized development environment in the form of a direnv file as well as a Nix dev-shell. Please ensure that you have direnv and nix, the package manager, installed on your system.

We recommend using Zed to edit the codebase, as that is what we have been using for development. But anything that supports direnv should be fine.

## UI / UX

Sonami aims to follow the [GNOME Human Interface Guidelines (HIG)](https://developer.gnome.org/hig/) and use existing widgets and functionality where possible. In case you need to create something entirely new, please consider the GNOME HIG and try to follow their recommendations.

Whenever possible implement any UI components using the Schwifty UI library. Sonami is built on CGO-free GTK bindings which require manual reference management. Schwifty will take care of this for you and ensure proper memory management.

## Updating App Metadata & Screenshots
When contributing to Sonami, please ensure that you update the app metadata and screenshots accordingly. This includes updating the app icon, splash screen, and any other visual elements that may be affected by your changes.

### Screenshots
We aim to make screenshots follow the recommendations by [Flatpak](https://docs.flathub.org/docs/for-app-authors/metainfo-guidelines/quality-guidelines#screenshots) and [GNOME](https://gitlab.gnome.org/GNOME/Initiatives/-/wikis/Update-App-Screenshots). Some important recommendations include:

- Show off the app with realistic and good-looking example content
- Use GNOME system defaults (stylesheet, font, icons, window decorations, accent color, etc.)
- Do not manipulate the screenshot to make it look better than it is
- Use the GNOME screenshot tool in window mode to capture screenshots
- Write a one-line caption that describes the screenshot

It is recommended to make screenshots in the 16:9 aspect ratio. We recommend scaling the window to 1270x815px resolution when refreshing the screenshots.
