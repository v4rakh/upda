{ self, lib }:
{
  config,
  lib,
  pkgs,
  ...
}:
let
  cfg = config.services.upda;
in
{
  options.services.upda = {
    enable = lib.mkEnableOption "[upda](https://git.myservermanager.com/varakh/upda) - Update Dashboard in Go";

    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.server;
      defaultText = lib.literalExpression "self.packages.\${system}.server";
      description = "The upda package to use.";
    };

    environment = lib.mkOption {
      type = lib.types.attrsOf lib.types.str;
      default = { };
      example = {
        SERVER_LISTEN = "127.0.0.1";
        SERVER_PORT = "8080";
      };
      description = ''
        Environment variables for upda. Non-sensitive values go here.
        Secrets (SECRET, AUTH_SESSION_SECRET, DB_POSTGRES_PASSWORD, etc.) must be
        set via {option}`environmentFiles` so they are not stored in the nix store.
        See [configuration reference](https://git.myservermanager.com/varakh/upda/src/branch/main/_doc/Configuration.md) for all options.
      '';
    };

    environmentFiles = lib.mkOption {
      type = lib.types.listOf lib.types.path;
      default = [ ];
      example = [ "/run/secrets/upda.env" ];
      description = ''
        Files containing additional environment variables for upda.
        Secrets such as SECRET, AUTH_SESSION_SECRET, AUTH_SESSION_PASSWORD,
        DB_POSTGRES_PASSWORD, and PROMETHEUS_SECURE_TOKEN must be provided here
        rather than in {option}`environment` to avoid storing them in the nix store.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.upda = {
      description = "upda - Update Dashboard in Go";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      environment = cfg.environment;

      serviceConfig = {
        ExecStart = "${cfg.package}/bin/upda server serve";
        EnvironmentFile = cfg.environmentFiles;
        DynamicUser = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        RestrictNamespaces = true;
        PrivateDevices = true;
        ProtectKernelTunables = true;
        ProtectKernelModules = true;
        ProtectControlGroups = true;
        LockPersonality = true;
        MemoryDenyWriteExecute = true;
        RestrictRealtime = true;
      };
    };
  };
}
