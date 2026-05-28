{
  description = "upda flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      perSystem =
        { pkgs, self', ... }:
        let
          version = "8.3.0";
          frontend = pkgs.stdenv.mkDerivation (finalAttrs: {
            pname = "upda-ui";
            inherit version;
            src = ./internal/server/web;

            nativeBuildInputs = with pkgs; [
              nodejs_24
              pnpm_10
              pnpmConfigHook
            ];

            pnpmInstallFlags = [ "--frozen-lockfile" ];

            pnpmDeps = pkgs.fetchPnpmDeps {
              inherit (finalAttrs)
                pname
                version
                src
                pnpmInstallFlags
                ;
              fetcherVersion = 3;
              hash = "sha256-DbGib0ayJZirmluQ4D8R7HfcGMKlVLsWBA/grUFIMS8=";
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
            src = ./.;
            tags = [ "prod" ];
            doCheck = false;
            vendorHash = "sha256-24UzisWDZon6u59Z+Yf2XB3rNRoEEiyd861cKgZ+f4c=";

            preBuild = ''
              mkdir -p internal/server/web
              cp -r ${frontend}/build internal/server/web
            '';
            buildInputs = [ frontend ];
          };

          packages.default = self'.packages.server;
        };
    };
}
