{
  description = "upda flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs@{
      flake-parts,
      nixpkgs,
      self,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      perSystem =
        { pkgs, config, ... }:
        let
          version = "8.3.1";
          frontend = pkgs.stdenv.mkDerivation (finalAttrs: {
            pname = "upda-ui";
            inherit version;
            src = ./internal/server/web;

            nativeBuildInputs = with pkgs; [
              nodejs_24
              pnpm
              pnpmConfigHook
            ];

            pnpmDeps = pkgs.fetchPnpmDeps {
              inherit (finalAttrs)
                pname
                version
                src
                ;
              fetcherVersion = 4;
              hash = "sha256-X4oHsk2NXjpCO4KNIbu1d/tenICVrYW2aogVL68jpDE=";
            };

            buildPhase = ''
              pnpm build
              rm -rf build/conf
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p $out
              cp -r build $out/
              runHook postInstall
            '';
          });
        in
        {
          packages.frontend = frontend;

          packages.server = pkgs.buildGoModule {
            pname = "upda";
            inherit version;
            src = pkgs.lib.cleanSourceWith {
              src = ./.;
              filter = path: type: !(type == "directory" && builtins.baseNameOf path == "node_modules");
            };
            tags = [ "prod" ];
            doCheck = false;
            ldflags = [
              "-s"
              "-w"
            ];
            vendorHash = "sha256-c7oIKm9Uf8vOb019mdPD50+/2Kw/7F60lV2S/vJ/8W8=";

            preBuild = ''
              mkdir -p internal/server/web
              cp -r ${frontend}/build internal/server/web
            '';
            buildInputs = [ frontend ];
          };

          packages.default = config.packages.server;

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              git-cliff
              gnumake
              go
              golangci-lint
              grype
              nodejs_24
              pnpm
            ];
          };
        };

      flake = {
        nixosModules.default = import ./nix/module.nix {
          inherit self;
          lib = nixpkgs.lib;
        };
      };
    };
}
