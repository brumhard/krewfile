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
  args = if cfg.upgrade then "-upgrade" else "";
in
{
  options.programs.krewfile = {
    enable = mkEnableOption "krewfile";

    krewPackage = mkOption {
      type = types.package;
      default = pkgs.krew;
      defaultText = literalExpression "pkgs.krew";
      description = "Krew package to install.";
    };

    plugins = mkOption {
      type = with types; listOf str;
      default = [ ];
      defaultText = literalExpression "[ "edit-status" ]";
      description = "List of plugins to be installed.";
    };

    upgrade = mkOption {
      type = types.bool;
      default = false;
      defaultText = literalExpression "false";
      description = "Enable auto update of plugins.";
    };
  };

  config = mkIf cfg.enable {

    home.packages = [ finalPackage cfg.krewPackage ];

    home.sessionVariables.PATH = "$HOME/.krew/bin:$PATH";

    home.activation.krew = hm.dag.entryAfter [ "writeBoundary" ] ''
      echo $PATH
      PATH="$HOME/.krew/bin:$PATH"
      run ${cfg.krewPackage}/bin/${cfg.krewPackage.pname} update
      run ${finalPackage}/bin/${finalPackage.pname} \
        -command ${cfg.krewPackage}/bin/${cfg.krewPackage.pname} \
        -file ${krewfileContent} ${args}
    '';
  };
}
