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
  krewfileContent = pkgs.writeText "krewfile" (concatStringsSep "\n" cfg.plugins);
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

    plugins = mkOption {
      type = with types; listOf str;
      default = [ ];
      description = "List of plugins to be written to the krewfile.";
    };
  };

  config = mkIf cfg.enable {

    home.packages = ([ finalPackage ] ++ optionals cfg.installKrew [ cfg.krewPackage ]);

    home.activation.krew = hm.dag.entryAfter [ "writeBoundary" ] ''
      run ${finalPackage}/bin/${finalPackage.pname} \
        -command ${cfg.krewPackage}/bin/${cfg.krewPackage.pname} \
        -file ${krewfileContent} -upgrade
    '';
  };
}
