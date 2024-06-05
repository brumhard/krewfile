self:
{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.programs.krewfile;
  finalPackage = self.packages.${pkgs.system}.krewfile;
  krewfileContent = concatStringsSep "\n" cfg.plugins;
in
{
  options.programs.krewfile = {
    enable = mkEnableOption "krewfile";

    installKrew = mkEnableOption "installKrew";

    krewPackage = mkOption {
      type = types.package;
      default = pkgs.krew;
      defaultText = literalExpression "pkgs.krew";
      description = ''
        Krew package to install.
        Only relevant if `programs.krewfile.installKrew` is enabled.
      '';
    };

    path = mkOption {
      type = types.string;
      default = ".krewfile";
      defaultText = literalExpression ".krewfile";
      description = ''
        Specify the path where the `krewfile` is written.
        This is relative to the home directory.
      '';
    };

    plugins = mkOption {
      type = with types; listOf str;
      default = [ ];
      description = "List of plugins to be written to the krewfile.";
    };
  };

  config = mkIf cfg.enable {

    home.packages = ([ finalPackage ] ++ optionals cfg.installKrew [ cfg.krewPackage ]);

    home.file = {
      ${cfg.path}.text = krewfileContent;
    };

    home.activation.krew = hm.dag.entryAfter [ "writeBoundary" ] ''
      PATH=$PATH:${cfg.krewPackage}/bin
      run ${finalPackage}/bin/${finalPackage.pname}
    '';
  };
}
