{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    nix-filter.url = "github:numtide/nix-filter";
  };

  outputs = {
    self,
    nixpkgs,
    systems,
    nix-filter,
  }: let
    eachSystem = nixpkgs.lib.genAttrs (import systems);
    filter = nix-filter.lib;
  in {
    packages = eachSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
      in {
        default = pkgs.buildGoModule {
          pname = "rose-pine-bloom";
          version = "0.0.0";

          src = filter {
            root = ./.;
            include = [
              ./go.mod
              ./go.sum
              ./main.go
              ./builder
              ./cmd
              ./color
            ];
          };

          vendorHash = "sha256-LSy43oMhDnAySHXnF1ZsBgm7Whlq1rw5XQaUcBDBhg4=";

          meta.mainProgram = "rose-pine-bloom";
        };
        rose-pine-bloom = self.packages.${system}.default;
      }
    );

    devShells = eachSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
      in {
        default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gopls
          ];
        };
        bloom = pkgs.mkShell {
          packages = [self.packages.${system}.rose-pine-bloom];
        };
      }
    );

    overlays.default = final: prev: {
      rose-pine-bloom = self.packages.${final.stdenv.hostPlatform.system}.rose-pine-bloom;
    };
  };
}
