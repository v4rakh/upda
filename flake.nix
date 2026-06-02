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
          version = "8.3.0";
          frontend = pkgs.stdenv.mkDerivation (finalAttrs: {
            pname = "upda-ui";
            inherit version;
            src = ./internal/server/web;

            nativeBuildInputs = with pkgs; [
              nodejs_24
              pnpm
              pnpmConfigHook
            ];

            pnpmInstallFlags = [
              "--frozen-lockfile"
            ];

            pnpmDeps = pkgs.fetchPnpmDeps {
              inherit (finalAttrs)
                pname
                version
                src
                pnpmInstallFlags
                ;
              fetcherVersion = 3;
              hash = "sha256-6XRIp/Wten28PhGs39U+IEykflo3h+I3l7F8iwYf4Yg=";
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
            vendorHash = "sha256-QMqqDpr1khNbRBljDUL+UcAzJTEe6rvaXjIVx1Q99oY=";

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
