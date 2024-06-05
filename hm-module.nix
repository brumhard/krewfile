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
  krewfile = pkgs.writeText "krewfile" ''
    ${cfg.plugins}
  '';
in
{
  options.programs.krewfile = {
    enable = mkEnableOption "krewfile";

    installKrew = mkEnableOption "installKrew";

    krewPackage = mkOption {
      type = types.package;
      default = pkgs.krew;
      defaultText = literalExpression "pkgs.krew";
      description = "krew package to install.";
    };

    plugins = mkOption { type = with types; listOf str; };
  };

  config = mkIf cfg.enable {

    home.packages = ([ finalPackage ] ++ optionals cfg.installKrew [ cfg.krewPackage ]);

    home.file.".krewfile".text = ''
      ${concatStringsSep "\n" cfg.plugins}
    '';

    home.activation.krew = hm.dag.entryAfter [ "writeBoundary" ] ''
      PATH=$PATH:${cfg.krewPackage}/bin
      run ${finalPackage}/bin/${finalPackage.pname}
    '';
  };
}
