{
  description = "upda flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
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
            fetcherVersion = 2;
            hash = "sha256-qpezvwgdKrkvw6/ahJ+PQDQDmwW6LQCv+tJipT5Co/I=";
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
          vendorHash = "sha256-smmCLswh1lmTEFUbm5iMQ8Vh4NGTmWtjOMLgAOVqfww=";

          preBuild = ''
            mkdir -p internal/server/web
            cp -r ${frontend}/build internal/server/web
          '';
          buildInputs = [ frontend ];
        };

        packages.default = self.packages.${system}.server;
      }
    );
}
